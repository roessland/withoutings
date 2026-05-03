package port_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/integrationtest"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
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
		// hypnogram should reference each state color used in the fixture
		assert.Contains(t, body, stateColorAssert(1)) // light
		assert.Contains(t, body, stateColorAssert(2)) // deep
		assert.Contains(t, body, stateColorAssert(3)) // rem
	})

	t.Run("missing session shows friendly empty state", func(t *testing.T) {
		beforeEach(t)
		// Insert a session, then ask for a different startdate.
		insertSession(t, it.Ctx, sleepSummaryFixture, sleepGetFixture)

		resp, body := doSessionReq(1700099999)
		assert.Equal(t, 200, resp.Code)
		assert.Contains(t, body, "No stored data for this session")
		assert.NotContains(t, body, "Hypnogram")
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

// stateColorAssert returns the hex color used for a Withings sleep state in the
// rendered SVG, mirroring the table in sleep_session.go (kept private). If the
// production map changes, this helper must move with it.
func stateColorAssert(state int) string {
	switch state {
	case 0:
		return "#bdbdbd"
	case 1:
		return "#90caf9"
	case 2:
		return "#1565c0"
	case 3:
		return "#ab47bc"
	}
	return ""
}
