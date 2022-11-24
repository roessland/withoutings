package adapters

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/roessland/withoutings/internal/repos/db"
	"github.com/roessland/withoutings/internal/withoutings/domain/account"
)

type AccountPostgresRepository struct {
	queries *db.Queries
}

func NewAccountPostgresRepository(queries *db.Queries) AccountPostgresRepository {
	return AccountPostgresRepository{
		queries: queries,
	}
}

func (r AccountPostgresRepository) GetAccountByWithingsUserID(ctx context.Context, withingsUserID string) (account.Account, error) {
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

func (r AccountPostgresRepository) CreateAccount(ctx context.Context, account account.Account) error {
	return r.queries.CreateAccount(ctx, db.CreateAccountParams{
		WithingsUserID:            account.WithingsUserID,
		WithingsAccessToken:       account.WithingsAccessToken,
		WithingsRefreshToken:      account.WithingsRefreshToken,
		WithingsAccessTokenExpiry: account.WithingsAccessTokenExpiry,
		WithingsScopes:            account.WithingsScopes,
	})
}
