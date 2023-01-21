package command

import (
	"context"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type SubscribeAccount struct {
	Account account.Account
}

type SubscribeAccountHandler interface {
	Handle(ctx context.Context, cmd SubscribeAccount) error
}

func (h subscribeAccountHandler) Handle(ctx context.Context, cmd SubscribeAccount) (err error) {
	_, err = h.accountRepo.GetAccountByWithingsUserID(ctx, cmd.Account.WithingsUserID)

	return h.accountRepo.CreateAccount(ctx, account.Account{
		WithingsUserID:            cmd.Account.WithingsUserID,
		WithingsAccessToken:       cmd.Account.WithingsAccessToken,
		WithingsRefreshToken:      cmd.Account.WithingsRefreshToken,
		WithingsAccessTokenExpiry: cmd.Account.WithingsAccessTokenExpiry,
		WithingsScopes:            cmd.Account.WithingsScopes,
	})
}

func NewSubscribeAccountHandler(accountRepo account.Repo) SubscribeAccountHandler {
	return subscribeAccountHandler{
		accountRepo: accountRepo,
	}
}

type subscribeAccountHandler struct {
	accountRepo account.Repo
}
