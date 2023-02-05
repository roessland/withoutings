package handlers_test

import (
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/testdb"
	accountAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/account"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/roessland/withoutings/web"
	"github.com/roessland/withoutings/web/middleware"
	"github.com/roessland/withoutings/web/templates"
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
	ctx := testctx.New()
	database := testdb.New(ctx)
	defer database.Drop(ctx)

	svc := &app.App{}
	svc.Log = ctx.Logger
	queries := db.New(database)

	var router *mux.Router

	var mockWithingsRepo *withings.MockRepo

	var accountID int64
	var withingsUserID string
	var accountRepo account.Repo

	beforeEach := func(t *testing.T) {
		accountRepo = accountAdapter.NewPgRepo(database.Pool, queries)
		svc.AccountRepo = accountRepo
		withingsUserID = uuid.NewString()
		require.NoError(t, accountRepo.CreateAccount(ctx, account.Account{
			WithingsUserID:       withingsUserID,
			WithingsAccessToken:  "bob",
			WithingsRefreshToken: "kåre",
		}))
		acc, err := accountRepo.GetAccountByWithingsUserID(ctx, withingsUserID)
		require.NoError(t, err)
		accountID = acc.AccountID

		mockWithingsRepo = withings.NewMockRepo(t)
		svc.WithingsRepo = mockWithingsRepo

		svc.Queries = app.Queries{
			AccountForUserID:         query.NewAccountByIDHandler(accountRepo),
			AccountForWithingsUserID: query.NewAccountByWithingsUserIDHandler(accountRepo),
			AllAccounts:              query.NewAllAccountsHandler(accountRepo),
		}

		svc.Commands = app.Commands{
			RefreshAccessToken: command.NewRefreshAccessTokenHandler(accountRepo, mockWithingsRepo),
		}

		svc.Templates = templates.LoadTemplates()

		svc.Sessions = scs.New()
		svc.Sessions.Store = pgxstore.New(database.Pool)

		router = web.Router(svc)
	}

	t.Run("with expired token refreshes token", func(t *testing.T) {
		beforeEach(t)

		mockWithingsRepo.EXPECT().
			RefreshAccessToken(mock.Anything, mock.Anything).
			Return(&withings.Token{
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
		req = req.WithContext(middleware.AddAccountToContext(ctx, account.Account{
			AccountID:            accountID,
			WithingsUserID:       withingsUserID,
			WithingsAccessToken:  "bob",
			WithingsRefreshToken: "kåre",
		}))

		// Should be success
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		respBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.Code, string(respBody))

		accUpdated, err := accountRepo.GetAccountByWithingsUserID(ctx, withingsUserID)
		require.NoError(t, err)
		require.Equal(t, "a075f8c14fb8df40b08ebc8508533dc332a6910a", accUpdated.WithingsAccessToken)
	})

}
