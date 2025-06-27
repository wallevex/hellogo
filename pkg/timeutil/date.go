package timeutil

import "time"

func GetMonday(t time.Time) time.Time {
	var delta int
	weekday := t.Weekday()
	if weekday != time.Sunday {
		delta = int(time.Monday - weekday)
	} else {
		delta = -6
	}
	return time.Date(t.Year(), t.Month(), t.Day()+delta, 0, 0, 0, 0, time.Local)
}

func GetSunday(t time.Time) time.Time {
	var delta int
	weekday := t.Weekday()
	if weekday == time.Sunday {
		delta = 0
	} else {
		delta = 7 - int(weekday)
	}
	return time.Date(t.Year(), t.Month(), t.Day()+delta, 0, 0, 0, 0, time.Local)
}

func GetDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
