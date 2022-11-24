package command

import (
	"context"
	"github.com/roessland/withoutings/internal/withoutings/domain/account"
)

type CreateOrUpdateAccount struct {
	Account account.Account
}

type CreateOrUpdateAccountHandler interface {
	Handle(ctx context.Context, cmd CreateOrUpdateAccount) error
}

func (h createOrUpdateAccountHandler) Handle(ctx context.Context, cmd CreateOrUpdateAccount) (err error) {
	_, err = h.accountRepo.GetAccountByWithingsUserID(ctx, cmd.Account.WithingsUserID)

	return h.accountRepo.CreateAccount(ctx, account.Account{
		WithingsUserID:            cmd.Account.WithingsUserID,
		WithingsAccessToken:       cmd.Account.WithingsAccessToken,
		WithingsRefreshToken:      cmd.Account.WithingsRefreshToken,
		WithingsAccessTokenExpiry: cmd.Account.WithingsAccessTokenExpiry,
		WithingsScopes:            cmd.Account.WithingsScopes,
	})
}

func NewCreateOrUpdateAccountHandler(accountRepo account.Repository) CreateOrUpdateAccountHandler {
	return createOrUpdateAccountHandler{
		accountRepo: accountRepo,
	}
}

type createOrUpdateAccountHandler struct {
	accountRepo account.Repository
}
