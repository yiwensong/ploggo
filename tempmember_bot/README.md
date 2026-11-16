# tempmember-bot

A Discord bot for a server with temporary members. Anyone can be invited, and
once their time is up (24 hours by default) they're removed at the next daily
sweep, at **07:30 America/New_York** by default.

On each sweep the bot removes every member who joined at least `member_ttl_hours`
ago, **except**:

- anyone with the **admin role** (default `Admin`)
- anyone with the **permanent member role** (default `Permanent Member`)
- the server owner, and bots

It goes by Discord's own "joined server" timestamp, so it needs no database.
Someone who leaves and is re-invited starts a fresh window.

Safety: if either protected role can't be found by name (renamed or deleted),
the sweep is aborted rather than run without protection. `setup` and `run`
create both roles if they don't exist.

## Discord setup

1. Create an application at https://discord.com/developers/applications and add a bot.
2. Under **Bot**, turn on **Server Members Intent** (needed to list members).
   Consider turning off **Public Bot** so nobody else can invite it.
3. Invite it with the `bot` scope and the **Kick Members** and **Manage Roles**
   permissions:
   `https://discord.com/oauth2/authorize?client_id=<APP_ID>&scope=bot&permissions=268435458`
4. In **Server Settings → Roles**, drag the bot's role **above** the roles your
   temporary members get. Discord won't let a bot remove anyone whose highest
   role is above its own (those show up as `failed` in the logs).
5. Copy the server's ID (enable Developer Mode, right-click the server → Copy Server ID).

## Configuration

Each setting is read from its environment variable first, then from the local
store, then falls back to the default.

| Setting            | Environment variable        | Default            |
|--------------------|-----------------------------|--------------------|
| `token`            | `DISCORD_TOKEN`             | required           |
| `guild_id`         | `TEMPBOT_GUILD_ID`          | required           |
| `member_ttl_hours` | `TEMPBOT_MEMBER_TTL_HOURS`  | `24` (fractions OK)|
| `admin_role`       | `TEMPBOT_ADMIN_ROLE`        | `Admin`            |
| `permanent_role`   | `TEMPBOT_PERMANENT_ROLE`    | `Permanent Member` |
| `sweep_time`       | `TEMPBOT_SWEEP_TIME`        | `07:30` (24h)      |
| `timezone`         | `TEMPBOT_TIMEZONE`          | `America/New_York` |
| `dry_run`          | `TEMPBOT_DRY_RUN`           | `false`            |

The local store is a TOML file, `~/.config/tempmember-bot/config.toml` on Linux
(override with `TEMPBOT_CONFIG`). It's written with `0600` permissions.

```bash
tempmember-bot config set token        # prompts on stdin, keeps it out of shell history
tempmember-bot config set guild_id 123456789012345678
tempmember-bot config set member_ttl_hours 48
tempmember-bot config show             # effective values and where each came from
```

## Running

```bash
cargo build --release
./target/release/tempmember-bot setup            # check access, create roles
./target/release/tempmember-bot sweep --dry-run  # who would be removed right now
./target/release/tempmember-bot run              # daily scheduler (default command)
```

`run` only uses Discord's REST API (no gateway connection), so the bot appears
offline in the member list; that's expected. Logging level follows `RUST_LOG`
(default `info`).

Note: a sweep that's missed because the bot was down isn't made up on restart;
it waits for the next scheduled time. Run `sweep` by hand if needed.

## Tests

```bash
cargo test
```
