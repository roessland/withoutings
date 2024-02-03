package ptrof

import "time"

func String(s string) *string {
	return &s
}

func Time(t time.Time) *time.Time {
	return &t
}
