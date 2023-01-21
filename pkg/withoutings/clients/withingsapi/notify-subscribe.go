package withingsapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// https://developer.withings.com/api-reference#operation/notify-subscribe

// NotifySubscribeParams are the parameters for Notify - Subscribe
type NotifySubscribeParams struct {
	Action      string `json:"action" url:"action"`
	Callbackurl int64  `json:"startdate" url:"startdate"`
	Appli       int64  `json:"enddate" url:"enddate"`
	Comment     string `json:"data_fields" url:"data_fields"`
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
		return nil, err
	}

	err = json.Unmarshal(resp.Raw, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
