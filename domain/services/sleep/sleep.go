package sleep

import (
	"cloud.google.com/go/civil"
	"context"
	"github.com/roessland/withoutings/domain/entities"
	"github.com/roessland/withoutings/logging"
	"github.com/roessland/withoutings/withingsapi"
	"time"
)

type Sleep struct {
	Withings *withingsapi.Client
}

func NewSleep(withings *withingsapi.Client) *Sleep {
	return &Sleep{
		Withings: withings,
	}
}

type GetSleepSummaryInput struct {
	Year  int
	Month int
}

type GetSleepSummaryOutput struct {
	Summaries entities.SleepSummaries
	Raw       []byte
}

func (sleepSvc *Sleep) GetSleepSummaries(ctx context.Context, accessToken string, in GetSleepSummaryInput) (GetSleepSummaryOutput, error) {
	log := logging.MustGetLoggerFromContext(ctx)
	authClient := sleepSvc.Withings.WithAccessToken(accessToken)

	// todo: use input
	params := withingsapi.NewSleepGetsummaryParams()
	params.Startdateymd = "2022-11-01"
	params.Enddateymd = "2022-11-12"

	resp, err := authClient.SleepGetsummary(ctx, params)
	if err != nil {
		return GetSleepSummaryOutput{}, err
	}

	// Convert to domain entities
	out := GetSleepSummaryOutput{}

	out.Raw = resp.Raw

	for i := range resp.Body.Series {
		sleep := resp.Body.Series[i]
		sleepSummary := entities.SleepSummary{}
		sleepSummary.Date, err = civil.ParseDate(sleep.Date)
		if err != nil {
			log.WithError(err).WithField("event", "warn.sleepsummary.parsedate").Warn()
		}
		sleepSummary.SleepScore = &sleep.Data.SleepScore
		duration := time.Second * time.Duration(sleep.Data.TotalSleepTime)
		sleepSummary.TotalSleepTime = &duration
		out.Summaries = append(out.Summaries, sleepSummary)
	}

	return out, nil
}
