package query

import (
	"context"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type AccountByID struct {
	AccountID int64
}

type AccountByIDHandler interface {
	Handle(ctx context.Context, query AccountByID) (account.Account, error)
}

type accountByIDHandler struct {
	readModel accountByIDReadModel
}

func NewAccountByIDHandler(
	readModel accountByIDReadModel,
) AccountByIDHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return accountByIDHandler{readModel: readModel}
}

type accountByIDReadModel interface {
	GetAccountByID(ctx context.Context, accountID int64) (account.Account, error)
}

func (h accountByIDHandler) Handle(ctx context.Context, query AccountByID) (account account.Account, err error) {
	return h.readModel.GetAccountByID(ctx, query.AccountID)
}
