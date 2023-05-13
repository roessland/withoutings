package withings

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"io"
	"net/http"
)

func (c *HTTPClient) SleepGetsummary(ctx context.Context, accessToken string, params withings.SleepGetsummaryParams) (*withings.SleepGetsummaryResponse, error) {
	return c.WithAccessToken(accessToken).SleepGetsummary(ctx, params)
}

// NewSleepGetsummaryRequest creates a new SleepGetsummary request.
func (c *HTTPClient) NewSleepGetsummaryRequest(params withings.SleepGetsummaryParams) (*http.Request, error) {
	return c.NewRequest("/v2/sleep", params)
}

// SleepGetsummary gets a sleep summary
// TODO test that it returns error for non-0 Status.
func (c *AuthenticatedClient) SleepGetsummary(ctx context.Context, params withings.SleepGetsummaryParams) (*withings.SleepGetsummaryResponse, error) {
	log := logging.MustGetLoggerFromContext(ctx)

	req, err := c.NewSleepGetsummaryRequest(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.WithAccessToken(c.AccessToken).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp withings.SleepGetsummaryResponse

	resp.Raw, err = io.ReadAll(httpResp.Body)
	if err != nil {
		log.WithField("event", "SleepGetsummary io.ReadAll failed").
			WithError(err).
			Error()
		return nil, err
	}

	err = json.Unmarshal(resp.Raw, &resp)
	if err != nil {
		log.WithField("response_body", string(resp.Raw)).
			WithField("event", "SleepGetsummary io.ReadAll failed").
			WithError(err).
			Error()
		return nil, err
	}

	if resp.Status != 0 {
		return &resp, withings.APIError(resp.Status)
	}

	return &resp, nil
}
