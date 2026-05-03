package port

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
)

// SleepSession renders detailed visualizations for a single sleep session.
// The session is identified by its Withings `startdate` (Unix seconds).
// Data is read from notification_data populated by appli=44 webhooks; if no
// matching session is stored for the user, a "no data" message is shown.
//
// Methods: GET
func SleepSession(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		acc := account.GetFromContext(ctx)
		if acc == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "You must be logged in to view this page.")
			return
		}

		startdate, err := strconv.ParseInt(mux.Vars(r)["startdate"], 10, 64)
		if err != nil {
			http.Error(w, "invalid startdate", http.StatusBadRequest)
			return
		}

		view, err := buildSleepSessionView(ctx, svc.SubscriptionRepo, acc.UUID(), startdate)
		if err != nil {
			log.WithError(err).WithField("event", "error.SleepSession.build").Error()
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := svc.Templates.RenderSleepSession(ctx, w, view); err != nil {
			log.WithError(err).WithField("event", "error.SleepSession.render").Error()
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

// buildSleepSessionView locates the matching Sleep v2 - Getsummary entry by
// startdate, gathers Sleep v2 - Get segments inside the session window, and
// produces a render-ready view with prebuilt inline SVGs.
func buildSleepSessionView(
	ctx context.Context,
	repo subscription.Repo,
	accountUUID uuid.UUID,
	startdate int64,
) (SleepSessionView, error) {
	view := SleepSessionView{Startdate: startdate}

	summaryRows, err := repo.GetNotificationDataByAccountUUIDAndService(ctx, accountUUID, subscription.NotificationDataServiceSleepv2Getsummary)
	if err != nil {
		return view, fmt.Errorf("get summary rows: %w", err)
	}

	var summary withings.SleepGetsummaryEntry
	var summaryFound bool
	for _, row := range summaryRows {
		var resp withings.SleepGetsummaryResponse
		if err := json.Unmarshal(row.Data(), &resp); err != nil {
			continue
		}
		for _, entry := range resp.Body.Series {
			if int64(entry.Startdate) == startdate {
				summary = entry
				summaryFound = true
				break
			}
		}
		if summaryFound {
			break
		}
	}
	if !summaryFound {
		return view, nil
	}
	view.Found = true
	view.populateSummary(summary)

	// Pull all Sleep v2 - Get rows for the account; segments overlapping the
	// session window contribute to the charts. Withings webhooks can deliver
	// data spanning multiple sessions, so we filter inside.
	getRows, err := repo.GetNotificationDataByAccountUUIDAndService(ctx, accountUUID, subscription.NotificationDataServiceSleepv2Get)
	if err != nil {
		return view, fmt.Errorf("get sleep-get rows: %w", err)
	}

	sessionStart := int64(summary.Startdate)
	sessionEnd := int64(summary.Enddate)
	var segments []withings.SleepGetEntry
	seen := make(map[string]bool)
	for _, row := range getRows {
		var resp withings.SleepGetResponse
		if err := json.Unmarshal(row.Data(), &resp); err != nil {
			continue
		}
		for _, seg := range resp.Body.Series {
			if seg.Enddate <= sessionStart || seg.Startdate >= sessionEnd {
				continue
			}
			key := fmt.Sprintf("%d-%d-%d", seg.Startdate, seg.Enddate, seg.State)
			if seen[key] {
				continue
			}
			seen[key] = true
			segments = append(segments, seg)
		}
	}
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].Startdate < segments[j].Startdate
	})

	view.buildCharts(segments, sessionStart, sessionEnd, summary.Timezone)
	return view, nil
}

// SleepSessionView is the render-ready struct for the sleep detail template.
type SleepSessionView struct {
	Found     bool
	Startdate int64

	Date          string
	Timezone      string
	StartLocal    string
	EndLocal      string
	DurationStr   string
	SleepScore    int
	Efficiency    int    // 0–100
	Light         string // "2h 21m"
	Deep          string
	REM           string
	Awake         string
	HRMin         int
	HRAvg         int
	HRMax         int
	RRMin         int
	RRAvg         int
	RRMax         int
	BreathingDist int
	Snoring       string
	SnoringCount  int
	AHI           float64

	Hypnogram template.HTML
	HRChart   template.HTML
	RRChart   template.HTML
	HRVChart  template.HTML
	Snoring01 template.HTML
}

