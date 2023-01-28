-- name: GetSubscriptionByID :one
SELECT *
FROM subscription
WHERE subscription_id = $1;

-- name: GetSubscriptionByWebhookSecret :one
SELECT *
FROM subscription
WHERE webhook_secret = $1;

-- name: GetSubscriptionsByAccountID :many
SELECT *
FROM subscription
WHERE account_id = $1;

-- name: GetPendingSubscriptions :many
SELECT *
FROM subscription
WHERE status = 'pending';

-- name: ListSubscriptions :many
SELECT *
FROM subscription
ORDER BY account_id;

-- name: CreateSubscription :exec
INSERT INTO subscription (account_id, appli, callbackurl, webhook_secret, status, comment)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: DeleteSubscription :exec
DELETE
FROM subscription
WHERE subscription_id = $1;

-- name: CreateRawNotification :exec
INSERT INTO raw_notification (source, status, data)
VALUES ($1, $2, $3);

-- name: GetPendingRawNotifications :many
SELECT *
FROM raw_notification
WHERE status == 'pending';

