package query

import (
	"context"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type AccountByWithingsUserID struct {
	WithingsUserID string
}

type AccountByWithingsUserIDHandler interface {
	Handle(ctx context.Context, query AccountByWithingsUserID) (account.Account, error)
}

type accountByWithingsUserIDHandler struct {
	readModel accountByWithingsUserIDReadModel
}

func NewAccountByWithingsUserIDHandler(
	readModel accountByWithingsUserIDReadModel,
) AccountByWithingsUserIDHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return accountByWithingsUserIDHandler{readModel: readModel}
}

type accountByWithingsUserIDReadModel interface {
	GetAccountByWithingsUserID(ctx context.Context, withingsUserID string) (account.Account, error)
}

func (h accountByWithingsUserIDHandler) Handle(ctx context.Context, query AccountByWithingsUserID) (account account.Account, err error) {
	return h.readModel.GetAccountByWithingsUserID(ctx, query.WithingsUserID)
}
