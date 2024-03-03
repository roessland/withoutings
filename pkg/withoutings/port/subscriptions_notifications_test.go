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
