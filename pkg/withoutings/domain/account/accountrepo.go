package account

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type NotFoundError struct {
	WithingsUserID string
	AccountUUID    uuid.UUID
}

func (e NotFoundError) Error() string {
	if e.WithingsUserID != "" {
		return fmt.Sprintf("account with Withings user ID '%s' not found", e.WithingsUserID)
	}
	if e.AccountUUID != uuid.Nil {
		return fmt.Sprintf("account with UUID '%s' not found", e.AccountUUID)
	}
	return fmt.Sprintf("account not found")

}

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
