package withings

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

//	params := &openapi2.Sleepv2GetsummaryParams{
//		Startdateymd:  "2022-05-01",
//		Enddateymd:    "2022-05-10",
//		Lastupdate:    0,
//		DataFields:    ptrof.String(strings.Join(fields, ",")),
//		Authorization: "Bearer " + in.AccessToken,
//	}
type SleepGetSummaryParams struct {
	Action        string `json:"action"`
	Startdateymd  string `json:"startdateymd"`
	Enddateymd    string `json:"enddateymd"`
	Lastupdate    int64  `json:"lastupdate"`
	DataFields    string `json:"data_fields"`
	Authorization string `json:"authorization"`
}

type SleepGetsummaryResponse struct {
	Status int `json:"status"`
	Body   SleepGetsummaryBody
	Raw    []byte
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

type SleepGetsummaryData struct {
	ApneaHypopneaIndex             int           `json:"apnea_hypopnea_index"`
	Asleepduration                 int           `json:"asleepduration"`
	BreathingDisturbancesIntensity int           `json:"breathing_disturbances_intensity"`
	Deepsleepduration              int           `json:"deepsleepduration"`
	Durationtosleep                int           `json:"durationtosleep"`
	Durationtowakeup               int           `json:"durationtowakeup"`
	HrAverage                      int           `json:"hr_average"`
	HrMax                          int           `json:"hr_max"`
	HrMin                          int           `json:"hr_min"`
	Lightsleepduration             int           `json:"lightsleepduration"`
	NbRemEpisodes                  int           `json:"nb_rem_episodes"`
	NightEvents                    []interface{} `json:"night_events"`
	OutOfBedCount                  int           `json:"out_of_bed_count"`
	Remsleepduration               int           `json:"remsleepduration"`
	RrAverage                      int           `json:"rr_average"`
	RrMax                          int           `json:"rr_max"`
	RrMin                          int           `json:"rr_min"`
	SleepEfficiency                int           `json:"sleep_efficiency"`
	SleepLatency                   int           `json:"sleep_latency"`
	SleepScore                     int           `json:"sleep_score"`
	Snoring                        int           `json:"snoring"`
	Snoringepisodecount            int           `json:"snoringepisodecount"`
	TotalSleepTime                 int           `json:"total_sleep_time"`
	TotalTimeinbed                 int           `json:"total_timeinbed"`
	WakeupLatency                  int           `json:"wakeup_latency"`
	Wakeupcount                    int           `json:"wakeupcount"`
	Wakeupduration                 int           `json:"wakeupduration"`
	Waso                           int           `json:"waso"`
}

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

// NewSleepGetsummaryParams creates new SleepGetSummaryParams with some defaults.
func NewSleepGetsummaryParams() SleepGetSummaryParams {
	return SleepGetSummaryParams{
		Action:     "getsummary",
		DataFields: SleepGetSummaryAllDataFields,
	}
}

// NewSleepGetsummaryRequest creates a new SleepGetsummary request.
func (c *Client) NewSleepGetsummaryRequest(params SleepGetSummaryParams) (*http.Request, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, c.APIBase+"/v2/sleep", io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}
	return req, nil
}

// SleepGetsummary gets a sleep summary
func (c *AuthenticatedClient) SleepGetsummary(ctx context.Context, params SleepGetSummaryParams) (*SleepGetsummaryResponse, error) {
	req, err := c.NewSleepGetsummaryRequest(params)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.WithAccessToken(c.AccessToken).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp SleepGetsummaryResponse

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
