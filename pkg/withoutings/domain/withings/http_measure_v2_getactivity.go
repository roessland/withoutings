package withings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// https://developer.withings.com/api-reference#operation/measurev2-getactivity

var MeasureGetactivityAllDataFields = strings.Join([]string{
	"steps",
	"distance",
	"elevation",
	"soft",
	"moderate",
	"intense",
	"active",
	"calories",
	"totalcalories",
	"hr_average",
	"hr_min",
	"hr_max",
	"hr_zone_0",
	"hr_zone_1",
	"hr_zone_2",
	"hr_zone_3",
}, ",")

// NewMeasureGetactivityParams creates new MeasureGetactivityParams with some defaults.
func NewMeasureGetactivityParams() MeasureGetactivityParams {
	return MeasureGetactivityParams{
		Action:     "getactivity",
		DataFields: MeasureGetactivityAllDataFields,
	}
}

// MeasureGetactivityParams
// Don't set Lastupdate and Startdateymd or Enddateymd at the same time.
type MeasureGetactivityParams struct {
	Action       string `json:"action" url:"action"`
	Startdateymd string `json:"startdateymd,omitempty" url:"startdateymd,omitempty"`
	Enddateymd   string `json:"enddateymd,omitempty" url:"enddateymd,omitempty"`
	Lastupdate   int64  `json:"lastupdate,omitempty" url:"lastupdate,omitempty"`
	Offset       int    `json:"offset,omitempty" url:"offset,omitempty"`
	DataFields   string `json:"data_fields" url:"data_fields"`
}

type MeasureGetactivityResponse struct {
	Status int                    `json:"status"`
	Body   MeasureGetactivityBody `json:"body"`
	Raw    []byte                 `json:"-"`
}

func MustNewMeasureGetactivityResponse(raw []byte) *MeasureGetactivityResponse {
	var resp MeasureGetactivityResponse
	err := json.Unmarshal(raw, &resp)
	resp.Raw = raw

	if err != nil {
		panic(fmt.Errorf(`couldn't unmarshal MeasureGetactivityResponse: %w`, err))
	}

	return &resp
}

type MeasureGetactivityBody struct {
	Activities []json.RawMessage `json:"activities"`
	More       bool              `json:"more"`
	Offset     int               `json:"offset"`
}
