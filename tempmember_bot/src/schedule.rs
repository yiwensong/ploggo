//! Computing the next sweep time in a given time zone (DST-aware).

use chrono::{DateTime, Duration, LocalResult, NaiveDate, NaiveTime, TimeZone, Utc};
use chrono_tz::Tz;

/// The first instant strictly after `now` whose wall-clock time in `tz` is `at`.
///
/// If `at` falls in a DST gap on some day, that day's sweep happens at the first
/// valid instant after the gap; if it is ambiguous (DST ends), the earlier one is used.
pub fn next_run(now: DateTime<Utc>, at: NaiveTime, tz: Tz) -> DateTime<Utc> {
    let today = now.with_timezone(&tz).date_naive();
    // Start a day back so a sweep moved by DST handling can't be skipped.
    let mut day = today - Duration::days(1);
    loop {
        let candidate = instant_on(day, at, tz);
        if candidate > now {
            return candidate;
        }
        day = day.succ_opt().expect("date overflow");
    }
}

fn instant_on(day: NaiveDate, at: NaiveTime, tz: Tz) -> DateTime<Utc> {
    let mut naive = day.and_time(at);
    loop {
        match tz.from_local_datetime(&naive) {
            LocalResult::Single(t) => return t.with_timezone(&Utc),
            LocalResult::Ambiguous(earliest, _) => return earliest.with_timezone(&Utc),
            // Inside a DST gap: walk forward a minute at a time until it's valid.
            LocalResult::None => naive += Duration::minutes(1),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use chrono_tz::America::New_York;

    fn t(h: u32, m: u32) -> NaiveTime {
        NaiveTime::from_hms_opt(h, m, 0).unwrap()
    }

    fn utc(y: i32, mo: u32, d: u32, h: u32, mi: u32) -> DateTime<Utc> {
        Utc.with_ymd_and_hms(y, mo, d, h, mi, 0).unwrap()
    }

    #[test]
    fn later_today_in_winter() {
        // 06:00 EST == 11:00 UTC; next 07:30 EST is 12:30 UTC the same day.
        assert_eq!(next_run(utc(2026, 1, 15, 11, 0), t(7, 30), New_York), utc(2026, 1, 15, 12, 30));
    }

    #[test]
    fn tomorrow_when_already_past() {
        assert_eq!(next_run(utc(2026, 1, 15, 13, 0), t(7, 30), New_York), utc(2026, 1, 16, 12, 30));
    }

    #[test]
    fn exactly_now_schedules_tomorrow() {
        assert_eq!(next_run(utc(2026, 1, 15, 12, 30), t(7, 30), New_York), utc(2026, 1, 16, 12, 30));
    }

    #[test]
    fn summer_uses_edt() {
        // 07:30 EDT == 11:30 UTC.
        assert_eq!(next_run(utc(2026, 7, 1, 0, 0), t(7, 30), New_York), utc(2026, 7, 1, 11, 30));
    }

    #[test]
    fn across_spring_forward() {
        // DST starts 2026-03-08. Sweep on the 7th is 12:30 UTC, on the 8th 11:30 UTC.
        assert_eq!(next_run(utc(2026, 3, 7, 13, 0), t(7, 30), New_York), utc(2026, 3, 8, 11, 30));
    }

    #[test]
    fn across_fall_back() {
        // DST ends 2026-11-01. Sweep on Oct 31 is 11:30 UTC, on Nov 1 12:30 UTC.
        assert_eq!(next_run(utc(2026, 10, 31, 12, 0), t(7, 30), New_York), utc(2026, 11, 1, 12, 30));
    }

    #[test]
    fn time_in_dst_gap_runs_after_gap() {
        // 02:30 doesn't exist on 2026-03-08 in New York; run at 03:00 EDT (07:00 UTC).
        assert_eq!(next_run(utc(2026, 3, 8, 5, 0), t(2, 30), New_York), utc(2026, 3, 8, 7, 0));
    }

    #[test]
    fn ambiguous_time_uses_first_occurrence() {
        // 01:30 happens twice on 2026-11-01; the first is 01:30 EDT == 05:30 UTC.
        assert_eq!(next_run(utc(2026, 11, 1, 4, 0), t(1, 30), New_York), utc(2026, 11, 1, 5, 30));
    }
}
