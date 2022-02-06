package timetools

import (
	"time"
)

// takes time t and returns the monday 00:00:00 and sunday 23:59:59 of the week the t is in.
func WeekRange(t time.Time) (mon, sun time.Time) {
	switch t.Weekday() {
	case time.Sunday:
		mon = t.AddDate(0, 0, -6)
		sun = t
	case time.Monday:
		mon = t
		sun = t.AddDate(0, 0, 6)
	case time.Tuesday:
		mon = t.AddDate(0, 0, -1)
		sun = t.AddDate(0, 0, 5)
	case time.Wednesday:
		mon = t.AddDate(0, 0, -2)
		sun = t.AddDate(0, 0, 4)
	case time.Thursday:
		mon = t.AddDate(0, 0, -3)
		sun = t.AddDate(0, 0, 3)
	case time.Friday:
		mon = t.AddDate(0, 0, -4)
		sun = t.AddDate(0, 0, 2)
	case time.Saturday:
		mon = t.AddDate(0, 0, -5)
		sun = t.AddDate(0, 0, 1)
	}

	// reset to 00:00:00
	mon = mon.Add((time.Duration(mon.Hour())*time.Hour + time.Duration(mon.Minute())*time.Minute + time.Duration(mon.Second())*time.Second + time.Duration(mon.Nanosecond())*time.Nanosecond) * -1)
	// reset to 00:00:00
	sun = sun.Add((time.Duration(sun.Hour())*time.Hour + time.Duration(sun.Minute())*time.Minute + time.Duration(sun.Second())*time.Second + time.Duration(sun.Nanosecond())*time.Nanosecond) * -1)
	// add 23:59:59.000
	sun = sun.Add(time.Duration(23)*time.Hour + time.Duration(59)*time.Minute + time.Duration(59)*time.Second)

	return
}

// same as WeekRange but it works just with year and week number
func WeekRangeFromISO(y, w int) (mon, sun time.Time) {
	return WeekRange(WeekStart(y, w))

}

// returns the 00:00:00:00 monday for the provided week
// copied from github.com/icza/gox
// https://stackoverflow.com/a/52303730
func WeekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

// returns the "proper" day of the week ranging from 1 to 7 where 1 is monday and 7 is sunday.
func WeekDay(t time.Time) int {
	return NormalizeWeekday(t.Weekday())
}

// normalizes the time's weekday into "proper" day of the week ranging from 1 to 7 where 1 is monday and 7 is sunday.
func NormalizeWeekday(d time.Weekday) int {
	days := map[time.Weekday]int{
		time.Monday:    1,
		time.Tuesday:   2,
		time.Wednesday: 3,
		time.Thursday:  4,
		time.Friday:    5,
		time.Saturday:  6,
		time.Sunday:    7,
	}
	return days[d]
}

// normalizes the time's weekday into "proper" day of the week ranging from 1 to 7 where 1 is monday and 7 is sunday.
func NormalizeWeekdayNumber(d int) int {
	return NormalizeWeekday(time.Weekday(d))
}

// returns a timestamp of the beginning and an end of the day of provided time t
func DayRangeTimestamp(t time.Time) (from, to int64) {
	// reset to 00:00:00
	t = t.Add((time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second + time.Duration(t.Nanosecond())*time.Nanosecond) * -1)
	from = t.Unix()
	// add 23:59:59
	t = t.Add(time.Duration(23)*time.Hour + time.Duration(59)*time.Minute + time.Duration(59)*time.Second)
	to = t.Unix()
	return
}
