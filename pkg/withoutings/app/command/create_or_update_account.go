package command

import (
	"context"
	"errors"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type CreateOrUpdateAccount struct {
	Account *account.Account
}

type CreateOrUpdateAccountHandler interface {
	Handle(ctx context.Context, cmd CreateOrUpdateAccount) error
}

func (h createOrUpdateAccountHandler) Handle(ctx context.Context, cmd CreateOrUpdateAccount) (err error) {
	existingAcc, err := h.accountRepo.GetAccountByWithingsUserID(ctx, cmd.Account.WithingsUserID())
	if err != nil && !errors.Is(err, account.ErrAccountNotFound) {
		return err
	}
	accountAlreadyExists := err == nil

	if accountAlreadyExists {
		return h.accountRepo.Update(ctx, existingAcc,
			func(ctx context.Context, acc *account.Account) (*account.Account, error) {
				acc = cmd.Account
				return acc, nil
			})
	} else {
		return h.accountRepo.CreateAccount(ctx, cmd.Account)
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
