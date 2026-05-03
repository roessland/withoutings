package port_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/integrationtest"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// One coherent night fixture: Getsummary entry references the same window that
// the Get segments cover. We construct minimal payloads inline so the assertions
// below check the parser and chart wiring, not Withings JSON nuances.
const sleepSummaryFixture = `{
  "status": 0,
  "body": {
    "series": [{
      "id": 1,
      "timezone": "UTC",
      "model": 32,
      "model_id": 63,
      "startdate": 1700000000,
      "enddate": 1700025200,
      "date": "2023-11-14",
      "data": {
        "total_sleep_time": 25200,
        "total_timeinbed": 25200,
        "lightsleepduration": 14400,
        "deepsleepduration": 7200,
        "remsleepduration": 3600,
        "sleep_efficiency": 0.92,
        "sleep_score": 87,
        "wakeupduration": 600,
        "wakeupcount": 1,
        "hr_average": 58,
        "hr_min": 50,
        "hr_max": 72,
        "rr_average": 16,
        "rr_min": 12,
        "rr_max": 22,
        "snoring": 120,
        "snoringepisodecount": 2,
        "apnea_hypopnea_index": 1.5,
        "breathing_disturbances_intensity": 3
      }
    }]
  }
}`

const sleepGetFixture = `{
  "status": 0,
  "body": {
    "series": [
      {"startdate": 1700000000, "enddate": 1700003600, "state": 1, "model": "Aura", "model_id": 63,
       "hr": {"1700000000": 60, "1700001800": 58},
       "rr": {"1700000000": 17, "1700001800": 16},
       "snoring": {"1700000000": 0, "1700001800": 30},
       "sdnn_1": {"1700000000": 40, "1700001800": 42},
       "rmssd": {"1700000000": 35, "1700001800": 36}},
      {"startdate": 1700003600, "enddate": 1700014400, "state": 2, "model": "Aura", "model_id": 63,
       "hr": {"1700003600": 52, "1700010000": 50},
       "rr": {"1700003600": 14, "1700010000": 13},
       "snoring": {"1700003600": 0},
       "sdnn_1": {"1700003600": 55},
       "rmssd": {"1700003600": 50}},
      {"startdate": 1700014400, "enddate": 1700020000, "state": 3, "model": "Aura", "model_id": 63,
       "hr": {"1700014400": 65, "1700017000": 70},
       "rr": {"1700014400": 18, "1700017000": 19},
       "snoring": {"1700014400": 0},
       "sdnn_1": {"1700014400": 48},
       "rmssd": {"1700014400": 55}},
      {"startdate": 1700020000, "enddate": 1700025200, "state": 0, "model": "Aura", "model_id": 63,
       "hr": {"1700020000": 70},
       "rr": {"1700020000": 20},
       "snoring": {"1700020000": 0},
       "sdnn_1": {"1700020000": 30},
       "rmssd": {"1700020000": 25}}
    ]
  }
}`

