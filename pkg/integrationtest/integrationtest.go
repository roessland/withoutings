package integrationtest

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/testdb"
	"github.com/roessland/withoutings/pkg/web"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/worker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type IntegrationTest struct {
	Ctx      context.Context
	Logger   *logrus.Logger
	App      *app.App
	Mocks    *app.MockApp
	Router   *mux.Router
	Database testdb.TestDatabase
	Worker   *worker.Worker
}

// WithFreshDatabase returns a new Context.
func WithFreshDatabase(t *testing.T) IntegrationTest {
	ctx := testctx.New()
	database := testdb.New(ctx)
	t.Cleanup(func() {
		database.Drop(ctx)
	})

	it := IntegrationTest{
		Ctx:      ctx,
		Logger:   ctx.Logger,
		Database: database,
	}

	it.ResetMocks(t)
	return it
}

// ResetMocks resets all mocks in the test.
func (it *IntegrationTest) ResetMocks(t *testing.T) {
	mockApp := app.NewTestApplication(t, it.Ctx, it.Database.Pool)
	it.App = mockApp.App
	it.Mocks = mockApp
	it.Router = web.Router(it.App)
	it.Worker = worker.NewWorker(it.App)
}

func (it *IntegrationTest) MakeNewAccount(t *testing.T) *account.Account {
	t.Helper()
	// Insert a user with an expired withings access token
	accountUUID := uuid.New()
	withingsUserID := uuid.NewString()
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
	return acc
}

// DoRequest does a request against the test server.
func (it *IntegrationTest) DoRequest(req *http.Request) (*httptest.ResponseRecorder, string) {
	resp := httptest.NewRecorder()
	it.Router.ServeHTTP(resp, req)
	respBody, _ := io.ReadAll(resp.Body)
	return resp, string(respBody)
}
