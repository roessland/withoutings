package adapter_test

import (
	"github.com/roessland/withoutings/pkg/withoutings/adapter"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

var _ account.Repo = adapter.AccountPgRepo{}
