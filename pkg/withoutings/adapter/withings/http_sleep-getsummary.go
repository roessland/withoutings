package withings

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"io"
	"net/http"
)

func (c *HTTPClient) SleepGetsummary(ctx context.Context, accessToken string, params withings.SleepGetSummaryParams) (*withings.SleepGetsummaryResponse, error) {
	return c.WithAccessToken(accessToken).SleepGetsummary(ctx, params)
}

// NewSleepGetsummaryRequest creates a new SleepGetsummary request.
func (c *HTTPClient) NewSleepGetsummaryRequest(params withings.SleepGetSummaryParams) (*http.Request, error) {
	return c.NewRequest("/v2/sleep", params)
}

// SleepGetsummary gets a sleep summary
func (c *AuthenticatedClient) SleepGetsummary(ctx context.Context, params withings.SleepGetSummaryParams) (*withings.SleepGetsummaryResponse, error) {
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
		return nil, err
	}

	err = json.Unmarshal(resp.Raw, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
