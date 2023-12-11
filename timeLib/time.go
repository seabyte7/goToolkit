package timeLib

import (
	"time"
)

// GetTimeFromTimestampSec Get the time from the timestamp (seconds)
func GetTimeFromTimestampSec(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// GetTimeFromTimestampMs Get the time from the timestamp (milliseconds)
func GetTimeFromTimestampMs(timestamp int64) time.Time {
	return time.Unix(timestamp/SecondMillisecondRatio, (timestamp%SecondMillisecondRatio)*MillisecondNanosecondRatio)
}

// GetTimeZeroTime Get the timestamp (seconds) from the time
func GetTimeZeroTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetTimeLastSecond Get the timestamp (seconds) from the time
func GetTimeLastSecond(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}
