package db_test

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/roessland/withoutings/internal/repos/db"
	"github.com/roessland/withoutings/internal/testctx"
	"github.com/roessland/withoutings/internal/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAccountRepo(t *testing.T) {
	ctx := testctx.New()
	database := testdb.New(ctx)
	defer database.Drop(ctx)

	queries := db.New(database)

	t.Run("CreateOrUpdateAccount", func(t *testing.T) {
		accessTokenExpiry := time.Now().Add(time.Hour)
		createAccountParams := db.CreateAccountParams{
			WithingsUserID:            "userid",
			WithingsAccessToken:       "accesstoken",
			WithingsRefreshToken:      "refreshtoken",
			WithingsAccessTokenExpiry: accessTokenExpiry,
			WithingsScopes:            "scope1,scope2,scope3",
		}
		account, err := queries.CreateOrUpdateAccount(ctx, createAccountParams)
		require.NoError(t, err)

		assert.True(t, account.AccountID > 0)
		assert.Equal(t, "userid", account.WithingsUserID)
		assert.Equal(t, "accesstoken", account.WithingsAccessToken)
		assert.Equal(t, "refreshtoken", account.WithingsRefreshToken)
		assert.Equal(t, accessTokenExpiry.Truncate(time.Second), account.WithingsAccessTokenExpiry.Truncate(time.Second))
		assert.Equal(t, "scope1,scope2,scope3", account.WithingsScopes)
	})
}
