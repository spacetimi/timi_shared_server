package time_utils

import "time"

func TimeToUnixTimeStamp(t time.Time) int64 {
	return t.Unix()
}

func UnixTimeStampToLocalTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func DurationToSeconds(d time.Duration) int64 {
	return int64(d.Seconds())
}

func DurationToDays(d time.Duration) int64 {
	return DurationToSeconds(d) / (3600 * 24)
}

func GetDurationBetweenTimes(first time.Time, second time.Time) time.Duration {
	return second.Sub(first)
}

func GetLocalYear() int {
	return time.Now().Year()
}

func GetLocalMonth() int {
	return int(time.Now().Month())
}

func GetLocalDayOfMonth() int {
	return time.Now().Day()
}

func BeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
