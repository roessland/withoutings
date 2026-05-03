package withings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// https://developer.withings.com/api-reference#operation/measurev2-getworkouts

var MeasureGetworkoutsAllDataFields = strings.Join([]string{
	"calories",
	"intensity",
	"manual_intensity",
	"manual_distance",
	"manual_calories",
	"hr_average",
	"hr_min",
	"hr_max",
	"hr_zone_0",
	"hr_zone_1",
	"hr_zone_2",
	"hr_zone_3",
	"pause_duration",
	"algo_pause_duration",
	"spo2_average",
	"core_body_temperature_avg",
	"core_body_temperature_min",
	"core_body_temperature_max",
	"core_body_temperature_status",
	"steps",
	"distance",
	"elevation",
	"pool_laps",
	"strokes",
	"pool_length",
}, ",")

// NewMeasureGetworkoutsParams creates new MeasureGetworkoutsParams with some defaults.
func NewMeasureGetworkoutsParams() MeasureGetworkoutsParams {
	return MeasureGetworkoutsParams{
		Action:     "getworkouts",
		DataFields: MeasureGetworkoutsAllDataFields,
	}
}

// MeasureGetworkoutsParams
// Don't set Lastupdate and Startdateymd or Enddateymd at the same time.
type MeasureGetworkoutsParams struct {
	Action       string `json:"action" url:"action"`
	Startdateymd string `json:"startdateymd,omitempty" url:"startdateymd,omitempty"`
	Enddateymd   string `json:"enddateymd,omitempty" url:"enddateymd,omitempty"`
	Lastupdate   int64  `json:"lastupdate,omitempty" url:"lastupdate,omitempty"`
	Offset       int    `json:"offset,omitempty" url:"offset,omitempty"`
	DataFields   string `json:"data_fields" url:"data_fields"`
}

type MeasureGetworkoutsResponse struct {
	Status int                    `json:"status"`
	Body   MeasureGetworkoutsBody `json:"body"`
	Raw    []byte                 `json:"-"`
}

func MustNewMeasureGetworkoutsResponse(raw []byte) *MeasureGetworkoutsResponse {
	var resp MeasureGetworkoutsResponse
	err := json.Unmarshal(raw, &resp)
	resp.Raw = raw

	if err != nil {
		panic(fmt.Errorf(`couldn't unmarshal MeasureGetworkoutsResponse: %w`, err))
	}

	return &resp
}

type MeasureGetworkoutsBody struct {
	Series []json.RawMessage `json:"series"`
	More   bool              `json:"more"`
	Offset int               `json:"offset"`
}
