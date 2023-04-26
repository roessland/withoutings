-- name: AllNotificationCategories :many
SELECT *
FROM notification_category
ORDER BY appli;

-- name: GetSubscriptionByUUID :one
SELECT *
FROM subscription
WHERE subscription_uuid = $1;

-- name: GetSubscriptionByWebhookSecret :one
SELECT *
FROM subscription
WHERE webhook_secret = $1;

-- name: GetSubscriptionsByAccountUUID :many
SELECT *
FROM subscription
WHERE account_uuid = $1
ORDER BY appli;

-- name: GetSubscriptionByAccountUUIDAndAppli :one
SELECT *
FROM subscription
WHERE account_uuid = $1 AND appli = $2;

-- name: GetPendingSubscriptions :many
SELECT *
FROM subscription
WHERE status = 'pending'
ORDER BY subscription.account_uuid;

-- name: ListSubscriptions :many
SELECT *
FROM subscription
ORDER BY account_uuid;

-- name: CreateSubscription :exec
INSERT INTO subscription (subscription_uuid, account_uuid, appli, callbackurl, webhook_secret, status, comment)
VALUES ($1, $2, $3, $4, $5, $6, $7);

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
WHERE status == 'pending'
ORDER BY raw_notification_id;

