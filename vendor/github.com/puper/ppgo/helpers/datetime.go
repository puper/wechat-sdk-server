package helpers

import (
	"time"
)

func FormatDatetime(sec int64) string {
	return time.Unix(sec, 0).Format("2006-01-02 15:04:05")
}

func ParseDatetime(d string) (int64, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", d, time.Local)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func ParseTime(d string) (int64, error) {
	t, err := time.ParseInLocation("15:04:05", d, time.Local)
	if err != nil {
		return 0, err
	}
	dayStart, err := DayStart(t)
	if err != nil {
		return 0, err
	}
	return t.Unix() - dayStart.Unix(), nil
}

func ParseDatetimeMinute(d string) (int64, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04", d, time.Local)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func DayStart(t time.Time) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", t.Format("2006-01-02"), time.Local)
}

func DayEnd(t time.Time) (time.Time, error) {
	start, err := DayStart(t)
	if err != nil {
		return start, err
	}
	return time.Now().Add(time.Hour * 24), nil
}
