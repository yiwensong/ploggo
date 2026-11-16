//! Configuration loading.
//!
//! Every setting is resolved in this order:
//!   1. an environment variable (`TEMPBOT_*`, or `DISCORD_TOKEN` for the token)
//!   2. the local store, a TOML file (see [`store_path`])
//!   3. a built-in default (required settings have none)

use std::collections::BTreeMap;
use std::fs;
use std::io::Write;
use std::path::{Path, PathBuf};

use anyhow::{anyhow, bail, Context, Result};
use chrono::{Duration, NaiveTime};
use chrono_tz::Tz;

/// Environment variable that overrides the location of the local store.
pub const STORE_PATH_ENV: &str = "TEMPBOT_CONFIG";

/// A setting the bot understands: its key in the local store, the environment
/// variable that overrides it, and its default (if any).
pub struct Key {
    pub name: &'static str,
    pub env: &'static str,
    pub default: Option<&'static str>,
    pub secret: bool,
}

pub const TOKEN: Key = Key { name: "token", env: "DISCORD_TOKEN", default: None, secret: true };
pub const GUILD_ID: Key = Key { name: "guild_id", env: "TEMPBOT_GUILD_ID", default: None, secret: false };
pub const MEMBER_TTL_HOURS: Key =
    Key { name: "member_ttl_hours", env: "TEMPBOT_MEMBER_TTL_HOURS", default: Some("24"), secret: false };
pub const ADMIN_ROLE: Key =
    Key { name: "admin_role", env: "TEMPBOT_ADMIN_ROLE", default: Some("Admin"), secret: false };
pub const PERMANENT_ROLE: Key =
    Key { name: "permanent_role", env: "TEMPBOT_PERMANENT_ROLE", default: Some("Permanent Member"), secret: false };
pub const SWEEP_TIME: Key =
    Key { name: "sweep_time", env: "TEMPBOT_SWEEP_TIME", default: Some("07:30"), secret: false };
pub const TIMEZONE: Key =
    Key { name: "timezone", env: "TEMPBOT_TIMEZONE", default: Some("America/New_York"), secret: false };
pub const DRY_RUN: Key = Key { name: "dry_run", env: "TEMPBOT_DRY_RUN", default: Some("false"), secret: false };

pub const ALL_KEYS: &[&Key] =
    &[&TOKEN, &GUILD_ID, &MEMBER_TTL_HOURS, &ADMIN_ROLE, &PERMANENT_ROLE, &SWEEP_TIME, &TIMEZONE, &DRY_RUN];

pub fn find_key(name: &str) -> Option<&'static Key> {
    ALL_KEYS.iter().copied().find(|k| k.name == name)
}

/// Fully resolved bot configuration.
#[derive(Clone)]
pub struct Config {
    pub token: String,
    pub guild_id: u64,
    pub member_ttl: Duration,
    pub admin_role: String,
    pub permanent_role: String,
    pub sweep_time: NaiveTime,
    pub timezone: Tz,
    pub dry_run: bool,
}

impl std::fmt::Debug for Config {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.debug_struct("Config")
            .field("token", &"<redacted>")
            .field("guild_id", &self.guild_id)
            .field("member_ttl", &self.member_ttl)
            .field("admin_role", &self.admin_role)
            .field("permanent_role", &self.permanent_role)
            .field("sweep_time", &self.sweep_time)
            .field("timezone", &self.timezone)
            .field("dry_run", &self.dry_run)
            .finish()
    }
}

/// Where a resolved value came from (for `config show`).
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Source {
    Env,
    Store,
    Default,
    Unset,
}

/// The local store: a flat TOML table of string values.
#[derive(Debug, Default, Clone)]
pub struct Store {
    values: BTreeMap<String, String>,
}

impl Store {
    pub fn load(path: &Path) -> Result<Self> {
        if !path.exists() {
            return Ok(Self::default());
        }
        let text = fs::read_to_string(path).with_context(|| format!("reading {}", path.display()))?;
        let table: toml::Table = text.parse().with_context(|| format!("parsing {}", path.display()))?;
        let mut values = BTreeMap::new();
        for (k, v) in table {
            let s = match v {
                toml::Value::String(s) => s,
                toml::Value::Integer(i) => i.to_string(),
                toml::Value::Float(f) => f.to_string(),
                toml::Value::Boolean(b) => b.to_string(),
                other => bail!("{}: unsupported value for {k}: {other}", path.display()),
            };
            values.insert(k, s);
        }
        Ok(Self { values })
    }

