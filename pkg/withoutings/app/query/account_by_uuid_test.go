package query_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/stretchr/testify/require"
	"testing"
)

type getAccountByUUIDReadModelMock struct{}

func (readModel getAccountByUUIDReadModelMock) GetAccountByUUID(ctx context.Context, accountUUID uuid.UUID) (*account.Account, error) {
	panic("should not be called")
}

func TestAccountByUUIDHandler_NilUUIDNotFound(t *testing.T) {
	readModel := getAccountByUUIDReadModelMock{}
	handler := query.NewAccountByUUIDHandler(readModel)
	acc, err := handler.Handle(context.Background(), query.AccountByUUID{
		AccountUUID: uuid.Nil,
	})
	require.Nil(t, acc)
	require.Equal(t, account.ErrAccountNotFound, err)
}
