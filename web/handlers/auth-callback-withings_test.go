package handlers_test

import (
	"fmt"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/roessland/withoutings/internal/repos/db"
	"github.com/roessland/withoutings/internal/services/withoutings"
	"github.com/roessland/withoutings/internal/testctx"
	"github.com/roessland/withoutings/internal/testdb"
	"github.com/roessland/withoutings/web"
	"github.com/roessland/withoutings/withingsapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCallback(t *testing.T) {
	ctx := testctx.New()
	database := testdb.New(ctx)
	defer database.Drop(ctx)

	mockWithingsTokenEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example response from docs. Not an actual token.
		fmt.Fprintln(w, `
		{
			"status": 0,
			"body": {
				"userid": "363",
				"access_token": "a075f8c14fb8df40b08ebc8508533dc332a6910a",
				"refresh_token": "f631236f02b991810feb774765b6ae8e6c6839ca",
				"expires_in": 10800,
				"scope": "user.info,user.metrics",
				"csrf_token": "PACnnxwHTaBQOzF7bQqwFUUotIuvtzSM",
				"token_type": "Bearer"
			}
		}`)
	}))

	defer mockWithingsTokenEndpoint.Close()

	svc := &withoutings.Service{}
	svc.Log = ctx.Logger
	svc.DB = database.Pool
	svc.Queries = db.New(svc.DB)
	svc.AccountRepo = svc.Queries
	svc.SubscriptionRepo = svc.Queries

	svc.Sessions = scs.New()
	svc.Sessions.Store = pgxstore.New(svc.DB)

	svc.Withings = withingsapi.NewClient("testclientid", "testclientsecret", "testredirecturl")
	svc.Withings.APIBase = mockWithingsTokenEndpoint.URL
	svc.Withings.OAuth2Config.Endpoint.TokenURL = mockWithingsTokenEndpoint.URL
	svc.Withings.OAuth2Config.Endpoint.AuthURL = mockWithingsTokenEndpoint.URL

	router := web.Router(svc)

	t.Run("without code yields bad request", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/callback", nil)
		router.ServeHTTP(resp, req)
		require.Equal(t, 400, resp.Code)
	})

	t.Run("without cookie yields bad request", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=qwerty", nil)

		router.ServeHTTP(resp, req)
		require.Equal(t, 400, resp.Code)
	})

	t.Run("with correct code and wrong state yields bad request", func(t *testing.T) {
		// Store state in session
		exampleDeadline := time.Now().Add(time.Hour)
		encodedValue, err := svc.Sessions.Codec.Encode(exampleDeadline, map[string]interface{}{
			"state": "e0GANQxF1SG",
		})
		require.NoError(t, err)
		err = svc.Sessions.Store.Commit("some-session-id", encodedValue, exampleDeadline)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=qwerty&state=WRONGSTATE", nil)

		// Add cookie with correct session_id, referring to session state stored earlier
		cookie := http.Cookie{Name: svc.Sessions.Cookie.Name, Value: "some-session-id"}
		req.AddCookie(&cookie)

		// Should be success and redirect
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, 400, resp.Code)

		accounts, err := svc.AccountRepo.ListAccounts(ctx)
		require.NoError(t, err)
		require.Len(t, accounts, 0)
	})

	t.Run("with correct code and state creates account", func(t *testing.T) {
		// Store state in session
		exampleDeadline := time.Now().Add(time.Hour)
		encodedValue, err := svc.Sessions.Codec.Encode(exampleDeadline, map[string]interface{}{
			"state": "e0GANQxF1SG",
		})
		require.NoError(t, err)
		err = svc.Sessions.Store.Commit("some-session-id", encodedValue, exampleDeadline)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=qwerty&state=e0GANQxF1SG", nil)

		// Add cookie with correct session_id, referring to session state stored earlier
		cookie := http.Cookie{Name: svc.Sessions.Cookie.Name, Value: "some-session-id"}
		req.AddCookie(&cookie)

		// Should be success and redirect
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, 302, resp.Code)

		accounts, err := svc.AccountRepo.ListAccounts(ctx)
		require.NoError(t, err)
		require.Len(t, accounts, 1)
		acc := accounts[0]
		assert.Equal(t, "363", acc.WithingsUserID)
		assert.Equal(t, "a075f8c14fb8df40b08ebc8508533dc332a6910a", acc.WithingsAccessToken)
		assert.Equal(t, "f631236f02b991810feb774765b6ae8e6c6839ca", acc.WithingsRefreshToken)
		assert.WithinDuration(t, time.Now().Add(10800*time.Second), acc.WithingsAccessTokenExpiry, time.Minute)
		assert.Equal(t, "user.info,user.metrics", acc.WithingsScopes)
	})

}
