package account_test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/testdb"
	accountadapter "github.com/roessland/withoutings/pkg/withoutings/adapter/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

var _ account.Repo = accountadapter.PgRepo{}

func TestAccountPgRepo_UpdateAccount(t *testing.T) {
	ctx := testctx.New()
	database := testdb.New(ctx)
	defer database.Drop(ctx)
	queries := db.New(database)
	repo := accountadapter.NewPgRepo(database.Pool, queries)

	// Insert test user with default field values
	withingsUserID := fmt.Sprintf("%d", rand.Intn(10000))
	err := repo.CreateAccount(ctx, account.NewAccountOrPanic(
		uuid.New(),
		withingsUserID,
		"gibberish",
		"whatever",
		time.Now().Add(time.Hour),
		"some_scope",
	))
	require.NoError(t, err)

	// Retrieve inserted account ID
	acc, err := repo.GetAccountByWithingsUserID(ctx, withingsUserID)
	require.NoError(t, err)
	withingsID := acc.WithingsUserID()

	t.Run("updates all fields", func(t *testing.T) {
		err := repo.UpdateAccount(ctx, withingsID, func(ctx context.Context, accNext *account.Account) (*account.Account, error) {
			require.NoError(t, accNext.UpdateWithingsToken(
				"gibberish-updated",
				"whatever-updated",
				time.Now().Add(time.Minute),
				"some_scope-updated",
			))
			return accNext, nil
		})
		require.NoError(t, err)

		accUpdated, err := repo.GetAccountByWithingsUserID(ctx, withingsID)
		require.NoError(t, err)
		require.EqualValues(t, withingsUserID, accUpdated.WithingsUserID())
		require.EqualValues(t, "gibberish-updated", accUpdated.WithingsAccessToken())
		require.EqualValues(t, "whatever-updated", accUpdated.WithingsRefreshToken())
		require.True(t, accUpdated.WithingsAccessTokenExpiry().After(time.Now()))
		require.EqualValues(t, "some_scope-updated", accUpdated.WithingsScopes())

	})
}
