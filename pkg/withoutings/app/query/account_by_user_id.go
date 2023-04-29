package query

import (
	"context"
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type AccountByUUID struct {
	AccountUUID uuid.UUID
}

type AccountByUUIDHandler interface {
	Handle(ctx context.Context, query AccountByUUID) (*account.Account, error)
}

type accountByUUIDHandler struct {
	readModel accountByUUIDReadModel
}

func NewAccountByUUIDHandler(
	readModel accountByUUIDReadModel,
) AccountByUUIDHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return accountByUUIDHandler{readModel: readModel}
}

type accountByUUIDReadModel interface {
	GetAccountByUUID(ctx context.Context, accountUUID uuid.UUID) (*account.Account, error)
}

func (h accountByUUIDHandler) Handle(ctx context.Context, query AccountByUUID) (*account.Account, error) {
	if query.AccountUUID == uuid.Nil {
		return nil, account.NotFoundError{}
	}
	return h.readModel.GetAccountByUUID(ctx, query.AccountUUID)
}
