package sleep

import (
	"cloud.google.com/go/civil"
	"context"
	"errors"
	"fmt"
	"github.com/roessland/withoutings/pkg/domain/entities"
	"github.com/roessland/withoutings/pkg/logging"
	withingsService "github.com/roessland/withoutings/pkg/withoutings/app/service/withings"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"time"
)

// TODO: convert to query.

type Sleep struct {
	Withings withingsService.Service
}

func NewSleep(withingsSvc withingsService.Service) *Sleep {
	return &Sleep{
		Withings: withingsSvc,
	}
}

type GetSleepSummaryInput struct {
	From time.Time // required
	To   time.Time // defaults to now
}

func (in GetSleepSummaryInput) Validate() error {
	if in.From.IsZero() {
		return errors.New("from cannot be zero")
	}

	if in.To.IsZero() {
		return errors.New("to cannot be zero")
	}

	if in.From.After(in.To) {
		return errors.New("from cannot be after to")
	}

	if in.To.Sub(in.From) > 32*24*time.Hour {
		return errors.New("date range cannot be longer than 1 month")
	}

	return nil
}

type GetSleepSummaryOutput struct {
	Summaries entities.SleepSummaries
	Raw       []byte
}

func (sleepSvc *Sleep) GetSleepSummaries(
	ctx context.Context,
	acc *account.Account,
	in GetSleepSummaryInput,
) (GetSleepSummaryOutput, error) {
	log := logging.MustGetLoggerFromContext(ctx)

	if in.To.IsZero() {
		in.To = time.Now()
	}
	if err := in.Validate(); err != nil {
		return GetSleepSummaryOutput{}, fmt.Errorf("input validation failed: %w", err)
	}

	params := withings.NewSleepGetsummaryParams()
	params.Startdateymd = in.From.Format("2006-01-02")
	params.Enddateymd = in.To.Format("2006-01-02")

	resp, err := sleepSvc.Withings.SleepGetsummary(ctx, acc, params)
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
