package timeLib

import "time"

// GetTodayZeroTime Get today's 0 o 'clock time
func GetTodayZeroTime() time.Time {
	// 获取当前时间
	now := time.Now()
	// 获取当前时间的年月日
	year, month, day := now.Date()
	// 获取今天0点的时间
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location())
}

// GetTodayLastSecond Get the last second of today
func GetTodayLeftSec() int64 {
	return GetTimeLastSecond(time.Now()).Unix() - time.Now().Unix()
}

// GetTodayPastSec Get the last second of today
func GetTodayPastSec() int64 {
	return time.Now().Unix() - GetTimeZeroTime(time.Now()).Unix()
}

// GetTodayLastSecond Get the last second of today
func GetDayStartAndEndTimeSec(t time.Time) (startTimeSec, endTimeSec int64) {
	startTime := GetTimeZeroTime(t)
	endTime := startTime.AddDate(0, 0, 1)
	startTimeSec = startTime.Unix()
	endTimeSec = endTime.Unix()

	return
}
