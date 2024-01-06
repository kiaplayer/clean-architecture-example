package helpers

import "time"

func StringToTime(value string) (time.Time, error) {
	return time.ParseInLocation(time.DateTime, value, time.Local)
}

func TimeToString(value time.Time) string {
	return value.Local().Format(time.DateTime)
}
