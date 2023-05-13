package withings

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"io"
	"net/http"
)

func (c *HTTPClient) SleepGet(ctx context.Context, accessToken string, params withings.SleepGetParams) (*withings.SleepGetResponse, error) {
	return c.WithAccessToken(accessToken).SleepGet(ctx, params)
}

// NewSleepGetRequest creates a new SleepGet request.
func (c *HTTPClient) NewSleepGetRequest(params withings.SleepGetParams) (*http.Request, error) {
	return c.NewRequest("/v2/sleep", params)
}

// SleepGet gets a sleep summary
func (c *AuthenticatedClient) SleepGet(ctx context.Context, params withings.SleepGetParams) (*withings.SleepGetResponse, error) {
	req, err := c.NewSleepGetRequest(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.WithAccessToken(c.AccessToken).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp withings.SleepGetResponse

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