func (v *SleepSessionView) populateSummary(s withings.SleepGetsummaryEntry) {
	v.Date = s.Date
	v.Timezone = s.Timezone

	loc, err := time.LoadLocation(s.Timezone)
	if err != nil {
		loc = time.UTC
	}
	startT := time.Unix(int64(s.Startdate), 0).In(loc)
	endT := time.Unix(int64(s.Enddate), 0).In(loc)
	v.StartLocal = startT.Format("15:04")
	v.EndLocal = endT.Format("15:04")
	v.DurationStr = formatDuration(time.Duration(s.Data.TotalSleepTime) * time.Second)

	v.SleepScore = int(s.Data.SleepScore)
	v.Efficiency = int(s.Data.SleepEfficiency * 100)
	v.Light = formatDuration(time.Duration(s.Data.Lightsleepduration) * time.Second)
	v.Deep = formatDuration(time.Duration(s.Data.Deepsleepduration) * time.Second)
	v.REM = formatDuration(time.Duration(s.Data.Remsleepduration) * time.Second)
	v.Awake = formatDuration(time.Duration(s.Data.Wakeupduration) * time.Second)
	v.HRMin = int(s.Data.HrMin)
	v.HRAvg = int(s.Data.HrAverage)
	v.HRMax = int(s.Data.HrMax)
	v.RRMin = int(s.Data.RrMin)
	v.RRAvg = int(s.Data.RrAverage)
	v.RRMax = int(s.Data.RrMax)
	v.BreathingDist = int(s.Data.BreathingDisturbancesIntensity)
	v.Snoring = formatDuration(time.Duration(s.Data.Snoring) * time.Second)
	v.SnoringCount = s.Data.Snoringepisodecount
	v.AHI = s.Data.ApneaHypopneaIndex
}

func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "0m"
	}
	h := int(d / time.Hour)
	m := int((d % time.Hour) / time.Minute)
	if h == 0 {
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%dh %02dm", h, m)
}

// timePoint is a single (unix-second, value) sample.
type timePoint struct {
	t int64
	v float64
}

const (
	chartW         = 800
	chartH         = 120
	hypnoH         = 90
	chartMarginX   = 40
	chartMarginTop = 14
	chartMarginBot = 22
)

// buildCharts assembles the inline SVGs from the session segments.
func (v *SleepSessionView) buildCharts(segments []withings.SleepGetEntry, sessionStart, sessionEnd int64, tz string) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}

	var hr, rr, sdnn, rmssd, snoring []timePoint
	for _, seg := range segments {
		hr = append(hr, decodeSamples(seg.HR)...)
		rr = append(rr, decodeSamples(seg.RR)...)
		sdnn = append(sdnn, decodeSamples(seg.SDNN1)...)
		rmssd = append(rmssd, decodeSamples(seg.RMSSD)...)
		snoring = append(snoring, decodeSamples(seg.Snoring)...)
	}
	sortPoints(hr)
	sortPoints(rr)
	sortPoints(sdnn)
	sortPoints(rmssd)
	sortPoints(snoring)

	v.Hypnogram = template.HTML(buildHypnogramSVG(segments, sessionStart, sessionEnd, loc))
	v.HRChart = template.HTML(buildLineChartSVG("Heart rate (bpm)", hr, sessionStart, sessionEnd, loc, "#e91e63"))
	v.RRChart = template.HTML(buildLineChartSVG("Respiratory rate (rpm)", rr, sessionStart, sessionEnd, loc, "#3f51b5"))
	v.HRVChart = template.HTML(buildDualLineChartSVG("HRV (sdnn_1, rmssd)", sdnn, rmssd, sessionStart, sessionEnd, loc, "#009688", "#ff9800"))
	v.Snoring01 = template.HTML(buildSnoringSVG("Snoring (s/min)", snoring, sessionStart, sessionEnd, loc))
}

func decodeSamples(raw []byte) []timePoint {
	if len(raw) == 0 {
		return nil
	}
	var m map[string]float64
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	out := make([]timePoint, 0, len(m))
	for k, val := range m {
		ts, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			continue
		}
		out = append(out, timePoint{t: ts, v: val})
	}
	return out
}

