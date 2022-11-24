package adapters_test

import (
	"github.com/roessland/withoutings/internal/withoutings/adapters"
	"github.com/roessland/withoutings/internal/withoutings/domain/account"
)

var _ account.Repository = adapters.AccountPostgresRepository{}
