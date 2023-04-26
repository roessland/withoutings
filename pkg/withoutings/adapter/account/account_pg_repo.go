package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type PgRepo struct {
	db      *pgxpool.Pool
	queries *db.Queries
}

func NewPgRepo(db *pgxpool.Pool, queries *db.Queries) PgRepo {
	return PgRepo{
		db:      db,
		queries: queries,
	}
}

func (r PgRepo) WithTx(tx pgx.Tx) PgRepo {
	return PgRepo{
		db:      r.db,
		queries: r.queries.WithTx(tx),
	}
}

func (r PgRepo) GetAccountByWithingsUserID(ctx context.Context, withingsUserID string) (*account.Account, error) {
	if r.queries == nil {
		panic("queries was nil")
	}
	acc, err := r.queries.GetAccountByWithingsUserID(ctx, withingsUserID)
	if err == pgx.ErrNoRows {
		return nil, account.NotFoundError{WithingsUserID: withingsUserID}
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve account")
	}
	return account.NewAccount(
		acc.AccountUuid,
		acc.WithingsUserID,
		acc.WithingsAccessToken,
		acc.WithingsRefreshToken,
		acc.WithingsAccessTokenExpiry,
		acc.WithingsScopes,
	)
}

func (r PgRepo) GetAccountByUUID(ctx context.Context, accountUUID uuid.UUID) (*account.Account, error) {
	if r.queries == nil {
		panic("queries was nil")
	}
	acc, err := r.queries.GetAccountByAccountUUID(ctx, accountUUID)
	if err == pgx.ErrNoRows {
		return nil, account.NotFoundError{AccountUUID: accountUUID}
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve account")
	}
	return account.NewAccount(
		acc.AccountUuid,
		acc.WithingsUserID,
		acc.WithingsAccessToken,
		acc.WithingsRefreshToken,
		acc.WithingsAccessTokenExpiry,
		acc.WithingsScopes,
	)
}

func (r PgRepo) CreateAccount(ctx context.Context, account *account.Account) error {
	return r.queries.CreateAccount(ctx, db.CreateAccountParams{
		AccountUuid:               account.UUID(),
		WithingsUserID:            account.WithingsUserID(),
		WithingsAccessToken:       account.WithingsAccessToken(),
		WithingsRefreshToken:      account.WithingsRefreshToken(),
		WithingsAccessTokenExpiry: account.WithingsAccessTokenExpiry(),
		WithingsScopes:            account.WithingsScopes(),
	})
}

func (r PgRepo) ListAccounts(ctx context.Context) ([]*account.Account, error) {
	var accounts []*account.Account
	dbAccounts, err := r.queries.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}
	for _, dbAcc := range dbAccounts {
		acc, err := account.NewAccount(
			dbAcc.AccountUuid,
			dbAcc.WithingsUserID,
			dbAcc.WithingsAccessToken,
			dbAcc.WithingsRefreshToken,
			dbAcc.WithingsAccessTokenExpiry,
			dbAcc.WithingsScopes,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

// UpdateAccount updates an account in the database.
// updateFn is a function that takes the current account and returns the updated account.
// updateFn is called within a transaction, so it should not start its own transaction.
func (r PgRepo) UpdateAccount(
	ctx context.Context,
	withingsUserID string,
	updateFn func(ctx context.Context, acc *account.Account) (*account.Account, error),
) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	inTransaction := r.WithTx(tx)

	acc, err := inTransaction.GetAccountByWithingsUserID(ctx, withingsUserID)
	if err != nil {
		return err
	}
	updatedAcc, err := updateFn(ctx, acc)
	err = inTransaction.queries.UpdateAccount(ctx, db.UpdateAccountParams{
		WithingsUserID:            updatedAcc.WithingsUserID(),
		WithingsAccessToken:       updatedAcc.WithingsAccessToken(),
		WithingsRefreshToken:      updatedAcc.WithingsRefreshToken(),
		WithingsAccessTokenExpiry: updatedAcc.WithingsAccessTokenExpiry(),
		WithingsScopes:            updatedAcc.WithingsScopes(),
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
