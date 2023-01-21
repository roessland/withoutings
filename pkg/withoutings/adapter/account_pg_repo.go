package adapter

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/roessland/withoutings/pkg/repos/db"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type AccountPgRepo struct {
	queries *db.Queries
}

func NewAccountPgRepo(queries *db.Queries) AccountPgRepo {
	return AccountPgRepo{
		queries: queries,
	}
}

func (r AccountPgRepo) GetAccountByID(ctx context.Context, accountID int64) (account.Account, error) {
	acc, err := r.queries.GetAccountByID(ctx, accountID)
	if err == pgx.ErrNoRows {
		return account.Account{}, account.NotFoundError{}
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

func (r AccountPgRepo) GetAccountByWithingsUserID(ctx context.Context, withingsUserID string) (account.Account, error) {
	if r.queries == nil {
		panic("queries was nil")
	}
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

func (r AccountPgRepo) ListAccounts(ctx context.Context) ([]account.Account, error) {
	var accounts []account.Account
	dbAccounts, err := r.queries.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}
	for _, dbAcc := range dbAccounts {
		accounts = append(accounts, account.Account{
			AccountID:                 dbAcc.AccountID,
			WithingsUserID:            dbAcc.WithingsUserID,
			WithingsAccessToken:       dbAcc.WithingsAccessToken,
			WithingsRefreshToken:      dbAcc.WithingsRefreshToken,
			WithingsAccessTokenExpiry: dbAcc.WithingsAccessTokenExpiry,
			WithingsScopes:            dbAcc.WithingsScopes,
		})
	}
	return accounts, nil
}
