package withings

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"io"
	"net/http"
)

func (c *HTTPClient) NotifySubscribe(ctx context.Context, accessToken string, params withings.NotifySubscribeParams) (*withings.NotifySubscribeResponse, error) {
	return c.WithAccessToken(accessToken).NotifySubscribe(ctx, params)
}

// NewNotifySubscribeRequest creates a new NotifySubscribeRequest request.
func (c *HTTPClient) NewNotifySubscribeRequest(params withings.NotifySubscribeParams) (*http.Request, error) {
	return c.NewRequest("/notify", params)
}

// NotifySubscribe subscribes to health data events for the current user.
func (c *AuthenticatedClient) NotifySubscribe(ctx context.Context, params withings.NotifySubscribeParams) (*withings.NotifySubscribeResponse, error) {
	log := logging.MustGetLoggerFromContext(ctx)

	req, err := c.NewNotifySubscribeRequest(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.WithAccessToken(c.AccessToken).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp withings.NotifySubscribeResponse

	resp.Raw, err = io.ReadAll(httpResp.Body)
	if err != nil {
		log.WithField("event", "error.NotifySubscribe.readbody.failed").
			WithError(err).
			Error()
		return nil, err
	}

	err = json.Unmarshal(resp.Raw, &resp)
	if err != nil {
		log.WithField("response_body", string(resp.Raw)).
			WithField("event", "error.NotifySubscribe.unmarshal.failed").
			WithError(err).
			Error()
		return nil, err
	}

	if resp.Status != 0 {
		return &resp, withings.APIError(resp.Status)
	}

	return &resp, nil
}
