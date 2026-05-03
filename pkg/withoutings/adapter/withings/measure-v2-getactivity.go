package withings

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"io"
	"net/http"
)

func (c *HTTPClient) MeasureGetactivity(ctx context.Context, accessToken string, params withings.MeasureGetactivityParams) (*withings.MeasureGetactivityResponse, error) {
	return c.WithAccessToken(accessToken).MeasureGetactivity(ctx, params)
}

// NewMeasureGetactivityRequest creates a new MeasureGetactivity request.
func (c *HTTPClient) NewMeasureGetactivityRequest(params withings.MeasureGetactivityParams) (*http.Request, error) {
	return c.NewRequest("/v2/measure", params, nil)
}

// MeasureGetactivity gets daily aggregated activity data
func (c *AuthenticatedClient) MeasureGetactivity(ctx context.Context, params withings.MeasureGetactivityParams) (*withings.MeasureGetactivityResponse, error) {
	log := logging.MustGetLoggerFromContext(ctx)
	req, err := c.NewMeasureGetactivityRequest(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.WithAccessToken(c.AccessToken).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp withings.MeasureGetactivityResponse

	resp.Raw, err = io.ReadAll(httpResp.Body)
	if err != nil {
		log.WithField("event", "error.MeasureGetactivity.readbody.failed").
			WithError(err).
			Error()
		return nil, err
	}

	err = json.Unmarshal(resp.Raw, &resp)
	if err != nil {
		log.WithField("response_body", string(resp.Raw)).
			WithField("event", "error.MeasureGetactivity.unmarshal.failed").
			WithError(err).
			Error()
		return nil, err
	}

	if resp.Status != 0 {
		return &resp, withings.APIError(resp.Status)
	}

	return &resp, nil
}
