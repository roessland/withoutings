package services

import (
	"cloud.google.com/go/civil"
	"context"
	"fmt"
	"github.com/roessland/withoutings/logging"
	"github.com/roessland/withoutings/ptrof"
	"github.com/roessland/withoutings/server/domain/entities"
	"github.com/roessland/withoutings/withings"
	"github.com/roessland/withoutings/withingsapi2/openapi2"
)

type Sleep struct {
	Withings *withings.Client
}

func NewSleep(withings *withings.Client) *Sleep {
	return &Sleep{
		Withings: withings,
	}
}

type GetSleepSummaryInput struct {
	AccessToken string
	Year        int
	Month       int
}

type GetSleepSummaryOutput entities.SleepSummaries

func (sleepSvc *Sleep) GetSleepSummaries(ctx context.Context, in GetSleepSummaryInput) (GetSleepSummaryOutput, error) {
	log := logging.MustGetLoggerFromContext(ctx)

	// Get sleep summary
	params := &openapi2.Sleepv2GetsummaryParams{
		Startdateymd:  "2022-05-01",
		Enddateymd:    "2022-05-08",
		Lastupdate:    0,
		DataFields:    ptrof.String("total_sleep_time,sleep_score"),
		Authorization: "Bearer " + in.AccessToken,
	}

	httpResp, err := sleepSvc.Withings.API2.Sleepv2Getsummary(ctx, params)
	if err != nil {
		return GetSleepSummaryOutput{}, fmt.Errorf("fetching data: %w", err)
	}

	// Check response code
	err = openapi2.ParseErrorResponse(httpResp)
	if err != nil {
		return GetSleepSummaryOutput{}, fmt.Errorf("non-zero response code: %w", err)
	}

	// Decode sleep summary into struct
	apiResp, err := openapi2.ParseSleepv2GetsummaryResponse(httpResp)
	if err != nil {
		return GetSleepSummaryOutput{}, fmt.Errorf("parsing sleep response: %w", err)
	}

	// Convert to domain entities
	out := GetSleepSummaryOutput{}
	for _, sleep := range *apiResp.JSON200.Body.Series {
		sleepSummary := entities.SleepSummary{}
		sleepSummary.Date, err = civil.ParseDate(*sleep.Date)
		if err != nil {
			log.WithError(err).WithField("event", "warn.sleepsummary.parsedate").Warn()
		}
		sleepSummary.SleepScore = sleep.Data.SleepScore
		out = append(out, sleepSummary)
	}

	return out, nil
}
