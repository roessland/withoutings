package account_test

import (
	"context"
	"fmt"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/testdb"
	account2 "github.com/roessland/withoutings/pkg/withoutings/adapter/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

var _ account.Repo = account2.AccountPgRepo{}

func TestAccountPgRepo_UpdateAccount(t *testing.T) {
	ctx := testctx.New()
	database := testdb.New(ctx)
	defer database.Drop(ctx)
	queries := db.New(database)
	repo := account2.NewAccountPgRepo(database.Pool, queries)

	// Insert test user with default field values
	withingsUserID := fmt.Sprintf("%d", rand.Intn(10000))
	err := repo.CreateAccount(ctx, account.Account{
		WithingsUserID: withingsUserID,
	})
	require.NoError(t, err)

	// Retrieve inserted account ID
	acc, err := repo.GetAccountByWithingsUserID(ctx, withingsUserID)
	require.NoError(t, err)
	accountID := acc.AccountID

	t.Run("updates all fields", func(t *testing.T) {
		err := repo.UpdateAccount(ctx, accountID, func(ctx context.Context, accNext account.Account) (account.Account, error) {
			return account.Account{
				WithingsUserID:            "a",
				WithingsAccessToken:       "b",
				WithingsRefreshToken:      "c",
				WithingsAccessTokenExpiry: time.Now().Add(time.Minute),
				WithingsScopes:            "d",
			}, nil
		})
		require.NoError(t, err)

		accUpdated, err := repo.GetAccountByID(ctx, accountID)
		require.NoError(t, err)
		require.EqualValues(t, "a", accUpdated.WithingsUserID)
		require.EqualValues(t, "b", accUpdated.WithingsAccessToken)
		require.EqualValues(t, "c", accUpdated.WithingsRefreshToken)
		require.True(t, accUpdated.WithingsAccessTokenExpiry.After(time.Now()))
		require.EqualValues(t, "d", accUpdated.WithingsScopes)

	})
}
