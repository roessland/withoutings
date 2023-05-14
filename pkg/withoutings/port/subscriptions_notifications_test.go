package port_test

import (
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/integrationtest"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNotificationsPage(t *testing.T) {
	// TODO this is not really an integration test, since it doesn't use the database.
	// It should be converted to a unit test.
	it := integrationtest.WithFreshDatabase(t)

	var accountUUID uuid.UUID
	var withingsUserID string

	loggedInRequest := func() *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/notifications", nil)

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

	// TODO move into integrationtest package
	doRequest := func(req *http.Request) (*httptest.ResponseRecorder, string) {
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

	t.Run("should respond with 200 OK", func(t *testing.T) {
		beforeEach(t)
		req := loggedInRequest()

		resp, body := doRequest(req)
		assert.Equal(t, 200, resp.Code, body)
	})
}
