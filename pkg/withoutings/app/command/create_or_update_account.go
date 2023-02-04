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
	existingAcc, err := h.accountRepo.GetAccountByWithingsUserID(ctx, cmd.Account.WithingsUserID)
	if err != nil && !errors.As(err, &account.NotFoundError{}) {
		return err
	}
	accountAlreadyExists := err == nil

	if accountAlreadyExists {
		return h.accountRepo.UpdateAccount(ctx, existingAcc.AccountID,
			func(ctx context.Context, acc account.Account) (account.Account, error) {
				if acc.WithingsUserID != cmd.Account.WithingsUserID {
					return account.Account{}, errors.New("tried to change withings user ID in createOrUpdateAccountHandler")
				}
				acc.WithingsAccessToken = cmd.Account.WithingsAccessToken
				acc.WithingsRefreshToken = cmd.Account.WithingsRefreshToken
				acc.WithingsAccessTokenExpiry = cmd.Account.WithingsAccessTokenExpiry
				acc.WithingsScopes = cmd.Account.WithingsScopes
				return acc, nil
			})
	} else {
		return h.accountRepo.CreateAccount(ctx, account.Account{
			WithingsUserID:            cmd.Account.WithingsUserID,
			WithingsAccessToken:       cmd.Account.WithingsAccessToken,
			WithingsRefreshToken:      cmd.Account.WithingsRefreshToken,
			WithingsAccessTokenExpiry: cmd.Account.WithingsAccessTokenExpiry,
			WithingsScopes:            cmd.Account.WithingsScopes,
		})
	}
}

func NewCreateOrUpdateAccountHandler(accountRepo account.Repo) CreateOrUpdateAccountHandler {
	return createOrUpdateAccountHandler{
		accountRepo: accountRepo,
	}
}

type createOrUpdateAccountHandler struct {
	accountRepo account.Repo
}
