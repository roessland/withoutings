package adapter

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/roessland/withoutings/internal/repos/db"
	"github.com/roessland/withoutings/internal/withoutings/domain/account"
)

type AccountPgRepo struct {
	queries *db.Queries
}

func NewAccountPgRepo(queries *db.Queries) AccountPgRepo {
	return AccountPgRepo{
		queries: queries,
	}
}

func (r AccountPgRepo) GetAccountByWithingsUserID(ctx context.Context, withingsUserID string) (account.Account, error) {
	acc, err := r.queries.GetAccountByWithingsUserID(ctx, withingsUserID)
	if err == pgx.ErrNoRows {
		return account.Account{}, account.NotFoundError{WithingsUserID: withingsUserID}
	}
	if err != nil {
		return account.Account{}, errors.Wrap(err, "unable to retrieve account")
	}
	return account.Account{
		AccountID:                 acc.AccountID,
		WithingsUserID:            acc.WithingsUserID,
		WithingsAccessToken:       acc.WithingsAccessToken,
		WithingsRefreshToken:      acc.WithingsRefreshToken,
		WithingsAccessTokenExpiry: acc.WithingsAccessTokenExpiry,
		WithingsScopes:            acc.WithingsScopes,
	}, err
}

func (r AccountPgRepo) CreateAccount(ctx context.Context, account account.Account) error {
	return r.queries.CreateAccount(ctx, db.CreateAccountParams{
		WithingsUserID:            account.WithingsUserID,
		WithingsAccessToken:       account.WithingsAccessToken,
		WithingsRefreshToken:      account.WithingsRefreshToken,
		WithingsAccessTokenExpiry: account.WithingsAccessTokenExpiry,
		WithingsScopes:            account.WithingsScopes,
	})
}
