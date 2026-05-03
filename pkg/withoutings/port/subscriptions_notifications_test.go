package port_test

import (
	"context"
	"fmt"
	"github.com/roessland/withoutings/pkg/integrationtest"
	"github.com/roessland/withoutings/pkg/withoutings/adapter/withings/withingstestdata"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNotificationsPage(t *testing.T) {
	it := integrationtest.WithFreshDatabase(t)

	var acc *account.Account

	newListNotificationsReq := func() *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
		req = req.WithContext(account.AddToContext(it.Ctx, acc))
		return req
	}

	simulateIncomingWebhook := func(payload string) {
		req := httptest.NewRequest(http.MethodPost,
			"/withings/webhooks/qwerty1234",
			strings.NewReader(payload))
		resp, _ := it.DoRequest(req)
		assert.Equal(t, 200, resp.Code)
	}

	beforeEach := func(t *testing.T) {
		it.ResetMocks(t)
		acc = it.MakeNewAccount(t)
	}

	t.Run("should respond with 200 OK listing notifications", func(t *testing.T) {
		beforeEach(t)

		workerCtx, cancelWorker := context.WithCancel(it.Ctx)
		defer cancelWorker()

		it.Mocks.MockWithingsSvc.EXPECT().
			MeasureGetmeas(mock.Anything, mock.Anything, mock.Anything).
			Return(withings.MustNewMeasureGetmeasResponse(withingstestdata.MeasureGetmeasSuccess), nil)

		defer it.Mocks.MockWithingsSvc.AssertExpectations(t)

		go it.Worker.Work(workerCtx)

		weighInNotificationParams := fmt.Sprintf(`userid=%s&startdate=1684738999&enddate=1684739000&appli=1`, acc.WithingsUserID())
		simulateIncomingWebhook("") // Initial heartbeat
		simulateIncomingWebhook(weighInNotificationParams)

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			resp, body := it.DoRequest(newListNotificationsReq())
			assert.Equal(c, 200, resp.Code)
			assert.Contains(c, body, "Measure - Getmeas")
		}, 10*time.Second, 300*time.Millisecond, "should show received notifications with links to fetched payloads")
	})

	t.Run("should fetch new activity data when receiving appli=16 webhook", func(t *testing.T) {
		// appli=16 (New activity-related data) has three services to call:
		// - Measure v2 - Getactivity
		// - Measure v2 - Getintradayactivity
		// - Measure v2 - Getworkouts
		// The webhook payload uses date=YYYY-MM-DD (not startdate/enddate).

		beforeEach(t)

		workerCtx, cancelWorker := context.WithCancel(it.Ctx)
		defer cancelWorker()

		// Pin params so a regression in the YMD-vs-Unix split is caught.
		// Daily endpoints use the date directly. Intraday uses a precise
		// local-day window derived from the timezone field on the first
		// activity row in the Getactivity response (Europe/Oslo, +01:00 in
		// January) — see the success fixture.
		getactivityParams := withings.NewMeasureGetactivityParams()
		getactivityParams.Startdateymd = "2022-01-26"
		getactivityParams.Enddateymd = "2022-01-26"

		intradayParams := withings.NewMeasureGetintradayactivityParams()
		intradayParams.Startdate = 1643151600 // 2022-01-26 00:00 Europe/Oslo
		intradayParams.Enddate = 1643238000   // 2022-01-27 00:00 Europe/Oslo

		getworkoutsParams := withings.NewMeasureGetworkoutsParams()
		getworkoutsParams.Startdateymd = "2022-01-26"
		getworkoutsParams.Enddateymd = "2022-01-26"

		it.Mocks.MockWithingsSvc.EXPECT().
			MeasureGetactivity(mock.Anything, mock.Anything, getactivityParams).
			Return(withings.MustNewMeasureGetactivityResponse(withingstestdata.MeasureGetactivitySuccess), nil)
		it.Mocks.MockWithingsSvc.EXPECT().
			MeasureGetintradayactivity(mock.Anything, mock.Anything, intradayParams).
			Return(withings.MustNewMeasureGetintradayactivityResponse(withingstestdata.MeasureGetintradayactivitySuccess), nil)
		it.Mocks.MockWithingsSvc.EXPECT().
			MeasureGetworkouts(mock.Anything, mock.Anything, getworkoutsParams).
			Return(withings.MustNewMeasureGetworkoutsResponse(withingstestdata.MeasureGetworkoutsSuccess), nil)

		defer it.Mocks.MockWithingsSvc.AssertExpectations(t)

		go it.Worker.Work(workerCtx)

		activity16WebhookParams := fmt.Sprintf(`userid=%s&appli=16&date=2022-01-26`, acc.WithingsUserID())
		simulateIncomingWebhook(activity16WebhookParams)

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			resp, body := it.DoRequest(newListNotificationsReq())
			assert.Equal(c, 200, resp.Code)

			assert.Contains(c, body, "Measure v2 - Getactivity<")
			assert.Contains(c, body, "Measure v2 - Getintradayactivity<")
			assert.Contains(c, body, "Measure v2 - Getworkouts<")
		}, 10*time.Second, 300*time.Millisecond, "should show received activity notification with three fetched payloads")
	})

	t.Run("should still record appli=16 notification when Withings returns no data", func(t *testing.T) {
		beforeEach(t)

		workerCtx, cancelWorker := context.WithCancel(it.Ctx)
		defer cancelWorker()

		// With no activity rows we can't read a timezone, so the intraday
		// window falls back to UTC midnight-to-midnight.
		intradayParams := withings.NewMeasureGetintradayactivityParams()
		intradayParams.Startdate = 1643155200 // 2022-01-26 00:00 UTC
		intradayParams.Enddate = 1643241600   // 2022-01-27 00:00 UTC

		it.Mocks.MockWithingsSvc.EXPECT().
			MeasureGetactivity(mock.Anything, mock.Anything, mock.Anything).
			Return(withings.MustNewMeasureGetactivityResponse(withingstestdata.MeasureGetactivityNoData), nil)
		it.Mocks.MockWithingsSvc.EXPECT().
			MeasureGetintradayactivity(mock.Anything, mock.Anything, intradayParams).
			Return(withings.MustNewMeasureGetintradayactivityResponse(withingstestdata.MeasureGetintradayactivityNoData), nil)
		it.Mocks.MockWithingsSvc.EXPECT().
			MeasureGetworkouts(mock.Anything, mock.Anything, mock.Anything).
			Return(withings.MustNewMeasureGetworkoutsResponse(withingstestdata.MeasureGetworkoutsNoData), nil)

		defer it.Mocks.MockWithingsSvc.AssertExpectations(t)

		go it.Worker.Work(workerCtx)

		simulateIncomingWebhook(fmt.Sprintf(`userid=%s&appli=16&date=2022-01-26`, acc.WithingsUserID()))

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			resp, body := it.DoRequest(newListNotificationsReq())
			assert.Equal(c, 200, resp.Code)

			assert.Contains(c, body, "Measure v2 - Getactivity<")
			assert.Contains(c, body, "Measure v2 - Getintradayactivity<")
			assert.Contains(c, body, "Measure v2 - Getworkouts<")
		}, 10*time.Second, 300*time.Millisecond, "empty payloads should still produce three NotificationData rows")
	})

	t.Run("should not loop forever when appli=16 webhook has malformed date", func(t *testing.T) {
		// A malformed date= must not propagate as an error to watermill,
		// or the message will redeliver indefinitely.
		beforeEach(t)

		workerCtx, cancelWorker := context.WithCancel(it.Ctx)
		defer cancelWorker()

		// No EXPECT(): any Withings call would be a regression. AssertExpectations
		// also fails the test if any unexpected call is made.
		defer it.Mocks.MockWithingsSvc.AssertExpectations(t)

		go it.Worker.Work(workerCtx)

		simulateIncomingWebhook(fmt.Sprintf(`userid=%s&appli=16&date=not-a-date`, acc.WithingsUserID()))

		// The notification should be processed (handler returns nil) so the
		// page renders, and no Withings service rows show up for it.
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			resp, body := it.DoRequest(newListNotificationsReq())
			assert.Equal(c, 200, resp.Code)
			assert.NotContains(c, body, "Measure v2 - ")
		}, 10*time.Second, 300*time.Millisecond, "malformed date should be handled without retries or service rows")
	})

	t.Run("should fetch new available data from multiple services", func(t *testing.T) {
		// Some notification categories have multiple services to call,
		// in order to retrieve all the relevant data.
		// For example appli=44 (New sleep-related data) has these services
		// to call:
		// - Sleep v2 - Get
		// - Sleep v2 - Getsummary

		beforeEach(t)

		workerCtx, cancelWorker := context.WithCancel(it.Ctx)
		defer cancelWorker()

		it.Mocks.MockWithingsSvc.EXPECT().
			SleepGetsummary(mock.Anything, mock.Anything, mock.Anything).
			Return(withings.MustNewSleepGetsummaryResponse(withingstestdata.Sleepv2GetsummarySuccess), nil)
		sleepGetParams := withings.SleepGetParams{
			Action:     "get",
			Startdate:  1709336580,
			Enddate:    1709368500,
			DataFields: withings.SleepGetAllDataFields,
		}
		it.Mocks.MockWithingsSvc.EXPECT().
			SleepGet(mock.Anything, mock.Anything, sleepGetParams).
			Return(withings.MustNewSleepGetResponse(withingstestdata.Sleepv2GetSuccess), nil)

		defer it.Mocks.MockWithingsSvc.AssertExpectations(t)

		go it.Worker.Work(workerCtx)

		sleep44WebhookParams := fmt.Sprintf(`userid=%s&startdate=1709336580&enddate=1709368500&appli=44`, acc.WithingsUserID())
		simulateIncomingWebhook(sleep44WebhookParams)

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			resp, body := it.DoRequest(newListNotificationsReq())
			assert.Equal(c, 200, resp.Code)

			assert.Contains(c, body, "Sleep v2 - Getsummary<")
			assert.Contains(c, body, "Sleep v2 - Get<")
		}, 10*time.Second, 300*time.Millisecond, "should show received notifications with multiple fetched payloads")
	})
}
