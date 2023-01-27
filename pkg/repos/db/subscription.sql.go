// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: subscription.sql

package db

import (
	"context"
)

const createSubscription = `-- name: CreateSubscription :exec
INSERT INTO subscription (account_id, appli, callbackurl, webhook_secret, comment)
VALUES ($1, $2, $3, $4, $5)
`

type CreateSubscriptionParams struct {
	AccountID     int64
	Appli         int32
	Callbackurl   string
	WebhookSecret string
	Comment       string
}

func (q *Queries) CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) error {
	_, err := q.db.Exec(ctx, createSubscription,
		arg.AccountID,
		arg.Appli,
		arg.Callbackurl,
		arg.WebhookSecret,
		arg.Comment,
	)
	return err
}

const deleteSubscription = `-- name: DeleteSubscription :exec
DELETE
FROM subscription
WHERE subscription_id = $1
`

func (q *Queries) DeleteSubscription(ctx context.Context, subscriptionID int64) error {
	_, err := q.db.Exec(ctx, deleteSubscription, subscriptionID)
	return err
}

const getSubscriptionByID = `-- name: GetSubscriptionByID :one
SELECT subscription_id, account_id, appli, callbackurl, webhook_secret, comment
FROM subscription
WHERE subscription_id = $1
`

func (q *Queries) GetSubscriptionByID(ctx context.Context, subscriptionID int64) (Subscription, error) {
	row := q.db.QueryRow(ctx, getSubscriptionByID, subscriptionID)
	var i Subscription
	err := row.Scan(
		&i.SubscriptionID,
		&i.AccountID,
		&i.Appli,
		&i.Callbackurl,
		&i.WebhookSecret,
		&i.Comment,
	)
	return i, err
}

const getSubscriptionsByAccountID = `-- name: GetSubscriptionsByAccountID :many
SELECT subscription_id, account_id, appli, callbackurl, webhook_secret, comment
FROM subscription
WHERE account_id = $1
`

func (q *Queries) GetSubscriptionsByAccountID(ctx context.Context, accountID int64) ([]Subscription, error) {
	rows, err := q.db.Query(ctx, getSubscriptionsByAccountID, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subscription
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.SubscriptionID,
			&i.AccountID,
			&i.Appli,
			&i.Callbackurl,
			&i.WebhookSecret,
			&i.Comment,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSubscriptions = `-- name: ListSubscriptions :many
SELECT subscription_id, account_id, appli, callbackurl, webhook_secret, comment
FROM subscription
ORDER BY account_id
`

func (q *Queries) ListSubscriptions(ctx context.Context) ([]Subscription, error) {
	rows, err := q.db.Query(ctx, listSubscriptions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subscription
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.SubscriptionID,
			&i.AccountID,
			&i.Appli,
			&i.Callbackurl,
			&i.WebhookSecret,
			&i.Comment,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
