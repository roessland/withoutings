package account_test

import (
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccount_CanRefreshAccessToken_CannotRefreshNilAccount(t *testing.T) {
	var acc *account.Account
	require.False(t, acc.CanRefreshAccessToken())
}