    /// Writes the store, readable only by the current user (it holds the token).
    pub fn save(&self, path: &Path) -> Result<()> {
        if let Some(dir) = path.parent() {
            fs::create_dir_all(dir).with_context(|| format!("creating {}", dir.display()))?;
        }
        let table: toml::Table = self.values.iter().map(|(k, v)| (k.clone(), toml::Value::String(v.clone()))).collect();
        let text = toml::to_string(&table)?;

        let mut opts = fs::OpenOptions::new();
        opts.write(true).create(true).truncate(true);
        #[cfg(unix)]
        {
            use std::os::unix::fs::OpenOptionsExt;
            opts.mode(0o600);
        }
        let mut file = opts.open(path).with_context(|| format!("writing {}", path.display()))?;
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            file.set_permissions(fs::Permissions::from_mode(0o600))?;
        }
        file.write_all(text.as_bytes())?;
        Ok(())
    }

    pub fn get(&self, key: &str) -> Option<&str> {
        self.values.get(key).map(String::as_str)
    }

    pub fn set(&mut self, key: &str, value: String) {
        self.values.insert(key.to_string(), value);
    }

    pub fn unset(&mut self, key: &str) -> bool {
        self.values.remove(key).is_some()
    }
}

/// Default location of the local store:
/// `$TEMPBOT_CONFIG`, else `<config dir>/tempmember-bot/config.toml`
/// (`~/.config/tempmember-bot/config.toml` on Linux).
pub fn store_path() -> Result<PathBuf> {
    if let Some(p) = std::env::var_os(STORE_PATH_ENV).filter(|p| !p.is_empty()) {
        return Ok(PathBuf::from(p));
    }
    let dir = dirs::config_dir().ok_or_else(|| anyhow!("no config directory; set {STORE_PATH_ENV}"))?;
    Ok(dir.join("tempmember-bot").join("config.toml"))
}

/// Resolves settings from an environment lookup and a store.
pub struct Resolver<'a, E: Fn(&str) -> Option<String>> {
    pub env: E,
    pub store: &'a Store,
}

impl<'a, E: Fn(&str) -> Option<String>> Resolver<'a, E> {
    pub fn lookup(&self, key: &Key) -> (Option<String>, Source) {
        if let Some(v) = (self.env)(key.env).filter(|v| !v.trim().is_empty()) {
            return (Some(v), Source::Env);
        }
        if let Some(v) = self.store.get(key.name).filter(|v| !v.trim().is_empty()) {
            return (Some(v.to_string()), Source::Store);
        }
        match key.default {
            Some(d) => (Some(d.to_string()), Source::Default),
            None => (None, Source::Unset),
        }
    }

    fn required(&self, key: &Key) -> Result<String> {
        self.lookup(key).0.map(|v| v.trim().to_string()).ok_or_else(|| {
            anyhow!(
                "missing setting `{}`: set the {} environment variable or run \
                 `tempmember-bot config set {}`",
                key.name,
                key.env,
                key.name
            )
        })
    }

    pub fn resolve(&self) -> Result<Config> {
        let token = self.required(&TOKEN)?;
        let guild_id = parse_guild_id(&self.required(&GUILD_ID)?)?;
        let member_ttl = parse_ttl_hours(&self.required(&MEMBER_TTL_HOURS)?)?;
        let admin_role = self.required(&ADMIN_ROLE)?;
        let permanent_role = self.required(&PERMANENT_ROLE)?;
        let sweep_time = parse_time(&self.required(&SWEEP_TIME)?)?;
        let timezone = parse_tz(&self.required(&TIMEZONE)?)?;
        let dry_run = parse_bool(&self.required(&DRY_RUN)?).with_context(|| format!("invalid `{}`", DRY_RUN.name))?;
        Ok(Config { token, guild_id, member_ttl, admin_role, permanent_role, sweep_time, timezone, dry_run })
    }
}

/// Loads configuration from the process environment and the default store.
pub fn load() -> Result<Config> {
    let store = Store::load(&store_path()?)?;
    Resolver { env: |k: &str| std::env::var(k).ok(), store: &store }.resolve()
}

/// Validates a value for `key` before it is written to the store.
pub fn validate(key: &Key, value: &str) -> Result<()> {
    let v = value.trim();
    if v.is_empty() {
        bail!("`{}` cannot be empty", key.name);
    }
    match key.name {
        "guild_id" => parse_guild_id(v).map(|_| ()),
        "member_ttl_hours" => parse_ttl_hours(v).map(|_| ()),
        "sweep_time" => parse_time(v).map(|_| ()),
        "timezone" => parse_tz(v).map(|_| ()),
        "dry_run" => parse_bool(v).map(|_| ()),
        _ => Ok(()),
    }
}

fn parse_guild_id(v: &str) -> Result<u64> {
    match v.trim().parse::<u64>() {
        Ok(id) if id > 0 => Ok(id),
        _ => bail!("invalid guild_id {v:?}: expected the server's numeric ID"),
    }
}

fn parse_ttl_hours(v: &str) -> Result<Duration> {
    let hours: f64 = v.trim().parse().map_err(|_| anyhow!("invalid member_ttl_hours {v:?}: expected a number"))?;
    if !hours.is_finite() || hours <= 0.0 {
        bail!("invalid member_ttl_hours {v:?}: must be greater than 0");
    }
    if hours > 24.0 * 365.0 * 10.0 {
        bail!("invalid member_ttl_hours {v:?}: unreasonably large");
    }
    Ok(Duration::seconds((hours * 3600.0).round() as i64))
}

