package integrationtest

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/testdb"
	"github.com/roessland/withoutings/pkg/web"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/sirupsen/logrus"
	"testing"
)

type IntegrationTest struct {
	Ctx      context.Context
	Logger   *logrus.Logger
	App      *app.App
	Mocks    *app.MockApp
	Router   *mux.Router
	Database testdb.TestDatabase
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
}
