package port

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/web/templates"
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
) (templates.SleepSessionView, error) {
	view := templates.SleepSessionView{Startdate: startdate}
	log := logging.MustGetLoggerFromContext(ctx)

	// Single indexed lookup against the JSONB GIN index — pulls the one row
	// (if any) whose body.series array contains an entry with this startdate.
	// Avoids loading every Sleep v2 - Getsummary blob for the user.
	summaryRow, err := repo.GetNotificationDataByAccountAndServiceAndSeriesStartdate(
		ctx, accountUUID, subscription.NotificationDataServiceSleepv2Getsummary, startdate,
	)
	if err != nil {
		return view, fmt.Errorf("get summary row: %w", err)
	}
	if summaryRow == nil {
		return view, nil
	}

	var summary withings.SleepGetsummaryEntry
	var summaryFound bool
	var resp withings.SleepGetsummaryResponse
	if err := json.Unmarshal(summaryRow.Data(), &resp); err != nil {
		log.WithError(err).
			WithField("event", "warn.SleepSession.summary_unmarshal_failed").
			WithField("notification_data_uuid", summaryRow.UUID()).
			Warn()
		return view, nil
	}
	for _, entry := range resp.Body.Series {
		if int64(entry.Startdate) == startdate {
			summary = entry
			summaryFound = true
			break
		}
	}
	if !summaryFound {
		// JSONB containment matched but the Go-side scan didn't — would only
		// happen if the index says yes but the row's content disagrees.
		return view, nil
	}
	view.Found = true
	populateSummary(ctx, &view, summary)

	sessionStart := int64(summary.Startdate)
	sessionEnd := int64(summary.Enddate)

	// Pull only the Sleep v2 - Get rows whose body.series contains at least
	// one segment overlapping the session window. Filtering in SQL avoids
	// transferring (and Go-unmarshalling) every stored Get blob — each row is
	// hundreds of KB once a user has months of webhook history.
	getRows, err := repo.GetNotificationDataByAccountAndServiceAndOverlappingWindow(
		ctx, accountUUID, subscription.NotificationDataServiceSleepv2Get, sessionStart, sessionEnd,
	)
	if err != nil {
		return view, fmt.Errorf("get sleep-get rows: %w", err)
	}

	var segments []withings.SleepGetEntry
	seen := make(map[string]bool)
	for _, row := range getRows {
		var resp withings.SleepGetResponse
		if err := json.Unmarshal(row.Data(), &resp); err != nil {
			log.WithError(err).
				WithField("event", "warn.SleepSession.get_unmarshal_failed").
				WithField("notification_data_uuid", row.UUID()).
				Warn()
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

	buildCharts(ctx, &view, segments, sessionStart, sessionEnd, summary.Timezone)
	return view, nil
}

func populateSummary(ctx context.Context, v *templates.SleepSessionView, s withings.SleepGetsummaryEntry) {
	v.Date = s.Date
	v.Timezone = s.Timezone

	loc := loadLocationOrUTC(ctx, s.Timezone)
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
func buildCharts(ctx context.Context, v *templates.SleepSessionView, segments []withings.SleepGetEntry, sessionStart, sessionEnd int64, tz string) {
	loc := loadLocationOrUTC(ctx, tz)

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
	v.SnoringChart = template.HTML(buildSnoringSVG("Snoring (s/min)", snoring, sessionStart, sessionEnd, loc))
}

// loadLocationOrUTC resolves a Withings timezone name with a UTC fallback,
// emitting one warn-level log per failure so distroless / missing-tzdata
// regressions and bad summary payloads aren't silently invisible.
func loadLocationOrUTC(ctx context.Context, tz string) *time.Location {
	loc, err := time.LoadLocation(tz)
	if err == nil {
		return loc
	}
	logging.MustGetLoggerFromContext(ctx).
		WithError(err).
		WithField("event", "warn.SleepSession.load_location_failed").
		WithField("tz", tz).
		Warn()
	return time.UTC
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
		// Drop NaN/Inf early. Otherwise yMin/yMax seeding in the line-chart
		// builder seeds from a non-finite value (NaN comparisons are always
		// false), every projected coordinate is NaN, and the polyline renders
		// as `NaN,NaN ...` — invisible and indistinguishable from no data.
		if math.IsNaN(val) || math.IsInf(val, 0) {
			continue
		}
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

// SleepStateInfo is the display metadata for a Withings sleep state.
// Exposed so tests can assert against the same source of truth the SVG uses.
type SleepStateInfo struct {
	Color string
	Label string
}

// SleepStateInfoByState maps Withings sleep states (0=awake, 1=light, 2=deep,
// 3=REM) to display metadata. States 4-6 (deprecated/manual per Withings docs)
// fall through to a neutral fallback in buildHypnogramSVG.
var SleepStateInfoByState = map[int]SleepStateInfo{
	0: {Color: "#bdbdbd", Label: "Awake"},
	1: {Color: "#90caf9", Label: "Light"},
	2: {Color: "#1565c0", Label: "Deep"},
	3: {Color: "#ab47bc", Label: "REM"},
}

const unknownStateColor = "#cccccc"

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
		info, known := SleepStateInfoByState[seg.State]
		color := info.Color
		label := info.Label
		if !known {
			color = unknownStateColor
			label = fmt.Sprintf("State %d", seg.State)
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
			label,
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
		info := SleepStateInfoByState[i]
		fmt.Fprintf(b, `<rect x="%d" y="2" width="10" height="10" fill="%s"/>`, x, info.Color)
		fmt.Fprintf(b, `<text x="%d" y="11" font-size="10" fill="#444">%s</text>`, x+13, info.Label)
		x += 60
	}
}

func writeTimeAxis(b *strings.Builder, start, end int64, loc *time.Location, y int) {
	if end <= start {
		return
	}
	span := float64(end - start)
	innerW := float64(chartW - 2*chartMarginX)
	// Compute the next local-hour boundary explicitly via time.Date so the
	// axis lands on wall-clock hours in fractional-offset zones (e.g.
	// Asia/Kolkata +05:30) and during DST transitions.
	startLocal := time.Unix(start, 0).In(loc)
	t := time.Date(startLocal.Year(), startLocal.Month(), startLocal.Day(), startLocal.Hour()+1, 0, 0, 0, loc)
	endT := time.Unix(end, 0).In(loc)
	seenLabel := make(map[string]bool)
	for ; !t.After(endT); t = t.Add(time.Hour) {
		ts := t.Unix()
		if ts <= start || ts >= end {
			continue
		}
		// On DST fall-back the loop walks two unix-hours that share a wall
		// clock label (e.g. 02:00 twice). Render the first only; the second
		// would stack a duplicate tick on top of the first.
		label := t.Format("15:04")
		if seenLabel[label] {
			continue
		}
		seenLabel[label] = true
		x := float64(chartMarginX) + (float64(ts-start)/span)*innerW
		fmt.Fprintf(b, `<line x1="%.1f" y1="%d" x2="%.1f" y2="%d" stroke="#ddd"/>`, x, chartMarginTop, x, y-12)
		fmt.Fprintf(b, `<text x="%.1f" y="%d" font-size="10" text-anchor="middle" fill="#666">%s</text>`, x, y, label)
	}
}

type lineSeries struct {
	points []timePoint
	color  string
}

func buildLineChartSVG(title string, points []timePoint, start, end int64, loc *time.Location, color string) string {
	return buildLineChartsSVG(title, start, end, loc, lineSeries{points: points, color: color})
}

func buildDualLineChartSVG(title string, a, b []timePoint, start, end int64, loc *time.Location, colorA, colorB string) string {
	return buildLineChartsSVG(title, start, end, loc,
		lineSeries{points: a, color: colorA},
		lineSeries{points: b, color: colorB},
	)
}

func buildLineChartsSVG(title string, start, end int64, loc *time.Location, series ...lineSeries) string {
	if end <= start {
		return ""
	}
	totalPoints := 0
	for _, s := range series {
		totalPoints += len(s.points)
	}
	if totalPoints == 0 {
		return emptyChartSVG(title)
	}
	var yMin, yMax float64
	seeded := false
	for _, s := range series {
		for _, p := range s.points {
			if !seeded {
				yMin, yMax, seeded = p.v, p.v, true
				continue
			}
			if p.v < yMin {
				yMin = p.v
			}
			if p.v > yMax {
				yMax = p.v
			}
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
	for _, s := range series {
		writePolyline(&sb, s.points, xPos, yPos, s.color, start, end)
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
	type xy struct{ x, y float64 }
	var inWindow []xy
	for _, p := range points {
		if p.t < start || p.t > end {
			continue
		}
		inWindow = append(inWindow, xy{x: xPos(p.t), y: yPos(p.v)})
	}
	if len(inWindow) == 0 {
		return
	}
	if len(inWindow) == 1 {
		// A single sample renders a polyline with one vertex, which draws
		// nothing. Emit a small circle so the data point is at least visible.
		fmt.Fprintf(b, `<circle cx="%.1f" cy="%.1f" r="1.5" fill="%s"/>`, inWindow[0].x, inWindow[0].y, color)
		return
	}
	for i, p := range inWindow {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%.1f,%.1f", p.x, p.y)
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
