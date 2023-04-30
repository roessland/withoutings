package account

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

var ErrAccountNotFound = errors.New("account not found")

//go:generate mockery --name Repo --filename accountrepo_mock.go
type Repo interface {
	GetAccountByWithingsUserID(ctx context.Context, withingsUserID string) (*Account, error)
	GetAccountByUUID(ctx context.Context, accountUUID uuid.UUID) (*Account, error)
	CreateAccount(ctx context.Context, account *Account) error
	ListAccounts(ctx context.Context) ([]*Account, error)
	UpdateAccount(
		ctx context.Context,
		withingsUserID string,
		updateFn func(ctx context.Context, acc *Account) (*Account, error),
	) error
}
