package timestamp

import (
	"strconv"
	"time"
)

type Timestamp int64

func (ts Timestamp) Unix() int64 {
	return int64(ts)
}

func (ts Timestamp) Before(b Timestamp) bool {
	return ts < b
}

func (ts Timestamp) BeforeTime(b time.Time) bool {
	return ts < Timestamp(b.Unix())
}

func (ts Timestamp) After(b Timestamp) bool {
	return ts > b
}

func (ts Timestamp) AfterTime(b time.Time) bool {
	return ts > Timestamp(b.Unix())
}

func (ts Timestamp) Equal(b Timestamp) bool {
	return ts == b
}

func (ts Timestamp) EqualTime(b time.Time) bool {
	return ts == Timestamp(b.Unix())
}

func (ts Timestamp) IsZero() bool {
	return ts == 0
}

func (ts Timestamp) Parse() time.Time {
	return time.Unix(int64(ts), 0)
}

func (ts Timestamp) Date() (day, month, year int) {
	y, m, d := ts.Parse().Date()
	return d, int(m), y
}

func (ts Timestamp) Year() int {
	return ts.Parse().Year()
}

func (ts Timestamp) Month() int {
	return int(ts.Parse().Month())
}

func (ts Timestamp) Day() int {
	return ts.Parse().Day()
}

// returns the first and last timestamp of the matching day
func (ts Timestamp) Range() (Timestamp, Timestamp) {
	y, m, d := ts.Date()
	from := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC).Unix()
	to := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1).Unix() - 1
	return Timestamp(from), Timestamp(to)
}

func (ts Timestamp) Overlap(b Timestamp) bool {
	ay, am, ad := ts.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func (ts Timestamp) AddDate(years, months, days int) Timestamp {
	return FromTime(ts.Parse().AddDate(years, months, days))
}

func (ts Timestamp) Add(d time.Duration) Timestamp {
	return FromTime(ts.Parse().Add(d))
}

func (ts Timestamp) MidDay() Timestamp {
	y, m, d := ts.Date()
	date := time.Date(y, time.Month(m), d, 12, 0, 0, 0, time.UTC)
	return FromTime(date)
}

func FromTime(t time.Time) Timestamp {
	return Timestamp(t.Unix())
}

// panics if year is negative, month is out of 1-12 range and day out of 1-31 range.
func FromDate(day, month, year int) Timestamp {
	if year < 0 {
		panic("year cannot be negative")
	}
	if month < 1 || month > 12 {
		panic("month must be within 1-12 range")
	}
	if day < 1 || day > 31 {
		panic("day must be within 1-31 range")
	}

	var t time.Time
	t = t.AddDate(year, month, day)

	return FromTime(t)
}

func Now() Timestamp {
	return Timestamp(time.Now().Unix())
}

func Blank() Timestamp {
	return 0
}

func (ts Timestamp) Format(layout string) string {
	return ts.Parse().Format(layout)
}

func (ts Timestamp) String() string {
	return strconv.FormatInt(ts.Unix(), 10)
}
