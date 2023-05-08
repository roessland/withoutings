package command

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"time"
)

type RefreshAccessToken struct {
	Account   *account.Account
	CommandID uuid.UUID
}

type RefreshAccessTokenHandler interface {
	Handle(ctx context.Context, cmd RefreshAccessToken) error
}

func (h refreshAccessTokenHandler) Handle(ctx context.Context, cmd RefreshAccessToken) (err error) {
	log := logging.MustGetLoggerFromContext(ctx)
	acc := cmd.Account
	if acc == nil {
		return fmt.Errorf("account is nil")
	}
	if acc.WithingsAccessTokenExpiry().After(time.Now()) {
		return
	}

	newToken, err := h.withingsRepo.RefreshAccessToken(ctx, acc.WithingsRefreshToken())
	if err != nil {
		return err
	}

	return h.accountRepo.Update(
		ctx,
		acc,
		func(ctx context.Context, accNext *account.Account) (*account.Account, error) {
			if accNext.WithingsRefreshToken() != acc.WithingsRefreshToken() {
				log.Warn("someone else updated withingsRefreshToken already")
				return accNext, nil
			}
			err = accNext.UpdateWithingsToken(newToken.AccessToken, newToken.RefreshToken, newToken.Expiry, newToken.Scope)
			if err != nil {
				return nil, fmt.Errorf("failed to update withings token: %w", err)
			}
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
