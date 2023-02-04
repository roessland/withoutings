package command

import (
	"context"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"time"
)

type RefreshAccessToken struct {
	Account account.Account
}

type RefreshAccessTokenHandler interface {
	Handle(ctx context.Context, cmd RefreshAccessToken) error
}

func (h refreshAccessTokenHandler) Handle(ctx context.Context, cmd RefreshAccessToken) (err error) {
	log := logging.MustGetLoggerFromContext(ctx)
	acc := cmd.Account
	if acc.WithingsAccessTokenExpiry.After(time.Now()) {
		return
	}

	newToken, err := h.withingsRepo.RefreshAccessToken(ctx, acc.WithingsRefreshToken)
	if err != nil {
		return err
	}

	return h.accountRepo.UpdateAccount(
		ctx,
		acc.AccountID,
		func(ctx context.Context, accNext account.Account) (account.Account, error) {
			if accNext.WithingsRefreshToken != acc.WithingsRefreshToken {
				log.Warn("someone else updated WithingsRefreshToken already")
				return accNext, nil
			}
			accNext.WithingsAccessToken = newToken.AccessToken
			accNext.WithingsRefreshToken = newToken.RefreshToken
			accNext.WithingsAccessTokenExpiry = newToken.Expiry
			return accNext, nil
		},
	)
}

func NewRefreshAccessTokenHandler(accountRepo account.Repo, withingsRepo withings.Repo) RefreshAccessTokenHandler {
	return refreshAccessTokenHandler{
		accountRepo:  accountRepo,
		withingsRepo: withingsRepo,
	}
}

type refreshAccessTokenHandler struct {
	accountRepo  account.Repo
	withingsRepo withings.Repo
}
