package timeLib

import "time"

// GetCurrentMs Gets the current time in milliseconds
func GetCurrentMs() int64 {
	return time.Now().UnixNano() / MillisecondNanosecondRatio
}

// GetCurrentSec Gets the current time in seconds
func GetCurrentSec() int64 {
	return time.Now().Unix()
}
