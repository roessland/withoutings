package withings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// https://developer.withings.com/api-reference#operation/measurev2-getintradayactivity

var MeasureGetintradayactivityAllDataFields = strings.Join([]string{
	"steps",
	"elevation",
	"calories",
	"distance",
	"stroke",
	"pool_lap",
	"duration",
	"heart_rate",
	"spo2_auto",
	"rmssd",
	"sdnn1",
	"hrv_quality",
	"core_body_temperature",
	"rr",
	"chest_movement_rate",
}, ",")

// NewMeasureGetintradayactivityParams creates new MeasureGetintradayactivityParams with some defaults.
func NewMeasureGetintradayactivityParams() MeasureGetintradayactivityParams {
	return MeasureGetintradayactivityParams{
		Action:     "getintradayactivity",
		DataFields: MeasureGetintradayactivityAllDataFields,
	}
}

// MeasureGetintradayactivityParams are the parameters for Measure v2 - Getintradayactivity.
type MeasureGetintradayactivityParams struct {
	Action     string `json:"action" url:"action"`
	Startdate  int64  `json:"startdate,omitempty" url:"startdate,omitempty"`
	Enddate    int64  `json:"enddate,omitempty" url:"enddate,omitempty"`
	DataFields string `json:"data_fields" url:"data_fields"`
}

type MeasureGetintradayactivityResponse struct {
	Status int                              `json:"status"`
	Body   MeasureGetintradayactivityBody   `json:"body"`
	Raw    []byte                           `json:"-"`
}

func MustNewMeasureGetintradayactivityResponse(raw []byte) *MeasureGetintradayactivityResponse {
	var resp MeasureGetintradayactivityResponse
	err := json.Unmarshal(raw, &resp)
	resp.Raw = raw

	if err != nil {
		panic(fmt.Errorf(`couldn't unmarshal MeasureGetintradayactivityResponse: %w`, err))
	}

	return &resp
}

// MeasureGetintradayactivityBody preserves the raw "series" object — the API returns
// an object keyed by unix timestamp (or sometimes {} when there is no data),
// not an array.
type MeasureGetintradayactivityBody struct {
	Series json.RawMessage `json:"series"`
}
