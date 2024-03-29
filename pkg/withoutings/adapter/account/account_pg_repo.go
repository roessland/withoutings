package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"time"
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

func fromPgTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		panic("time was not valid")
	}
	if t.InfinityModifier != pgtype.Finite {
		panic("time was not finite")
	}
	return t.Time
}

func toPgTime(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:             t,
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
}

func (r PgRepo) GetAccountByWithingsUserID(ctx context.Context, withingsUserID string) (*account.Account, error) {
	if r.queries == nil {
		panic("queries was nil")
	}
	acc, err := r.queries.GetAccountByWithingsUserID(ctx, withingsUserID)
	if err == pgx.ErrNoRows {
		return nil, account.ErrAccountNotFound
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
		return nil, account.ErrAccountNotFound
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

// Update updates an account in the database.
// updateFn is a function that takes the current account and returns the updated account.
// updateFn is called within a transaction, so it should not start its own transaction.
func (r PgRepo) Update(ctx context.Context, accountUUID uuid.UUID, updateFn func(ctx context.Context, acc *account.Account) (*account.Account, error)) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	inTransaction := r.WithTx(tx)

	acc, err := inTransaction.GetAccountByUUID(ctx, accountUUID)
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
