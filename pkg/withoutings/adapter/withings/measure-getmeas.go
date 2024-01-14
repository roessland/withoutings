package withings

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"io"
	"net/http"
)

func (c *HTTPClient) MeasureGetmeas(ctx context.Context, accessToken string, params withings.MeasureGetmeasParams) (*withings.MeasureGetmeasResponse, error) {
	return c.WithAccessToken(accessToken).MeasureGetmeas(ctx, params)
}

// NewMeasureGetmeasRequest creates a new MeasureGetmeas request.
func (c *HTTPClient) NewMeasureGetmeasRequest(params withings.MeasureGetmeasParams) (*http.Request, error) {
	return c.NewRequest("/measure", params)
}

// MeasureGetmeas gets measures
func (c *AuthenticatedClient) MeasureGetmeas(ctx context.Context, params withings.MeasureGetmeasParams) (*withings.MeasureGetmeasResponse, error) {
	log := logging.MustGetLoggerFromContext(ctx)
	req, err := c.NewMeasureGetmeasRequest(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.WithAccessToken(c.AccessToken).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp withings.MeasureGetmeasResponse

	resp.Raw, err = io.ReadAll(httpResp.Body)
	if err != nil {
		log.WithField("event", "error.MeasureGetmeas.readbody.failed").
			WithError(err).
			Error()
		return nil, err
	}

	err = json.Unmarshal(resp.Raw, &resp)
	if err != nil {
		log.WithField("response_body", string(resp.Raw)).
			WithField("event", "error.MeasureGetmeas.unmarshal.failed").
			WithError(err).
			Error()
		return nil, err
	}

	if resp.Status != 0 {
		return &resp, withings.APIError(resp.Status)
	}

	return &resp, nil
}
