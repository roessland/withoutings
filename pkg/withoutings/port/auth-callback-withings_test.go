package port_test

import (
	"github.com/roessland/withoutings/pkg/integrationtest"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCallback(t *testing.T) {
	it := integrationtest.WithFreshDatabase(t)

	beforeEach := func(t *testing.T) {
		it.ResetMocks(t)
	}

	t.Run("without code yields bad request", func(t *testing.T) {
		beforeEach(t)
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/callback", nil)
		it.Router.ServeHTTP(resp, req)
		require.Equal(t, 400, resp.Code)
	})

	t.Run("without cookie yields bad request", func(t *testing.T) {
		beforeEach(t)

		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=qwerty", nil)

		it.Router.ServeHTTP(resp, req)
		require.Equal(t, 400, resp.Code)
	})

	t.Run("with correct code and wrong state yields bad request", func(t *testing.T) {
		beforeEach(t)

		// Store state in session
		exampleDeadline := time.Now().Add(time.Hour)
		encodedValue, err := it.App.Sessions.Codec.Encode(exampleDeadline, map[string]interface{}{
			"state": "e0GANQxF1SG",
		})
		require.NoError(t, err)
		err = it.App.Sessions.Store.Commit("some-session-id1", encodedValue, exampleDeadline)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=qwerty&state=WRONGSTATE", nil)

		// Add cookie with correct session_id, referring to session state stored earlier
		cookie := http.Cookie{Name: it.App.Sessions.Cookie.Name, Value: "some-session-id1"}
		req.AddCookie(&cookie)

		// Should be success and redirect
		resp := httptest.NewRecorder()
		it.Router.ServeHTTP(resp, req)
		assert.Equal(t, 400, resp.Code)

		accounts, err := it.App.Queries.AllAccounts.Handle(it.Ctx, query.AllAccounts{})
		require.NoError(t, err)
		require.Len(t, accounts, 0)
	})

	t.Run("with correct code and state creates account", func(t *testing.T) {
		beforeEach(t)

		it.Mocks.MockWithingsRepo.EXPECT().
			GetAccessToken(mock.Anything, mock.Anything).
			Once().Return(&withings.Token{
			UserID:       "363",
			AccessToken:  "a075f8c14fb8df40b08ebc8508533dc332a6910a",
			RefreshToken: "f631236f02b991810feb774765b6ae8e6c6839ca",
			ExpiresIn:    10800,
			Scope:        "user.info,user.metrics",
			CSRFToken:    "PACnnxwHTaBQOzF7bQqwFUUotIuvtzSM",
			TokenType:    "Bearer",
			Expiry:       time.Now().Add(10800 * time.Second),
		}, nil)

		// Store state in session
		exampleDeadline := time.Now().Add(3 * time.Hour)
		encodedValue, err := it.App.Sessions.Codec.Encode(exampleDeadline, map[string]interface{}{
			"state": "e0GANQxF1SG",
		})
		require.NoError(t, err)
		err = it.App.Sessions.Store.Commit("some-session-id2", encodedValue, exampleDeadline)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=qwerty&state=e0GANQxF1SG", nil)

		// Add cookie with correct session_id, referring to session state stored earlier
		cookie := http.Cookie{Name: it.App.Sessions.Cookie.Name, Value: "some-session-id2"}
		req.AddCookie(&cookie)

		// Should be success and redirect
		resp := httptest.NewRecorder()
		it.Router.ServeHTTP(resp, req)
		assert.Equal(t, 302, resp.Code)

		accounts, err := it.App.AccountRepo.ListAccounts(it.Ctx)
		require.NoError(t, err)
		require.Len(t, accounts, 1)
		acc := accounts[0]
		assert.Equal(t, "363", acc.WithingsUserID())
		assert.Equal(t, "a075f8c14fb8df40b08ebc8508533dc332a6910a", acc.WithingsAccessToken())
		assert.Equal(t, "f631236f02b991810feb774765b6ae8e6c6839ca", acc.WithingsRefreshToken())
		assert.WithinDuration(t, time.Now().Add(10800*time.Second), acc.WithingsAccessTokenExpiry(), time.Minute)
		assert.Equal(t, "user.info,user.metrics", acc.WithingsScopes())
	})

}
