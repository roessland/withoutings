package withings

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"io"
	"net/http"
)

func (c *HTTPClient) MeasureGetintradayactivity(ctx context.Context, accessToken string, params withings.MeasureGetintradayactivityParams) (*withings.MeasureGetintradayactivityResponse, error) {
	return c.WithAccessToken(accessToken).MeasureGetintradayactivity(ctx, params)
}

// NewMeasureGetintradayactivityRequest creates a new MeasureGetintradayactivity request.
func (c *HTTPClient) NewMeasureGetintradayactivityRequest(params withings.MeasureGetintradayactivityParams) (*http.Request, error) {
	return c.NewRequest("/v2/measure", params, nil)
}

// MeasureGetintradayactivity gets activity data captured at high frequency
func (c *AuthenticatedClient) MeasureGetintradayactivity(ctx context.Context, params withings.MeasureGetintradayactivityParams) (*withings.MeasureGetintradayactivityResponse, error) {
	log := logging.MustGetLoggerFromContext(ctx)
	req, err := c.NewMeasureGetintradayactivityRequest(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.WithAccessToken(c.AccessToken).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp withings.MeasureGetintradayactivityResponse

	resp.Raw, err = io.ReadAll(httpResp.Body)
	if err != nil {
		log.WithField("event", "error.MeasureGetintradayactivity.readbody.failed").
			WithError(err).
			Error()
		return nil, err
	}

	err = json.Unmarshal(resp.Raw, &resp)
	if err != nil {
		log.WithField("response_body", string(resp.Raw)).
			WithField("event", "error.MeasureGetintradayactivity.unmarshal.failed").
			WithError(err).
			Error()
		return nil, err
	}

	if resp.Status != 0 {
		return &resp, withings.APIError(resp.Status)
	}

	return &resp, nil
}
