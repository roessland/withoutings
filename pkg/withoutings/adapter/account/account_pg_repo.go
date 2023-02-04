package account

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	db2 "github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type AccountPgRepo struct {
	db      *pgxpool.Pool
	queries *db2.Queries
}

func (r AccountPgRepo) WithTx(tx pgx.Tx) AccountPgRepo {
	return AccountPgRepo{
		db:      r.db,
		queries: r.queries.WithTx(tx),
	}
}

func NewAccountPgRepo(db *pgxpool.Pool, queries *db2.Queries) AccountPgRepo {
	return AccountPgRepo{
		db:      db,
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
	return r.queries.CreateAccount(ctx, db2.CreateAccountParams{
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

func (r AccountPgRepo) UpdateAccount(ctx context.Context, accountID int64, updateFn func(ctx context.Context, acc account.Account) (account.Account, error)) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	acc, err := r.WithTx(tx).GetAccountByID(ctx, accountID)
	if err != nil {
		return err
	}
	updatedAcc, err := updateFn(ctx, acc)
	err = r.WithTx(tx).queries.UpdateAccount(ctx, db2.UpdateAccountParams{
		AccountID:                 accountID,
		WithingsUserID:            updatedAcc.WithingsUserID,
		WithingsAccessToken:       updatedAcc.WithingsAccessToken,
		WithingsRefreshToken:      updatedAcc.WithingsRefreshToken,
		WithingsAccessTokenExpiry: updatedAcc.WithingsAccessTokenExpiry,
		WithingsScopes:            updatedAcc.WithingsScopes,
	})
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
