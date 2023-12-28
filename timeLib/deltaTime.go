package timeLib

import "time"

// GetDeltaTimeSec Get the difference between the specified time and the current time
func GetDeltaTimeSec(t time.Time) int64 {
	return GetCurrentSec() - t.Unix()
}

// GetDeltaTimeMs Get the difference between the specified time and the current time
func GetDeltaTimeMs(t time.Time) int64 {
	return GetCurrentMs() - t.UnixNano()/MillisecondNanosecondRatio
}
