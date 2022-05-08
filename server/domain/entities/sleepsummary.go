package entities

import "cloud.google.com/go/civil"

type SleepSummaries []SleepSummary

type SleepSummary struct {
	Date       civil.Date
	SleepScore *int
}
