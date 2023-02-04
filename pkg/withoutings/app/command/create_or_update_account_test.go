package command_test

import (
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/testdb"
	accountAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/account"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateOrUpdateAccountHandler(t *testing.T) {
	ctx := testctx.New()
	database := testdb.New(ctx)
	defer database.Drop(ctx)

	queries := db.New(database)

	var accountRepo account.Repo = accountAdapter.NewPgRepo(database.Pool, queries)

	var createOrUpdateAccount = command.NewCreateOrUpdateAccountHandler(accountRepo)

	now := time.Now()
	inOneHour := now.Add(time.Hour)

	t.Run("creates if no account exists", func(t *testing.T) {
		err := createOrUpdateAccount.Handle(ctx, command.CreateOrUpdateAccount{Account: account.Account{
			WithingsUserID:            "1",
			WithingsAccessToken:       "super",
			WithingsRefreshToken:      "secret",
			WithingsAccessTokenExpiry: now,
			WithingsScopes:            "user.info",
		}})
		require.NoError(t, err)

		acc1, err := accountRepo.GetAccountByWithingsUserID(ctx, "1")
		require.Positive(t, acc1.AccountID)
		require.EqualValues(t, "1", acc1.WithingsUserID)
		require.EqualValues(t, "super", acc1.WithingsAccessToken)
		require.EqualValues(t, "secret", acc1.WithingsRefreshToken)
		require.WithinDuration(t, now, acc1.WithingsAccessTokenExpiry, time.Second)
		require.EqualValues(t, "user.info", acc1.WithingsScopes)
	})

	t.Run("updates if account already exists", func(t *testing.T) {
		err := createOrUpdateAccount.Handle(ctx, command.CreateOrUpdateAccount{Account: account.Account{
			WithingsUserID:            "2",
			WithingsAccessToken:       "super",
			WithingsRefreshToken:      "secret",
			WithingsAccessTokenExpiry: now,
			WithingsScopes:            "user.info",
		}})
		require.NoError(t, err)

		err = createOrUpdateAccount.Handle(ctx, command.CreateOrUpdateAccount{Account: account.Account{
			WithingsUserID:            "2",
			WithingsAccessToken:       "SUPER",
			WithingsRefreshToken:      "SECRET",
			WithingsAccessTokenExpiry: inOneHour,
		}})
		require.NoError(t, err)

		acc2, err := accountRepo.GetAccountByWithingsUserID(ctx, "2")
		require.Positive(t, acc2.AccountID)
		require.EqualValues(t, "2", acc2.WithingsUserID)
		require.EqualValues(t, "SUPER", acc2.WithingsAccessToken)
		require.EqualValues(t, "SECRET", acc2.WithingsRefreshToken)
		require.WithinDuration(t, inOneHour, acc2.WithingsAccessTokenExpiry, time.Second)
		require.EqualValues(t, "", acc2.WithingsScopes) // it was missing, so it was zeroed

	})
}
