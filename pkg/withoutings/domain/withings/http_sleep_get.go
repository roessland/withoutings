package withings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// https://developer.withings.com/api-reference#operation/sleepv2-Get

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

func MustNewSleepGetResponse(raw []byte) *SleepGetResponse {
	var resp SleepGetResponse
	err := json.Unmarshal(raw, &resp)
	resp.Raw = raw

	if err != nil {
		panic(fmt.Errorf(`couldn't unmarshal SleepGetResponse: %w`, err))
	}

	return &resp
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
