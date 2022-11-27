package account

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	WithingsUserID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("account with Withings user ID '%s' not found", e.WithingsUserID)
}

type Repo interface {
	GetAccountByWithingsUserID(ctx context.Context, withingsUserID string) (Account, error)
	CreateAccount(ctx context.Context, account Account) error
}
