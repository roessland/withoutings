package withings

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"io"
	"net/http"
)

func (c *HTTPClient) NotifyList(ctx context.Context, accessToken string, params withings.NotifyListParams) (*withings.NotifyListResponse, error) {
	return c.WithAccessToken(accessToken).NotifyList(ctx, params)
}

// NewNotifyListRequest creates a new NotifyListRequest request.
func (c *HTTPClient) NewNotifyListRequest(params withings.NotifyListParams) (*http.Request, error) {
	return c.NewRequest("/notify", params)
}

// NotifyList gets all webhook subscriptions for a given notification category.
func (c *AuthenticatedClient) NotifyList(ctx context.Context, params withings.NotifyListParams) (*withings.NotifyListResponse, error) {
	log := logging.MustGetLoggerFromContext(ctx)

	req, err := c.NewNotifyListRequest(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.WithAccessToken(c.AccessToken).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp withings.NotifyListResponse

	resp.Raw, err = io.ReadAll(httpResp.Body)
	if err != nil {
		log.WithField("event", "error.NotifyList.readbody.failed").
			WithError(err).
			Error()
		return nil, err
	}

	err = json.Unmarshal(resp.Raw, &resp)
	if err != nil {
		log.WithField("response_body", string(resp.Raw)).
			WithField("event", "error.NotifyList.unmarshal.failed").
			WithError(err).
			Error()
		return nil, err
	}

	if resp.Status != 0 {
		return &resp, withings.APIError(resp.Status)
	}

	return &resp, nil
}
