package common

import "time"

func GetCurrentTimeStamp() *time.Time {
	timestamp := time.Now().UTC()
	return &timestamp
}