func sortPoints(p []timePoint) {
	sort.Slice(p, func(i, j int) bool { return p[i].t < p[j].t })
}

// stateColors maps Withings sleep states to display colors.
//
//	0 = awake, 1 = light, 2 = deep, 3 = REM
var stateColors = map[int]string{
	0: "#bdbdbd",
	1: "#90caf9",
	2: "#1565c0",
	3: "#ab47bc",
}

var stateLabels = map[int]string{
	0: "Awake",
	1: "Light",
	2: "Deep",
	3: "REM",
}

func buildHypnogramSVG(segments []withings.SleepGetEntry, start, end int64, loc *time.Location) string {
	if end <= start {
		return ""
	}
	span := float64(end - start)
	innerW := float64(chartW - 2*chartMarginX)
	x := func(t int64) float64 {
		if t < start {
			t = start
		}
		if t > end {
			t = end
		}
		return float64(chartMarginX) + (float64(t-start)/span)*innerW
	}
	bandTop := chartMarginTop
	bandBot := hypnoH - chartMarginBot

	var b strings.Builder
	fmt.Fprintf(&b, `<svg viewBox="0 0 %d %d" class="sleepviz">`, chartW, hypnoH)
	fmt.Fprintf(&b, `<text x="8" y="14" font-size="11" fill="#555">Hypnogram</text>`)
	for _, seg := range segments {
		color, ok := stateColors[seg.State]
		if !ok {
			color = "#cccccc"
		}
		segStart := seg.Startdate
		if segStart < start {
			segStart = start
		}
		segEnd := seg.Enddate
		if segEnd > end {
			segEnd = end
		}
		if segEnd <= segStart {
			continue
		}
		x0 := x(segStart)
		x1 := x(segEnd)
		fmt.Fprintf(&b, `<rect x="%.1f" y="%d" width="%.1f" height="%d" fill="%s"><title>%s %s–%s</title></rect>`,
			x0, bandTop, x1-x0, bandBot-bandTop, color,
			stateLabels[seg.State],
			time.Unix(segStart, 0).In(loc).Format("15:04"),
			time.Unix(segEnd, 0).In(loc).Format("15:04"),
		)
	}
	writeTimeAxis(&b, start, end, loc, hypnoH-chartMarginBot+12)
	writeLegend(&b, chartW)
	b.WriteString(`</svg>`)
	return b.String()
}

func writeLegend(b *strings.Builder, w int) {
	x := w - 240
	for i := 0; i < 4; i++ {
		fmt.Fprintf(b, `<rect x="%d" y="2" width="10" height="10" fill="%s"/>`, x, stateColors[i])
		fmt.Fprintf(b, `<text x="%d" y="11" font-size="10" fill="#444">%s</text>`, x+13, stateLabels[i])
		x += 60
	}
}

func writeTimeAxis(b *strings.Builder, start, end int64, loc *time.Location, y int) {
	if end <= start {
		return
	}
	span := float64(end - start)
	innerW := float64(chartW - 2*chartMarginX)
	startT := time.Unix(start, 0).In(loc).Truncate(time.Hour).Add(time.Hour)
	endT := time.Unix(end, 0).In(loc)
	for t := startT; !t.After(endT); t = t.Add(time.Hour) {
		ts := t.Unix()
		if ts <= start || ts >= end {
			continue
		}
		x := float64(chartMarginX) + (float64(ts-start)/span)*innerW
		fmt.Fprintf(b, `<line x1="%.1f" y1="%d" x2="%.1f" y2="%d" stroke="#ddd"/>`, x, chartMarginTop, x, y-12)
		fmt.Fprintf(b, `<text x="%.1f" y="%d" font-size="10" text-anchor="middle" fill="#666">%s</text>`, x, y, t.Format("15:04"))
	}
}

func buildLineChartSVG(title string, points []timePoint, start, end int64, loc *time.Location, color string) string {
	return buildDualLineChartSVG(title, points, nil, start, end, loc, color, "")
}

