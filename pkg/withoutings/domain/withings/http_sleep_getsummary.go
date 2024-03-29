package withings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// https://developer.withings.com/api-reference#operation/sleepv2-getsummary

var SleepGetSummaryAllDataFields = strings.Join([]string{
	"nb_rem_episodes",
	"sleep_efficiency",
	"sleep_latency",
	"total_sleep_time",
	"total_timeinbed",
	"wakeup_latency",
	"waso",
	"apnea_hypopnea_index",
	"breathing_disturbances_intensity",
	"asleepduration",
	"deepsleepduration",
	"durationtosleep",
	"durationtowakeup",
	"hr_average",
	"hr_max",
	"hr_min",
	"lightsleepduration",
	"night_events",
	"out_of_bed_count",
	"remsleepduration",
	"rr_average",
	"rr_max",
	"rr_min",
	"sleep_score",
	"snoring",
	"snoringepisodecount",
	"wakeupcount",
	"wakeupduration",
}, ",")

// NewSleepGetsummaryParams creates new SleepGetsummaryParams with some defaults.
func NewSleepGetsummaryParams() SleepGetsummaryParams {
	return SleepGetsummaryParams{
		Action:     "getsummary",
		DataFields: SleepGetSummaryAllDataFields,
	}
}

// SleepGetsummaryParams
// Don't set Lastupdate and Startdateymd or Enddateymd at the same time.
type SleepGetsummaryParams struct {
	Action       string `json:"action" url:"action"`
	Startdateymd string `json:"startdateymd" url:"startdateymd"`
	Enddateymd   string `json:"enddateymd" url:"enddateymd"`
	Lastupdate   int64  `json:"lastupdate,omitempty" url:"lastupdate,omitempty"`
	DataFields   string `json:"data_fields" url:"data_fields"`
}

type SleepGetsummaryResponse struct {
	Status int `json:"status"`
	Body   SleepGetsummaryBody
	Raw    []byte
}

func MustNewSleepGetsummaryResponse(raw []byte) *SleepGetsummaryResponse {
	var resp SleepGetsummaryResponse
	err := json.Unmarshal(raw, &resp)
	resp.Raw = raw

	if err != nil {
		panic(fmt.Errorf(`couldn't unmarshal SleepGetsummaryResponse: %w`, err))
	}

	return &resp
}

type SleepGetsummaryBody struct {
	Series []SleepGetsummaryEntry `json:"series"`
	More   bool                   `json:"more"`
	Offset int                    `json:"offset"`
}

type SleepGetsummaryEntry struct {
	Timezone  string `json:"timezone"`
	Model     int    `json:"model"`
	ModelId   int    `json:"model_id"`
	Startdate int    `json:"startdate"`
	Enddate   int    `json:"enddate"`
	Date      string `json:"date"`
	Created   int    `json:"created"`
	Modified  int    `json:"modified"`
	Data      SleepGetsummaryData
}

// SleepGetsummaryData is the response from the Sleep v2 - Getsummary API
// https://developer.withings.com/api-reference#operation/sleepv2-getsummary
type SleepGetsummaryData struct {
	ApneaHypopneaIndex             float64         `json:"apnea_hypopnea_index"`
	Asleepduration                 float64         `json:"asleepduration"`
	BreathingDisturbancesIntensity float64         `json:"breathing_disturbances_intensity"`
	Deepsleepduration              float64         `json:"deepsleepduration"` // deprecated
	Durationtosleep                float64         `json:"durationtosleep"`   // deprecated
	Durationtowakeup               float64         `json:"durationtowakeup"`  // deprecated
	HrAverage                      float64         `json:"hr_average"`
	HrMax                          float64         `json:"hr_max"`
	HrMin                          float64         `json:"hr_min"`
	Lightsleepduration             float64         `json:"lightsleepduration"`
	NbRemEpisodes                  int             `json:"nb_rem_episodes"`
	NightEvents                    json.RawMessage `json:"night_events"`
	OutOfBedCount                  int             `json:"out_of_bed_count"`
	Remsleepduration               float64         `json:"remsleepduration"`
	RrAverage                      float64         `json:"rr_average"`
	RrMax                          float64         `json:"rr_max"`
	RrMin                          float64         `json:"rr_min"`
	SleepEfficiency                float64         `json:"sleep_efficiency"`
	SleepLatency                   float64         `json:"sleep_latency"`
	SleepScore                     float64         `json:"sleep_score"`
	Snoring                        float64         `json:"snoring"`
	Snoringepisodecount            int             `json:"snoringepisodecount"`
	TotalSleepTime                 float64         `json:"total_sleep_time"`
	TotalTimeinbed                 int             `json:"total_timeinbed"`
	WakeupLatency                  float64         `json:"wakeup_latency"`
	Wakeupcount                    int             `json:"wakeupcount"`
	Wakeupduration                 float64         `json:"wakeupduration"`
	Waso                           float64         `json:"waso"`
}
