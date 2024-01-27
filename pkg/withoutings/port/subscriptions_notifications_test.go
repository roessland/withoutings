package port_test

import (
	"context"
	"fmt"
	"github.com/roessland/withoutings/pkg/integrationtest"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"html"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNotificationsPage(t *testing.T) {
	// TODO this is not really an integration test, since it doesn't use the database.
	// It should be converted to a unit test.
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
		go it.Worker.Work(workerCtx)

		// TODO: real response
		it.Mocks.MockWithingsSvc.EXPECT().MeasureGetmeas(mock.Anything, mock.Anything, mock.Anything).Return(withings.MustNewMeasureGetmeasResponse(
			[]byte(`
				{
					"status": 0,
					"body": {
						"updatetime": "string",
						"timezone": "string",
						"measuregrps": [
							{
								"grpid": 12,
								"attrib": 1,
								"date": 1594245600,
								"created": 1594246600,
								"modified": 1594257200,
								"category": 1594257200,
								"deviceid": "892359876fd8805ac45bab078c4828692f0276b1",
								"measures": [
									{
										"value": 65750,
										"type": 1,
										"unit": -3,
										"algo": 3425,
										"fm": 0,
										"position": 1
									}
								],
								"comment": "A measurement comment",
								"timezone": "Europe/Paris"
							}
						],
						"more": 0,
						"offset": 0
					}
				}
			`),
		), nil)

		// give worker time to register subscriptions, since gochannel pubsub is used for now.
		// TODO: replace with SQL-based pubsub
		time.Sleep(100 * time.Millisecond)
		weighInNotificationParams := fmt.Sprintf(`userid=%s&startdate=1684738999&enddate=1684739000&appli=1`, acc.WithingsUserID())
		simulateIncomingWebhook(weighInNotificationParams)
		escapedParams := html.EscapeString(weighInNotificationParams)

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			resp, body := it.DoRequest(newListNotificationsReq())
			assert.Equal(c, 200, resp.Code)
			assert.Contains(c, body, escapedParams)
		}, 1*time.Second, 100*time.Millisecond, "should show received notifications")
	})
}
