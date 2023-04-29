package handlers_test

import (
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/integrationtest"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/roessland/withoutings/web/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TODO simplify handler tests. extract shared code.
func TestRefreshWithingsAccessToken(t *testing.T) {
	it := integrationtest.WithFreshDatabase(t)

	var accountUUID uuid.UUID
	var withingsUserID string

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

	t.Run("with expired token refreshes token", func(t *testing.T) {
		beforeEach(t)

		it.Mocks.MockWithingsRepo.EXPECT().
			RefreshAccessToken(mock.Anything, mock.Anything).
			Once().Return(&withings.Token{
			UserID:       withingsUserID,
			AccessToken:  "a075f8c14fb8df40b08ebc8508533dc332a6910a",
			RefreshToken: "f631236f02b991810feb774765b6ae8e6c6839ca",
			ExpiresIn:    10800,
			Scope:        "user.info,user.metrics",
			CSRFToken:    "PACnnxwHTaBQOzF7bQqwFUUotIuvtzSM",
			TokenType:    "Bearer",
			Expiry:       time.Now().Add(10800 * time.Second),
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
		req = req.WithContext(middleware.AddAccountToContext(it.Ctx,
			account.NewAccountOrPanic(
				accountUUID,
				withingsUserID,
				"bob",
				"kåre",
				time.Now().Add(-time.Hour),
				"whatever",
			),
		))

		// Should be success
		resp := httptest.NewRecorder()
		it.Router.ServeHTTP(resp, req)
		respBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.Code, string(respBody))

		accUpdated, err := it.App.AccountRepo.GetAccountByWithingsUserID(it.Ctx, withingsUserID)
		require.NoError(t, err)
		require.Equal(t, "a075f8c14fb8df40b08ebc8508533dc332a6910a", accUpdated.WithingsAccessToken())
	})

	t.Run("without account on context returns bad request", func(t *testing.T) {
		beforeEach(t)

		req := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)

		// Should return bad request
		resp := httptest.NewRecorder()
		it.Router.ServeHTTP(resp, req)
		respBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 401, resp.Code, string(respBody), "should return bad request")
	})
}