func TestSleepSessionPage(t *testing.T) {
	it := integrationtest.WithFreshDatabase(t)

	var acc *account.Account

	doSessionReq := func(startdate int64) (*httptest.ResponseRecorder, string) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sleepsessions/%d", startdate), nil)
		req = req.WithContext(account.AddToContext(it.Ctx, acc))
		return it.DoRequest(req)
	}

	beforeEach := func(t *testing.T) {
		it.ResetMocks(t)
		acc = it.MakeNewAccount(t)
	}

	insertSession := func(t *testing.T, ctx context.Context, summaryJSON, getJSON string) uuid.UUID {
		t.Helper()
		notif := subscription.MustNewNotification(subscription.NewNotificationParams{
			NotificationUUID:    uuid.New(),
			AccountUUID:         acc.UUID(),
			ReceivedAt:          time.Now(),
			Params:              "userid=1&startdate=1700000000&enddate=1700025200&appli=44",
			DataStatus:          subscription.NotificationDataStatusFetched,
			FetchedAt:           ptrTime(time.Now()),
			RawNotificationUUID: uuid.New(),
			Source:              "test",
		})
		require.NoError(t, it.App.SubscriptionRepo.CreateNotification(ctx, notif))
		require.NoError(t, it.App.SubscriptionRepo.StoreNotificationData(ctx, subscription.MustNewNotificationData(subscription.NewNotificationDataParams{
			NotificationDataUUID: uuid.New(),
			NotificationUUID:     notif.UUID(),
			AccountUUID:          acc.UUID(),
			Service:              subscription.NotificationDataServiceSleepv2Getsummary,
			Data:                 []byte(summaryJSON),
			FetchedAt:            time.Now(),
		})))
		require.NoError(t, it.App.SubscriptionRepo.StoreNotificationData(ctx, subscription.MustNewNotificationData(subscription.NewNotificationDataParams{
			NotificationDataUUID: uuid.New(),
			NotificationUUID:     notif.UUID(),
			AccountUUID:          acc.UUID(),
			Service:              subscription.NotificationDataServiceSleepv2Get,
			Data:                 []byte(getJSON),
			FetchedAt:            time.Now(),
		})))
		return notif.UUID()
	}

	t.Run("renders summary stats and inline SVG charts", func(t *testing.T) {
		beforeEach(t)
		insertSession(t, it.Ctx, sleepSummaryFixture, sleepGetFixture)

		resp, body := doSessionReq(1700000000)
		assert.Equal(t, 200, resp.Code)

		// Summary card values
		assert.Contains(t, body, "Sleep session — 2023-11-14")
		assert.Contains(t, body, ">87<")    // sleep score
		assert.Contains(t, body, "7h 00m")  // total sleep
		assert.Contains(t, body, "92%")     // efficiency
		assert.Contains(t, body, ">58 bpm") // hr average
		assert.Contains(t, body, ">16 rpm") // rr average
		assert.Contains(t, body, "1.5")     // AHI

		// Inline charts
		assert.Contains(t, body, "Hypnogram")
		assert.Contains(t, body, "Heart rate (bpm)")
		assert.Contains(t, body, "Respiratory rate (rpm)")
		assert.Contains(t, body, "HRV (sdnn_1, rmssd)")
		// Hypnogram must include the same color the SVG renders for each
		// state present in the fixture. Sourcing from the production map
		// instead of duplicating hex codes keeps tests coupled to one map.
		for _, state := range []int{1, 2, 3} {
			assert.Contains(t, body, port.SleepStateInfoByState[state].Color, "missing color for state %d", state)
		}
		// Charts must include polylines (multi-sample series), not just titles.
		assert.Contains(t, body, "<polyline ")
		// y-axis bounds for HR fixture are 50..72 (from summary), but the
		// chart bounds come from the Sleep v2 - Get samples (50..60), so
		// assert the rendered bounds are present.
		assert.Contains(t, body, "stroke=\"#e91e63\"") // HR series color
	})

	t.Run("missing session shows friendly empty state", func(t *testing.T) {
		beforeEach(t)
		// Insert a session, then ask for a different startdate.
		insertSession(t, it.Ctx, sleepSummaryFixture, sleepGetFixture)

		resp, body := doSessionReq(1700099999)
		assert.Equal(t, 200, resp.Code)
		assert.Contains(t, body, "No stored data for this session")
		// Empty state must not render any chart chrome — it's how we verify
		// the handler returned `Found: false` rather than rendering with
		// missing-data noise.
		assert.NotContains(t, body, "Hypnogram")
		assert.NotContains(t, body, "<polyline")
		assert.NotContains(t, body, "<svg")
	})

	t.Run("unauthenticated request gets 401", func(t *testing.T) {
		beforeEach(t)
		// Request without an account in context — exercises the guard at the
		// top of the handler, which previously had no test.
		req := httptest.NewRequest(http.MethodGet, "/sleepsessions/1700000000", nil)
		req = req.WithContext(it.Ctx)
		resp, body := it.DoRequest(req)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, body, "must be logged in")
	})

	t.Run("non-UTC fixture renders with local hour-tick labels", func(t *testing.T) {
		// 2026-03-29 spans the spring-forward transition in Europe/Oslo
		// (02:00 -> 03:00 CET -> CEST). The hour-tick loop must skip the
		// jumped-over hour, not stack labels. We assert at minimum that the
		// rendered axis carries one local 02:00 label or one local 03:00
		// label, never both — that's the regression signal.
		beforeEach(t)

		// 2026-03-29 00:00 Europe/Oslo = 1774738800 (CET, +01:00).
		// 8h sleep ending 2026-03-29 08:00 local = 1774764000.
		const summary = `{"status":0,"body":{"series":[{"id":1,"timezone":"Europe/Oslo","model":32,"model_id":63,"startdate":1774738800,"enddate":1774764000,"date":"2026-03-29","data":{"total_sleep_time":25200,"sleep_score":80,"sleep_efficiency":0.9,"hr_average":58,"rr_average":16}}]}}`
		const get = `{"status":0,"body":{"series":[{"startdate":1774738800,"enddate":1774764000,"state":1,"hr":{"1774738800":58,"1774762800":60},"rr":{"1774738800":16,"1774762800":17}}]}}`

		insertSession(t, it.Ctx, summary, get)

		resp, body := doSessionReq(1774738800)
		assert.Equal(t, 200, resp.Code)
		assert.Contains(t, body, "Sleep session — 2026-03-29")
		assert.Contains(t, body, "Europe/Oslo")
		// Spring-forward skips 02:00 in Europe/Oslo; the spring-skipped label
		// must not appear at all. The 03:00 label must be present (once).
		assert.Equal(t, 0, countOccurrences(body, ">02:00<"), "02:00 must not be rendered during spring-forward")
		assert.GreaterOrEqual(t, countOccurrences(body, ">03:00<"), 1, "03:00 must be rendered as a real wall-clock hour")
		// Wall-clock 04:00 sanity check — appears in the morning hours of the night.
		assert.GreaterOrEqual(t, countOccurrences(body, ">04:00<"), 1)
	})

	t.Run("rejects non-numeric startdate", func(t *testing.T) {
		beforeEach(t)
		// gorilla/mux returns 404 for path that doesn't match the regex.
		req := httptest.NewRequest(http.MethodGet, "/sleepsessions/notanumber", nil)
		req = req.WithContext(account.AddToContext(it.Ctx, acc))
		resp, _ := it.DoRequest(req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}

func ptrTime(t time.Time) *time.Time { return &t }

func countOccurrences(haystack, needle string) int { return strings.Count(haystack, needle) }
