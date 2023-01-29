package withingsapi

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/logging"
	"io"
	"net/http"
)

// https://developer.withings.com/api-reference#operation/notify-subscribe

// NotifySubscribeParams are the parameters for Notify - Subscribe
type NotifySubscribeParams struct {
	Action      string `json:"action" url:"action"`
	Callbackurl string `json:"callbackurl" url:"callbackurl"`
	Appli       int    `json:"appli" url:"appli"`
	Comment     string `json:"comment" url:"comment"`
}

type NotifySubscribeResponse struct {
	Status int    `json:"status"`
	Raw    []byte `json:"-"`
}

// NewNotifySubscribeParams creates new NewNotifySubscribeParams with some defaults.
func NewNotifySubscribeParams() NotifySubscribeParams {
	return NotifySubscribeParams{
		Action: "subscribe",
	}
}

// NewNotifySubscribeRequest creates a new NotifySubscribeRequest request.
func (c *Client) NewNotifySubscribeRequest(params NotifySubscribeParams) (*http.Request, error) {
	return c.NewRequest("/notify", params)
}

// NotifySubscribe subscribes to health data events for the current user.
func (c *AuthenticatedClient) NotifySubscribe(ctx context.Context, params NotifySubscribeParams) (*NotifySubscribeResponse, error) {
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

	var resp NotifySubscribeResponse

	resp.Raw, err = io.ReadAll(httpResp.Body)
	if err != nil {
		log.WithField("event", "NotifySubscribe io.ReadAll failed").
			WithError(err).
			Error()
		return nil, err
	}

	log.WithField("event", "NotifySubscribe response read").
		WithField("response_body", string(resp.Raw)).
		WithField("http_status_code", httpResp.StatusCode).
		Info()

	err = json.Unmarshal(resp.Raw, &resp)
	if err != nil {
		log.WithField("response_body", string(resp.Raw)).
			WithField("event", "NotifySubscribe io.ReadAll failed").
			WithError(err).
			Error()
		return nil, err
	}

	if resp.Status != 0 {
		return &resp, APIError(resp.Status)
	}

	return &resp, nil
}
