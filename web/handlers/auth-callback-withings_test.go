package handlers_test

import (
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/config"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/testdb"
	accountAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/account"
	subscriptionAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/roessland/withoutings/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	svc := &app.App{}
	svc.Log = ctx.Logger
	queries := db.New(database)

	var router *mux.Router

	var mockWithingsRepo *withings.MockRepo

	beforeEach := func(t *testing.T) {
		var accountRepo account.Repo = accountAdapter.NewAccountPgRepo(database.Pool, queries)
		svc.AccountRepo = accountRepo

		var subscriptionRepo subscription.Repo = subscriptionAdapter.NewSubscriptionPgRepo(database.Pool, queries)
		svc.SubscriptionRepo = subscriptionRepo

		mockWithingsRepo = withings.NewMockRepo(t)
		svc.WithingsRepo = mockWithingsRepo

		svc.Queries = app.Queries{
			AccountForUserID:         query.NewAccountByIDHandler(accountRepo),
			AccountForWithingsUserID: query.NewAccountByWithingsUserIDHandler(accountRepo),
			AllAccounts:              query.NewAllAccountsHandler(accountRepo),
		}
		svc.Queries.Validate()

		cfg := &config.Config{}

		svc.Commands = app.Commands{
			// TODO replace withingsClient with interface
			SubscribeAccount:      command.NewSubscribeAccountHandler(accountRepo, subscriptionRepo, mockWithingsRepo, cfg),
			CreateOrUpdateAccount: command.NewCreateOrUpdateAccountHandler(accountRepo),
		}

		svc.Sessions = scs.New()
		svc.Sessions.Lifetime = time.Hour * 3
		svc.Sessions.IdleTimeout = time.Hour * 4
		svc.Sessions.Store = pgxstore.New(database.Pool)

		router = web.Router(svc)
	}

	t.Run("without code yields bad request", func(t *testing.T) {
		beforeEach(t)
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/callback", nil)
		router.ServeHTTP(resp, req)
		require.Equal(t, 400, resp.Code)
	})

	t.Run("without cookie yields bad request", func(t *testing.T) {
		beforeEach(t)

		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=qwerty", nil)

		router.ServeHTTP(resp, req)
		require.Equal(t, 400, resp.Code)
	})

	t.Run("with correct code and wrong state yields bad request", func(t *testing.T) {
		beforeEach(t)

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

		accounts, err := svc.Queries.AllAccounts.Handle(ctx, query.AllAccounts{})
		require.NoError(t, err)
		require.Len(t, accounts, 0)
	})

	t.Run("with correct code and state creates account", func(t *testing.T) {
		beforeEach(t)

		mockWithingsRepo.EXPECT().
			GetAccessToken(mock.Anything, mock.Anything).
			Return(&withings.Token{
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