func buildDualLineChartSVG(title string, a, b []timePoint, start, end int64, loc *time.Location, colorA, colorB string) string {
	if end <= start {
		return ""
	}
	all := make([]timePoint, 0, len(a)+len(b))
	all = append(all, a...)
	all = append(all, b...)
	if len(all) == 0 {
		return emptyChartSVG(title)
	}
	yMin, yMax := all[0].v, all[0].v
	for _, p := range all {
		if p.v < yMin {
			yMin = p.v
		}
		if p.v > yMax {
			yMax = p.v
		}
	}
	if yMax == yMin {
		yMax = yMin + 1
	}
	span := float64(end - start)
	innerW := float64(chartW - 2*chartMarginX)
	innerH := float64(chartH - chartMarginTop - chartMarginBot)
	xPos := func(t int64) float64 {
		return float64(chartMarginX) + (float64(t-start)/span)*innerW
	}
	yPos := func(v float64) float64 {
		return float64(chartMarginTop) + (1-(v-yMin)/(yMax-yMin))*innerH
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg viewBox="0 0 %d %d" class="sleepviz">`, chartW, chartH)
	fmt.Fprintf(&sb, `<text x="8" y="14" font-size="11" fill="#555">%s</text>`, template.HTMLEscapeString(title))
	fmt.Fprintf(&sb, `<text x="%d" y="%.1f" font-size="10" fill="#888" text-anchor="end">%.0f</text>`, chartMarginX-4, yPos(yMax)+3, yMax)
	fmt.Fprintf(&sb, `<text x="%d" y="%.1f" font-size="10" fill="#888" text-anchor="end">%.0f</text>`, chartMarginX-4, yPos(yMin)+3, yMin)
	writePolyline(&sb, a, xPos, yPos, colorA, start, end)
	if colorB != "" {
		writePolyline(&sb, b, xPos, yPos, colorB, start, end)
	}
	writeTimeAxis(&sb, start, end, loc, chartH-chartMarginBot+12)
	sb.WriteString(`</svg>`)
	return sb.String()
}

func writePolyline(b *strings.Builder, points []timePoint, xPos func(int64) float64, yPos func(float64) float64, color string, start, end int64) {
	if len(points) == 0 || color == "" {
		return
	}
	var sb strings.Builder
	first := true
	for _, p := range points {
		if p.t < start || p.t > end {
			continue
		}
		if !first {
			sb.WriteByte(' ')
		}
		first = false
		fmt.Fprintf(&sb, "%.1f,%.1f", xPos(p.t), yPos(p.v))
	}
	if first {
		return
	}
	fmt.Fprintf(b, `<polyline points="%s" fill="none" stroke="%s" stroke-width="1.2"/>`, sb.String(), color)
}

func buildSnoringSVG(title string, points []timePoint, start, end int64, loc *time.Location) string {
	if end <= start {
		return emptyChartSVG(title)
	}
	span := float64(end - start)
	innerW := float64(chartW - 2*chartMarginX)
	innerH := float64(chartH - chartMarginTop - chartMarginBot)
	yMax := 60.0
	for _, p := range points {
		if p.v > yMax {
			yMax = p.v
		}
	}
	xPos := func(t int64) float64 {
		return float64(chartMarginX) + (float64(t-start)/span)*innerW
	}
	yPos := func(v float64) float64 {
		return float64(chartMarginTop) + (1-v/yMax)*innerH
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<svg viewBox="0 0 %d %d" class="sleepviz">`, chartW, chartH)
	fmt.Fprintf(&b, `<text x="8" y="14" font-size="11" fill="#555">%s</text>`, template.HTMLEscapeString(title))
	for _, p := range points {
		if p.v <= 0 || p.t < start || p.t > end {
			continue
		}
		x := xPos(p.t)
		yTop := yPos(p.v)
		yBot := yPos(0)
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="1.5" height="%.1f" fill="#5d4037"/>`, x-0.75, yTop, yBot-yTop)
	}
	writeTimeAxis(&b, start, end, loc, chartH-chartMarginBot+12)
	b.WriteString(`</svg>`)
	return b.String()
}

func emptyChartSVG(title string) string {
	return fmt.Sprintf(`<svg viewBox="0 0 %d %d" class="sleepviz"><text x="8" y="14" font-size="11" fill="#555">%s</text><text x="%d" y="%d" font-size="10" fill="#aaa" text-anchor="middle">no samples in window</text></svg>`,
		chartW, chartH, template.HTMLEscapeString(title), chartW/2, chartH/2)
}
