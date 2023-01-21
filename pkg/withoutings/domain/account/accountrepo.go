package account

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	WithingsUserID string
	AccountID      string
}

func (e NotFoundError) Error() string {
	if e.WithingsUserID != "" {
		return fmt.Sprintf("account with Withings user ID '%s' not found", e.WithingsUserID)
	}
	if e.AccountID != "" {
		return fmt.Sprintf("account with ID '%s' not found", e.AccountID)
	}
	return fmt.Sprintf("account not found")

}

type Repo interface {
	GetAccountByID(ctx context.Context, accountID int64) (Account, error)
	GetAccountByWithingsUserID(ctx context.Context, withingsUserID string) (Account, error)
	CreateAccount(ctx context.Context, account Account) error
	ListAccounts(ctx context.Context) ([]Account, error)
}
