package db_test

import (
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAccountQueries(t *testing.T) {
	ctx := testctx.New()
	database := testdb.New(ctx)
	defer database.Drop(ctx)

	queries := db.New(database)

	t.Run("CreateAccount, getByID, getByUUID", func(t *testing.T) {
		params := db.CreateAccountParams{
			AccountUuid:               uuid.New(),
			WithingsUserID:            "userid1337",
			WithingsAccessToken:       "accesstoken",
			WithingsRefreshToken:      "refreshtoken",
			WithingsAccessTokenExpiry: time.Now().Add(time.Hour),
			WithingsScopes:            "scope1,scope2,scope3",
		}
		err := queries.CreateAccount(ctx, params)
		require.NoError(t, err)

		account, err := queries.GetAccountByWithingsUserID(ctx, "userid1337")
		require.NoError(t, err)
		assert.True(t, account.AccountID > 0)
		assert.Equal(t, params.AccountUuid, account.AccountUuid)
		assert.Equal(t, params.WithingsUserID, account.WithingsUserID)
		assert.Equal(t, params.WithingsAccessToken, account.WithingsAccessToken)
		assert.Equal(t, params.WithingsRefreshToken, account.WithingsRefreshToken)
		assert.Equal(t,
			params.WithingsAccessTokenExpiry.Truncate(time.Second),
			account.WithingsAccessTokenExpiry.Truncate(time.Second))
		assert.Equal(t, params.WithingsScopes, account.WithingsScopes)

		account2, err := queries.GetAccountByAccountUUID(ctx, params.AccountUuid)
		require.NoError(t, err)
		require.Equal(t, account.AccountUuid, account2.AccountUuid)
	})
}
