package port_test

import (
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/integrationtest"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"html"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWithingsMeasureGetmeasPage(t *testing.T) {
	it := integrationtest.WithFreshDatabase(t)

	var accountUUID uuid.UUID
	var withingsUserID string

	loggedInRequest := func() *http.Request {
		req := httptest.NewRequest(http.MethodPost, "/withings/measure/getmeas", nil)

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

	beforeEach := func(t *testing.T) {
		it.ResetMocks(t)

		// Insert a user with a valid withings access token
		accountUUID = uuid.New()
		withingsUserID = uuid.NewString()
		acc, err := account.NewAccount(
			accountUUID,
			withingsUserID,
			"bob",
			"kåre",
			time.Now().Add(time.Hour),
			"whatever",
		)
		require.NoError(t, err)
		require.NoError(t, it.App.AccountRepo.CreateAccount(it.Ctx, acc))

	}

	t.Run("should respond with 200 OK", func(t *testing.T) {
		beforeEach(t)
		it.Mocks.MockWithingsSvc.EXPECT().MeasureGetmeas(mock.Anything, mock.Anything, mock.Anything).Return(&withings.MeasureGetmeasResponse{
			Status: 0,
			Body:   withings.MeasureGetmeasBody{},
			Raw:    []byte(`{ "status": 0, "body": {} }`),
		}, nil)

		req := loggedInRequest()

		resp, body := it.DoRequest(req)
		assert.Equal(t, 200, resp.Code, body)
		require.Contains(t, body, html.EscapeString(`{ "status": 0, "body": {} }`))
	})
}
