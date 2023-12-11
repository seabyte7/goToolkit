package timeLib

import "time"

// GetMonthFirstDayZeroTime Get the 0 o 'clock time of the first day of the month
func GetMonthFirstDayZeroTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// GetMonthLastDayLastSecond Get the last second time on the last day of the month
func GetMonthLastDayLastSecond(t time.Time) time.Time {
	return GetTimeLastSecond(GetMonthFirstDayZeroTime(t).AddDate(0, 1, -1))
}
