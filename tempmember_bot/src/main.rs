mod bot;
mod config;
mod policy;
mod schedule;

use std::io::{BufRead, IsTerminal};

use anyhow::{anyhow, Context, Result};
use clap::{Parser, Subcommand};
use tracing_subscriber::EnvFilter;

/// Removes temporary members from a Discord server once their window is up.
#[derive(Parser)]
#[command(version, about)]
struct Cli {
    #[command(subcommand)]
    command: Option<Command>,
}

#[derive(Subcommand)]
enum Command {
    /// Run the bot: sweep every day at the configured time (the default).
    Run,
    /// Run a single sweep now and exit.
    Sweep {
        /// Only report who would be removed; remove nobody.
        #[arg(long)]
        dry_run: bool,
    },
    /// Check the token and server, and create the protected roles if missing.
    Setup,
    /// Read or change the local config store.
    #[command(subcommand)]
    Config(ConfigCommand),
}

#[derive(Subcommand)]
enum ConfigCommand {
    /// Print the store's location.
    Path,
    /// Show every setting, its effective value and where it came from.
    Show,
    /// Save a setting to the store. Omit the value to read it from stdin
    /// (recommended for the token, so it stays out of shell history).
    Set { key: String, value: Option<String> },
    /// Remove a setting from the store.
    Unset { key: String },
}

#[tokio::main]
async fn main() -> std::process::ExitCode {
    match real_main().await {
        Ok(()) => std::process::ExitCode::SUCCESS,
        Err(e) => {
            eprintln!("error: {e:#}");
            std::process::ExitCode::FAILURE
        }
    }
}

async fn real_main() -> Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::try_from_default_env().unwrap_or_else(|_| EnvFilter::new("info")))
        .init();

    match Cli::parse().command.unwrap_or(Command::Run) {
        Command::Run => bot::run(config::load()?).await,
        Command::Sweep { dry_run } => {
            let cfg = config::load()?;
            let http = bot::http(&cfg);
            bot::check_access(&http, &cfg).await?;
            let report = bot::sweep(&http, &cfg, dry_run || cfg.dry_run).await?;
            for name in &report.would_kick {
                println!("would remove: {name}");
            }
            for name in &report.kicked {
                println!("removed: {name}");
            }
            for name in &report.failed {
                println!("FAILED to remove: {name}");
            }
            Ok(())
        }
        Command::Setup => {
            let cfg = config::load()?;
            let http = bot::http(&cfg);
            bot::check_access(&http, &cfg).await?;
            bot::ensure_roles(&http, &cfg).await?;
            println!("ok: roles {:?} and {:?} are in place", cfg.admin_role, cfg.permanent_role);
            Ok(())
        }
        Command::Config(cmd) => config_command(cmd),
    }
}

fn config_command(cmd: ConfigCommand) -> Result<()> {
    let path = config::store_path()?;
    match cmd {
        ConfigCommand::Path => println!("{}", path.display()),
        ConfigCommand::Show => {
            let store = config::Store::load(&path)?;
            let resolver = config::Resolver { env: |k: &str| std::env::var(k).ok(), store: &store };
            println!("store: {}", path.display());
            for key in config::ALL_KEYS {
                let (value, source) = resolver.lookup(key);
                let shown = match (&value, key.secret) {
                    (None, _) => "<unset>".to_string(),
                    (Some(_), true) => "<set>".to_string(),
                    (Some(v), false) => v.clone(),
                };
                println!("{:<18} {:<28} ({:?}, env {})", key.name, shown, source, key.env);
            }
        }
        ConfigCommand::Set { key, value } => {
            let k = lookup_key(&key)?;
            let value = match value {
                Some(v) => v,
                None => {
                    if std::io::stdin().is_terminal() {
                        eprintln!("enter value for {key}, then press enter:");
                    }
                    let mut line = String::new();
                    std::io::stdin().lock().read_line(&mut line).context("reading stdin")?;
                    line
                }
            };
            let value = value.trim().to_string();
            config::validate(k, &value)?;
            let mut store = config::Store::load(&path)?;
            store.set(k.name, value);
            store.save(&path)?;
            println!("saved {} to {}", k.name, path.display());
        }
        ConfigCommand::Unset { key } => {
            let k = lookup_key(&key)?;
            let mut store = config::Store::load(&path)?;
            if store.unset(k.name) {
                store.save(&path)?;
                println!("removed {} from {}", k.name, path.display());
            } else {
                println!("{} was not set in {}", k.name, path.display());
            }
        }
    }
    Ok(())
}

fn lookup_key(name: &str) -> Result<&'static config::Key> {
    config::find_key(name).ok_or_else(|| {
        let names: Vec<_> = config::ALL_KEYS.iter().map(|k| k.name).collect();
        anyhow!("unknown setting {name:?}; expected one of: {}", names.join(", "))
    })
}
