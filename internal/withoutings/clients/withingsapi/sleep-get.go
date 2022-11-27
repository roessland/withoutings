package withingsapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// https://developer.withings.com/api-reference#operation/sleepv2-Get

// SleepGetParams are the parameters for Sleep v2 - Get
type SleepGetParams struct {
	Action     string `json:"action" url:"action"`
	Startdate  int64  `json:"startdate" url:"startdate"`
	Enddate    int64  `json:"enddate" url:"enddate"`
	DataFields string `json:"data_fields" url:"data_fields"`
}

type SleepGetResponse struct {
	Status int          `json:"status"`
	Body   SleepGetBody `json:"body"`
	Raw    []byte       `json:"-"`
}

type SleepGetBody struct {
	Series []SleepGetEntry `json:"series"`
}

type SleepGetEntry struct {
	Startdate int64           `json:"startdate"`
	Enddate   int64           `json:"enddate"`
	State     int             `json:"state"`
	Model     string          `json:"model"`
	ModelID   int             `json:"model_id"`
	HR        json.RawMessage `json:"hr"`
	RR        json.RawMessage `json:"rr"`
	Snoring   json.RawMessage `json:"snoring"`
	SDNN1     json.RawMessage `json:"sdnn_1"`
	RMSSD     json.RawMessage `json:"rmssd"`
}

var SleepGetAllDataFields = strings.Join([]string{
	"hr",
	"rr",
	"snoring",
	"sdnn_1",
	"rmssd",
}, ",")

// NewSleepGetParams creates new SleepGetParams with some defaults.
func NewSleepGetParams() SleepGetParams {
	return SleepGetParams{
		Action:     "get",
		DataFields: SleepGetAllDataFields,
	}
}

// NewSleepGetRequest creates a new SleepGet request.
func (c *Client) NewSleepGetRequest(params SleepGetParams) (*http.Request, error) {
	return c.NewRequest("/v2/sleep", params)
}

// SleepGet gets a sleep summary
func (c *AuthenticatedClient) SleepGet(ctx context.Context, params SleepGetParams) (*SleepGetResponse, error) {
	req, err := c.NewSleepGetRequest(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.WithAccessToken(c.AccessToken).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp SleepGetResponse

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
