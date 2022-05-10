package entities

import (
	"cloud.google.com/go/civil"
	"time"
)

type SleepSummaries []SleepSummary

type SleepSummary struct {
	Date           civil.Date
	SleepScore     *int
	TotalSleepTime *time.Duration
	RawResponse    string
}
