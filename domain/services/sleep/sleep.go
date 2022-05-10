package sleep

import (
	"bytes"
	"cloud.google.com/go/civil"
	"context"
	"fmt"
	"github.com/roessland/withoutings/domain/entities"
	"github.com/roessland/withoutings/logging"
	"github.com/roessland/withoutings/ptrof"
	"github.com/roessland/withoutings/withings"
	openapi2 "github.com/roessland/withoutings/withings/openapi"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
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

type GetSleepSummaryOutput struct {
	Summaries entities.SleepSummaries
	Raw       string
}

func (sleepSvc *Sleep) GetSleepSummaries(ctx context.Context, in GetSleepSummaryInput) (GetSleepSummaryOutput, error) {
	log := logging.MustGetLoggerFromContext(ctx)

	// Get sleep summary
	fields := []string{
		"nb_rem_episodes",
		"sleep_efficiency",
		"sleep_latency",
		"total_sleep_time",
		"total_timeinbed",
		"wakeup_latency",
		"waso",
		"apnea_hypopnea_index",
		"breathing_disturbances_intensity",
		"asleepduration",
		"deepsleepduration",
		"durationtosleep",
		"durationtowakeup",
		"hr_average",
		"hr_max",
		"hr_min",
		"lightsleepduration",
		"night_events",
		"out_of_bed_count",
		"remsleepduration",
		"rr_average",
		"rr_max",
		"rr_min",
		"sleep_score",
		"snoring",
		"snoringepisodecount",
		"wakeupcount",
		"wakeupduration",
	}
	params := &openapi2.Sleepv2GetsummaryParams{
		Startdateymd:  "2022-05-08",
		Enddateymd:    "2022-05-10",
		Lastupdate:    0,
		DataFields:    ptrof.String(strings.Join(fields, ",")),
		Authorization: "Bearer " + in.AccessToken,
	}

	httpResp, err := sleepSvc.Withings.API2.Sleepv2Getsummary(ctx, params)
	if err != nil {
		return GetSleepSummaryOutput{}, fmt.Errorf("fetching data: %w", err)
	}
	defer httpResp.Body.Close()

	f, err := os.Create("last-sleep-resp.json")
	defer f.Close()
	if err != nil {
		return GetSleepSummaryOutput{}, fmt.Errorf("opening file: %w", err)
	}

	bodyBuf, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return GetSleepSummaryOutput{}, fmt.Errorf("reading response: %w", err)
	}
	httpResp.Body = io.NopCloser(bytes.NewReader(bodyBuf))

	_, err = f.Write(bodyBuf)
	if err != nil {
		return GetSleepSummaryOutput{}, fmt.Errorf("writing to file: %w", err)
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

	rawBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return GetSleepSummaryOutput{}, err
	}
	out.Raw = string(rawBytes)

	for _, sleep := range *apiResp.JSON200.Body.Series {
		sleepSummary := entities.SleepSummary{}
		sleepSummary.Date, err = civil.ParseDate(*sleep.Date)
		if err != nil {
			log.WithError(err).WithField("event", "warn.sleepsummary.parsedate").Warn()
		}
		sleepSummary.SleepScore = sleep.Data.SleepScore

		if sleep.Data.TotalSleepTime != nil {
			duration := time.Second * time.Duration(*sleep.Data.TotalSleepTime)
			sleepSummary.TotalSleepTime = &duration
		}
		out.Summaries = append(out.Summaries, sleepSummary)
	}

	return out, nil
}
