package query

import (
	"context"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type AllAccounts struct {
}

type AllAccountsHandler interface {
	Handle(ctx context.Context, query AllAccounts) ([]*account.Account, error)
}

type accountsHandler struct {
	readModel allAccountsReadModel
}

func NewAllAccountsHandler(
	readModel allAccountsReadModel,
) AllAccountsHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return accountsHandler{readModel: readModel}
}

type allAccountsReadModel interface {
	ListAccounts(ctx context.Context) ([]*account.Account, error)
}

func (h accountsHandler) Handle(ctx context.Context, query AllAccounts) (accounts []*account.Account, err error) {
	return h.readModel.ListAccounts(ctx)
}
