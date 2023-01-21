package query

import (
	"context"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type Accounts struct {
}

type AccountsHandler interface {
	Handle(ctx context.Context, query Accounts) ([]account.Account, error)
}

type accountsHandler struct {
	readModel accountsReadModel
}

func NewAccountsHandler(
	readModel accountsReadModel,
) AccountsHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return accountsHandler{readModel: readModel}
}

type accountsReadModel interface {
	ListAccounts(ctx context.Context) ([]account.Account, error)
}

func (h accountsHandler) Handle(ctx context.Context, query Accounts) (accounts []account.Account, err error) {
	return h.readModel.ListAccounts(ctx)
}
