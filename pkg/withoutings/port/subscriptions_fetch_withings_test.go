package port_test

import (
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/integrationtest"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSubscriptionsWithingsPage(t *testing.T) {
	it := integrationtest.WithFreshDatabase(t)

	var accountUUID uuid.UUID
	var withingsUserID string

	notificationCategories, err := it.App.SubscriptionRepo.AllNotificationCategories(it.Ctx)
	require.NoError(t, err)

	mockBeingSubscribedToEverything := func() {
		// Mock a response for each notification category.
		for _, cat := range notificationCategories {
			it.Mocks.MockWithingsSvc.EXPECT().
				NotifyList(mock.Anything, mock.Anything, withings.NewNotifyListParams(cat.Appli)).
				Once().Return(&withings.NotifyListResponse{
				Status: 0,
				Body: withings.NotifyListBody{
					Profiles: []withings.NotifyListProfile{
						{
							Appli:       cat.Appli,
							CallbackURL: "https://abc/123",
							Comment:     "abc",
						},
					},
				},
				Raw: []byte("whatever"),
			}, nil)
		}
	}

	mockBeingSubscribedToNothing := func() {
		// Mock a response for each notification category.
		for _, cat := range notificationCategories {
			it.Mocks.MockWithingsSvc.EXPECT().
				NotifyList(mock.Anything, mock.Anything, withings.NewNotifyListParams(cat.Appli)).
				Once().Return(&withings.NotifyListResponse{
				Status: 0,
				Body: withings.NotifyListBody{
					Profiles: []withings.NotifyListProfile{},
				},
				Raw: []byte("whatever"),
			}, nil)
		}
	}

	loggedInRequest := func() *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/withings", nil)

		// TODO deduplicate this
		req = req.WithContext(account.AddToContext(it.Ctx,
			account.NewAccountOrPanic(
				accountUUID,
				withingsUserID,
				"bob",
				"kåre",
				time.Now().Add(-time.Hour),
				"whatever",
			),
		))
		return req
	}

	doRequest := func(req *http.Request) (*httptest.ResponseRecorder, string) {
		// Should be success
		resp := httptest.NewRecorder()
		it.Router.ServeHTTP(resp, req)
		respBody, _ := io.ReadAll(resp.Body)
		return resp, string(respBody)
	}

	beforeEach := func(t *testing.T) {
		it.ResetMocks(t)

		// Insert a user with an expired access token.
		accountUUID = uuid.New()
		withingsUserID = uuid.NewString()
		acc, err := account.NewAccount(
			accountUUID,
			withingsUserID,
			"bob",
			"kåre",
			time.Now().Add(-time.Hour),
			"whatever",
		)
		require.NoError(t, err)
		require.NoError(t, it.App.AccountRepo.CreateAccount(it.Ctx, acc))

	}

	t.Run("shows one row per active subscription", func(t *testing.T) {
		beforeEach(t)
		mockBeingSubscribedToEverything()
		req := loggedInRequest()

		resp, body := doRequest(req)
		assert.Equal(t, 200, resp.Code, body)

		require.Equal(t, strings.Count(body, "✅"), len(notificationCategories))
	})

	t.Run("shows one row per inactive subscription", func(t *testing.T) {
		beforeEach(t)
		mockBeingSubscribedToNothing()
		req := loggedInRequest()

		resp, body := doRequest(req)
		assert.Equal(t, 200, resp.Code, body)

		require.Equal(t, strings.Count(body, "❌"), len(notificationCategories))
	})
}
