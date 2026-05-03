package withings

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"io"
	"net/http"
)

func (c *HTTPClient) MeasureGetworkouts(ctx context.Context, accessToken string, params withings.MeasureGetworkoutsParams) (*withings.MeasureGetworkoutsResponse, error) {
	return c.WithAccessToken(accessToken).MeasureGetworkouts(ctx, params)
}

// NewMeasureGetworkoutsRequest creates a new MeasureGetworkouts request.
func (c *HTTPClient) NewMeasureGetworkoutsRequest(params withings.MeasureGetworkoutsParams) (*http.Request, error) {
	return c.NewRequest("/v2/measure", params, nil)
}

// MeasureGetworkouts gets workout summaries
func (c *AuthenticatedClient) MeasureGetworkouts(ctx context.Context, params withings.MeasureGetworkoutsParams) (*withings.MeasureGetworkoutsResponse, error) {
	log := logging.MustGetLoggerFromContext(ctx)
	req, err := c.NewMeasureGetworkoutsRequest(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.WithAccessToken(c.AccessToken).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp withings.MeasureGetworkoutsResponse

	resp.Raw, err = io.ReadAll(httpResp.Body)
	if err != nil {
		log.WithField("event", "error.MeasureGetworkouts.readbody.failed").
			WithError(err).
			Error()
		return nil, err
	}

	err = json.Unmarshal(resp.Raw, &resp)
	if err != nil {
		log.WithField("response_body", string(resp.Raw)).
			WithField("event", "error.MeasureGetworkouts.unmarshal.failed").
			WithError(err).
			Error()
		return nil, err
	}

	if resp.Status != 0 {
		return &resp, withings.APIError(resp.Status)
	}

	return &resp, nil
}
