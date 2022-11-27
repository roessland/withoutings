package adapter_test

import (
	"github.com/roessland/withoutings/internal/withoutings/adapter"
	"github.com/roessland/withoutings/internal/withoutings/domain/account"
)

var _ account.Repo = adapter.AccountPgRepo{}
