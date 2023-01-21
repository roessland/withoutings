package command

import (
	"context"
	"errors"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type CreateOrUpdateAccount struct {
	Account account.Account
}

type CreateOrUpdateAccountHandler interface {
	Handle(ctx context.Context, cmd CreateOrUpdateAccount) error
}

func (h createOrUpdateAccountHandler) Handle(ctx context.Context, cmd CreateOrUpdateAccount) (err error) {
	_, err = h.accountRepo.GetAccountByWithingsUserID(ctx, cmd.Account.WithingsUserID)
	if err != nil && errors.Is(err, account.NotFoundError{}) {
		return err
	}
	accountAlreadyExists := err == nil

	if accountAlreadyExists {
		return nil
	}

	return h.accountRepo.CreateAccount(ctx, account.Account{
		WithingsUserID:            cmd.Account.WithingsUserID,
		WithingsAccessToken:       cmd.Account.WithingsAccessToken,
		WithingsRefreshToken:      cmd.Account.WithingsRefreshToken,
		WithingsAccessTokenExpiry: cmd.Account.WithingsAccessTokenExpiry,
		WithingsScopes:            cmd.Account.WithingsScopes,
	})
}

func NewCreateOrUpdateAccountHandler(accountRepo account.Repo) CreateOrUpdateAccountHandler {
	return createOrUpdateAccountHandler{
		accountRepo: accountRepo,
	}
}

type createOrUpdateAccountHandler struct {
	accountRepo account.Repo
}
