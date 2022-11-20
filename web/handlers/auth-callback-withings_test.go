package handlers_test

import (
	"fmt"
	"github.com/roessland/withoutings/internal/domain/services/withoutings"
	"github.com/roessland/withoutings/internal/repos/db"
	"github.com/roessland/withoutings/internal/testctx"
	"github.com/roessland/withoutings/internal/testdb"
	"github.com/roessland/withoutings/web"
	"github.com/roessland/withoutings/web/sessions"
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
	svc.Sessions = sessions.NewManager([]byte("abc123"))
	svc.Withings = withingsapi.NewClient("testclientid", "testclientsecret", "testredirecturl")
	svc.Withings.APIBase = mockWithingsTokenEndpoint.URL
	svc.Withings.OAuth2Config.Endpoint.TokenURL = mockWithingsTokenEndpoint.URL
	svc.Withings.OAuth2Config.Endpoint.AuthURL = mockWithingsTokenEndpoint.URL
	svc.DB = database.Pool
	svc.AccountRepo = db.New(svc.DB)

	router := web.Router(svc)

	t.Run("without code yields bad request", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/callback", nil)
		router.ServeHTTP(resp, req)
		require.Equal(t, 400, resp.Code)
	})

	t.Run("without state yields bad request", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=qwerty", nil)
		router.ServeHTTP(resp, req)
		require.Equal(t, 400, resp.Code)
	})

	t.Run("with valid code creates user", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=qwerty&state=asdf", nil)
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