fn parse_time(v: &str) -> Result<NaiveTime> {
    NaiveTime::parse_from_str(v.trim(), "%H:%M")
        .map_err(|_| anyhow!("invalid sweep_time {v:?}: expected 24-hour HH:MM, e.g. 07:30"))
}

fn parse_tz(v: &str) -> Result<Tz> {
    v.trim().parse::<Tz>().map_err(|_| anyhow!("invalid timezone {v:?}: expected an IANA name like America/New_York"))
}

fn parse_bool(v: &str) -> Result<bool> {
    match v.trim().to_ascii_lowercase().as_str() {
        "1" | "true" | "yes" | "on" => Ok(true),
        "0" | "false" | "no" | "off" => Ok(false),
        _ => bail!("expected true or false, got {v:?}"),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashMap;

    fn env(pairs: &[(&str, &str)]) -> impl Fn(&str) -> Option<String> {
        let map: HashMap<String, String> = pairs.iter().map(|(k, v)| (k.to_string(), v.to_string())).collect();
        move |k| map.get(k).cloned()
    }

    fn store(pairs: &[(&str, &str)]) -> Store {
        let mut s = Store::default();
        for (k, v) in pairs {
            s.set(k, v.to_string());
        }
        s
    }

    #[test]
    fn env_takes_precedence_over_store() {
        let s = store(&[("token", "from-store"), ("guild_id", "1")]);
        let cfg = Resolver { env: env(&[("DISCORD_TOKEN", "from-env")]), store: &s }.resolve().unwrap();
        assert_eq!(cfg.token, "from-env");
        assert_eq!(cfg.guild_id, 1);
    }

    #[test]
    fn store_used_when_env_missing_or_blank() {
        let s = store(&[("token", "from-store"), ("guild_id", "42"), ("member_ttl_hours", "48")]);
        let cfg = Resolver { env: env(&[("DISCORD_TOKEN", "  ")]), store: &s }.resolve().unwrap();
        assert_eq!(cfg.token, "from-store");
        assert_eq!(cfg.member_ttl, Duration::hours(48));
    }

    #[test]
    fn defaults() {
        let s = store(&[("token", "t"), ("guild_id", "42")]);
        let cfg = Resolver { env: env(&[]), store: &s }.resolve().unwrap();
        assert_eq!(cfg.member_ttl, Duration::hours(24));
        assert_eq!(cfg.admin_role, "Admin");
        assert_eq!(cfg.permanent_role, "Permanent Member");
        assert_eq!(cfg.sweep_time, NaiveTime::from_hms_opt(7, 30, 0).unwrap());
        assert_eq!(cfg.timezone, chrono_tz::America::New_York);
        assert!(!cfg.dry_run);
    }

    #[test]
    fn missing_token_is_an_error() {
        let s = store(&[("guild_id", "42")]);
        let err = Resolver { env: env(&[]), store: &s }.resolve().unwrap_err().to_string();
        assert!(err.contains("DISCORD_TOKEN"), "{err}");
    }

    #[test]
    fn fractional_ttl() {
        assert_eq!(parse_ttl_hours("0.5").unwrap(), Duration::minutes(30));
        assert!(parse_ttl_hours("0").is_err());
        assert!(parse_ttl_hours("-3").is_err());
        assert!(parse_ttl_hours("abc").is_err());
    }

    #[test]
    fn validation() {
        assert!(validate(&SWEEP_TIME, "7:30").is_ok());
        assert!(validate(&SWEEP_TIME, "25:00").is_err());
        assert!(validate(&TIMEZONE, "Mars/Olympus").is_err());
        assert!(validate(&GUILD_ID, "abc").is_err());
        assert!(validate(&DRY_RUN, "maybe").is_err());
    }

    #[test]
    fn debug_redacts_token() {
        let s = store(&[("token", "super-secret"), ("guild_id", "42")]);
        let cfg = Resolver { env: env(&[]), store: &s }.resolve().unwrap();
        assert!(!format!("{cfg:?}").contains("super-secret"));
    }

    #[test]
    fn store_round_trip_and_permissions() {
        let dir = tempfile::tempdir().unwrap();
        let path = dir.path().join("nested").join("config.toml");
        let mut s = Store::default();
        s.set("token", "abc".into());
        s.set("member_ttl_hours", "12".into());
        s.save(&path).unwrap();

        let loaded = Store::load(&path).unwrap();
        assert_eq!(loaded.get("token"), Some("abc"));
        assert_eq!(loaded.get("member_ttl_hours"), Some("12"));

        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let mode = fs::metadata(&path).unwrap().permissions().mode() & 0o777;
            assert_eq!(mode, 0o600);
        }
    }

    #[test]
    fn store_accepts_native_toml_types() {
        let dir = tempfile::tempdir().unwrap();
        let path = dir.path().join("config.toml");
        fs::write(&path, "guild_id = 123\nmember_ttl_hours = 1.5\ndry_run = true\n").unwrap();
        let s = Store::load(&path).unwrap();
        assert_eq!(s.get("guild_id"), Some("123"));
        assert_eq!(s.get("member_ttl_hours"), Some("1.5"));
        assert_eq!(s.get("dry_run"), Some("true"));
    }
}
