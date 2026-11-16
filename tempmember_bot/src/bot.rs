//! Talking to Discord: role setup, sweeping, and the daily scheduler loop.
//!
//! The bot only uses Discord's REST API (no gateway connection), so it needs no
//! event handling: it wakes up at the sweep time, lists the server's members,
//! and removes the ones whose time is up.

use std::collections::HashSet;

use anyhow::{bail, Context, Result};
use chrono::{DateTime, Utc};
use serenity::builder::EditRole;
use serenity::http::Http;
use serenity::model::guild::Member;
use serenity::model::id::{GuildId, UserId};
use tracing::{error, info, warn};

use crate::config::Config;
use crate::policy::{Decision, KeepReason, MemberInfo, Policy};
use crate::schedule::next_run;

/// Discord's maximum page size for listing guild members.
const MEMBER_PAGE_SIZE: u64 = 1000;
/// Longest single sleep while waiting for the next sweep, so clock changes or a
/// suspended host don't push the sweep far past its time.
const MAX_SLEEP: std::time::Duration = std::time::Duration::from_secs(300);

#[derive(Debug, Default)]
pub struct SweepReport {
    pub scanned: usize,
    pub kicked: Vec<String>,
    pub would_kick: Vec<String>,
    pub failed: Vec<String>,
    pub protected: usize,
    pub within_window: usize,
    pub skipped_other: usize,
}

pub fn http(cfg: &Config) -> Http {
    Http::new(&cfg.token)
}

/// Checks the token works and the bot can see the configured server.
pub async fn check_access(http: &Http, cfg: &Config) -> Result<()> {
    let me = http.get_current_user().await.context("logging in to Discord (is the token right?)")?;
    let guild = http
        .get_guild(GuildId::new(cfg.guild_id))
        .await
        .with_context(|| format!("reading server {} (has the bot been invited to it?)", cfg.guild_id))?;
    info!(bot = %me.name, server = %guild.name, "connected");
    Ok(())
}

/// Creates the admin and permanent-member roles if the server doesn't have them yet.
pub async fn ensure_roles(http: &Http, cfg: &Config) -> Result<()> {
    let guild_id = GuildId::new(cfg.guild_id);
    let roles = http.get_guild_roles(guild_id).await.context("listing server roles")?;
    for name in [&cfg.admin_role, &cfg.permanent_role] {
        if roles.iter().any(|r| &r.name == name) {
            continue;
        }
        info!(role = %name, "creating missing role");
        guild_id
            .create_role(http, EditRole::new().name(name.as_str()).audit_log_reason("tempmember-bot setup"))
            .await
            .with_context(|| format!("creating role {name:?} (does the bot have Manage Roles?)"))?;
    }
    Ok(())
}

/// Removes every member whose window has passed and who has no protected role.
pub async fn sweep(http: &Http, cfg: &Config, dry_run: bool) -> Result<SweepReport> {
    let guild_id = GuildId::new(cfg.guild_id);
    let guild = http.get_guild(guild_id).await.context("reading server")?;

    // Resolve protected roles by name. Refuse to sweep if either is missing:
    // a renamed or deleted role must never turn into "kick everyone".
    let roles = http.get_guild_roles(guild_id).await.context("listing server roles")?;
    let mut protected_role_ids = HashSet::new();
    for name in [&cfg.admin_role, &cfg.permanent_role] {
        let ids: Vec<u64> = roles.iter().filter(|r| &r.name == name).map(|r| r.id.get()).collect();
        if ids.is_empty() {
            bail!("protected role {name:?} not found in the server; refusing to sweep");
        }
        protected_role_ids.extend(ids);
    }

    let policy = Policy { member_ttl: cfg.member_ttl, protected_role_ids, owner_id: guild.owner_id.get() };
    let members = list_members(http, guild_id).await?;
    let now = Utc::now();
    let reason = format!("Temporary membership expired ({} window)", format_ttl(cfg));

    let mut report = SweepReport { scanned: members.len(), ..Default::default() };
    for m in &members {
        match policy.decide(m, now) {
            Decision::Keep(KeepReason::ProtectedRole) => report.protected += 1,
            Decision::Keep(KeepReason::WithinWindow) => report.within_window += 1,
            Decision::Keep(_) => report.skipped_other += 1,
            Decision::Kick if dry_run => {
                info!(member = %m.display, joined_at = ?m.joined_at, "would remove (dry run)");
                report.would_kick.push(m.display.clone());
            }
            Decision::Kick => match http.kick_member(guild_id, UserId::new(m.user_id), Some(&reason)).await {
                Ok(()) => {
                    info!(member = %m.display, joined_at = ?m.joined_at, "removed");
                    report.kicked.push(m.display.clone());
                }
                Err(e) => {
                    // Usually: the member's highest role is above the bot's.
                    warn!(member = %m.display, error = %e, "failed to remove");
                    report.failed.push(m.display.clone());
                }
            },
        }
    }
    info!(
        scanned = report.scanned,
        removed = report.kicked.len(),
        would_remove = report.would_kick.len(),
        failed = report.failed.len(),
        protected = report.protected,
        within_window = report.within_window,
        "sweep finished"
    );
    Ok(report)
}

async fn list_members(http: &Http, guild_id: GuildId) -> Result<Vec<MemberInfo>> {
    let mut out = Vec::new();
    let mut after = None;
    loop {
        let page = http
            .get_guild_members(guild_id, Some(MEMBER_PAGE_SIZE), after)
            .await
            .context("listing members (is the Server Members Intent enabled for the bot?)")?;
        let len = page.len() as u64;
        after = page.iter().map(|m| m.user.id.get()).max();
        out.extend(page.iter().map(to_info));
        if len < MEMBER_PAGE_SIZE {
            return Ok(out);
        }
    }
}

fn to_info(m: &Member) -> MemberInfo {
    let display = match &m.nick {
        Some(nick) => format!("{} ({}, {})", nick, m.user.name, m.user.id),
        None => format!("{} ({})", m.user.name, m.user.id),
    };
    MemberInfo {
        user_id: m.user.id.get(),
        display,
        is_bot: m.user.bot,
        joined_at: m.joined_at.and_then(|t| DateTime::from_timestamp(t.unix_timestamp(), 0)),
        role_ids: m.roles.iter().map(|r| r.get()).collect(),
    }
}

fn format_ttl(cfg: &Config) -> String {
    let mins = cfg.member_ttl.num_minutes();
    if mins % 60 == 0 {
        format!("{}h", mins / 60)
    } else {
        format!("{}h{}m", mins / 60, mins % 60)
    }
}

/// Runs forever: sweeps once a day at the configured local time.
pub async fn run(cfg: Config) -> Result<()> {
    let http = http(&cfg);
    check_access(&http, &cfg).await?;
    ensure_roles(&http, &cfg).await?;
    info!(?cfg, "starting scheduler");

    loop {
        let next = next_run(Utc::now(), cfg.sweep_time, cfg.timezone);
        info!(next_sweep = %next.with_timezone(&cfg.timezone), "waiting for next sweep");
        loop {
            let remaining = next - Utc::now();
            if remaining <= chrono::Duration::zero() {
                break;
            }
            let nap = remaining.to_std().unwrap_or_default().min(MAX_SLEEP);
            tokio::select! {
                _ = tokio::time::sleep(nap) => {}
                _ = tokio::signal::ctrl_c() => {
                    info!("shutting down");
                    return Ok(());
                }
            }
        }
        if let Err(e) = sweep(&http, &cfg, cfg.dry_run).await {
            error!(error = format!("{e:#}"), "sweep failed");
        }
    }
}
