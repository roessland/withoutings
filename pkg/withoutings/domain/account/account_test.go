package account_test

import (
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAccount_CanRefreshAccessToken_CannotRefreshNilAccount(t *testing.T) {
	var acc *account.Account
	require.False(t, acc.CanRefreshAccessToken())
}

func TestAccount_UpdateWithingsToken(t *testing.T) {
	acc := newDummyAccount()
	require.True(t, acc.CanRefreshAccessToken())

	err := acc.UpdateWithingsToken("new-access-token", "new-refresh-token", time.Now().Add(time.Hour), "new-scopes")
	require.NoError(t, err)
	
	require.Equal(t, "new-access-token", acc.WithingsAccessToken())
	require.Equal(t, "new-refresh-token", acc.WithingsRefreshToken())
	require.Equal(t, "new-scopes", acc.WithingsScopes())
	require.False(t, acc.CanRefreshAccessToken())
}

func newDummyAccount() *account.Account {
	acc, err := account.NewAccount(
		uuid.New(),
		"withings-user-id",
		"access-token",
		"refresh-token",
		time.Now().Add(-time.Hour), // expired
		"scopes",
	)
	if err != nil {
		panic(err)
	}
	return acc
}
