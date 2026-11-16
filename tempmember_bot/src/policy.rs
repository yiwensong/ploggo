//! Pure decision logic: who gets removed during a sweep.

use std::collections::HashSet;

use chrono::{DateTime, Duration, Utc};

/// The parts of a guild member the policy cares about.
#[derive(Debug, Clone)]
pub struct MemberInfo {
    pub user_id: u64,
    pub display: String,
    pub is_bot: bool,
    pub joined_at: Option<DateTime<Utc>>,
    pub role_ids: Vec<u64>,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Decision {
    Kick,
    Keep(KeepReason),
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum KeepReason {
    Owner,
    Bot,
    ProtectedRole,
    UnknownJoinTime,
    WithinWindow,
}

pub struct Policy {
    /// Members who joined at or before `now - member_ttl` are removed.
    pub member_ttl: Duration,
    /// Holding any of these roles (admin / permanent member) protects a member.
    pub protected_role_ids: HashSet<u64>,
    pub owner_id: u64,
}

impl Policy {
    pub fn decide(&self, member: &MemberInfo, now: DateTime<Utc>) -> Decision {
        if member.user_id == self.owner_id {
            return Decision::Keep(KeepReason::Owner);
        }
        if member.is_bot {
            return Decision::Keep(KeepReason::Bot);
        }
        if member.role_ids.iter().any(|r| self.protected_role_ids.contains(r)) {
            return Decision::Keep(KeepReason::ProtectedRole);
        }
        let Some(joined_at) = member.joined_at else {
            // Never kick someone when we can't tell how long they've been here.
            return Decision::Keep(KeepReason::UnknownJoinTime);
        };
        if joined_at <= now - self.member_ttl {
            Decision::Kick
        } else {
            Decision::Keep(KeepReason::WithinWindow)
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::TimeZone;

    const ADMIN: u64 = 100;
    const PERMANENT: u64 = 200;
    const OTHER: u64 = 300;

    fn now() -> DateTime<Utc> {
        Utc.with_ymd_and_hms(2026, 11, 16, 12, 30, 0).unwrap()
    }

    fn policy() -> Policy {
        Policy { member_ttl: Duration::hours(24), protected_role_ids: [ADMIN, PERMANENT].into(), owner_id: 1 }
    }

    fn member(id: u64, hours_ago: i64, roles: &[u64]) -> MemberInfo {
        MemberInfo {
            user_id: id,
            display: format!("user{id}"),
            is_bot: false,
            joined_at: Some(now() - Duration::hours(hours_ago)),
            role_ids: roles.to_vec(),
        }
    }

    #[test]
    fn kicks_members_past_the_window() {
        assert_eq!(policy().decide(&member(2, 25, &[]), now()), Decision::Kick);
        assert_eq!(policy().decide(&member(2, 24 * 30, &[OTHER]), now()), Decision::Kick);
    }

    #[test]
    fn exactly_at_the_boundary_is_kicked() {
        assert_eq!(policy().decide(&member(2, 24, &[]), now()), Decision::Kick);
    }

    #[test]
    fn keeps_recent_members() {
        assert_eq!(policy().decide(&member(2, 23, &[]), now()), Decision::Keep(KeepReason::WithinWindow));
    }

    #[test]
    fn never_kicks_admins_or_permanent_members() {
        let p = policy();
        assert_eq!(p.decide(&member(2, 1000, &[ADMIN]), now()), Decision::Keep(KeepReason::ProtectedRole));
        assert_eq!(p.decide(&member(3, 1000, &[OTHER, PERMANENT]), now()), Decision::Keep(KeepReason::ProtectedRole));
    }

    #[test]
    fn never_kicks_owner_or_bots() {
        let p = policy();
        assert_eq!(p.decide(&member(1, 1000, &[]), now()), Decision::Keep(KeepReason::Owner));
        let mut bot = member(5, 1000, &[]);
        bot.is_bot = true;
        assert_eq!(p.decide(&bot, now()), Decision::Keep(KeepReason::Bot));
    }

    #[test]
    fn unknown_join_time_is_kept() {
        let mut m = member(2, 1000, &[]);
        m.joined_at = None;
        assert_eq!(policy().decide(&m, now()), Decision::Keep(KeepReason::UnknownJoinTime));
    }

    #[test]
    fn configurable_ttl() {
        let p = Policy { member_ttl: Duration::hours(72), ..policy() };
        assert_eq!(p.decide(&member(2, 48, &[]), now()), Decision::Keep(KeepReason::WithinWindow));
        assert_eq!(p.decide(&member(2, 73, &[]), now()), Decision::Kick);
    }
}
