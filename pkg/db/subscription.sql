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
WHERE account_uuid = $1
  AND appli = $2;

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
INSERT INTO subscription(subscription_uuid,
                         account_uuid,
                         appli,
                         callbackurl,
                         webhook_secret,
                         status,
                         comment,
                         status_last_checked_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateSubscription :exec
UPDATE subscription
SET account_uuid           = $1,
    appli                  = $2,
    callbackurl            = $3,
    webhook_secret         = $4,
    status                 = $5,
    comment                = $6,
    status_last_checked_at = $7
WHERE subscription_uuid = $8;

-- name: DeleteSubscription :exec
DELETE
FROM subscription
WHERE subscription_uuid = $1;

-- name: CreateRawNotification :exec
INSERT INTO raw_notification (raw_notification_uuid, source, status, data)
VALUES ($1, $2, $3, $4);

-- name: DeleteRawNotification :exec
DELETE
FROM raw_notification
WHERE raw_notification_uuid = $1;

-- name: GetPendingRawNotifications :many
SELECT *
FROM raw_notification
WHERE status = 'pending'
ORDER BY raw_notification_id;

-- name: GetRawNotification :one
SELECT *
FROM raw_notification
WHERE raw_notification_uuid = $1;


-- name: CreateNotification :exec
INSERT INTO notification(notification_uuid,
                         account_uuid,
                         received_at,
                         params,
                         data,
                         data_status,
                         fetched_at,
                         raw_notification_uuid,
                         source)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (notification_uuid) DO NOTHING;


-- name: UpdateRawNotification :exec
UPDATE raw_notification
SET source       = $1,
    status       = $2,
    data         = $3,
    received_at  = $4,
    processed_at = $5
WHERE raw_notification_uuid = $6;

-- name: GetNotificationsByAccountUUID :many
SELECT *
FROM notification
WHERE account_uuid = $1
ORDER BY received_at DESC;

-- name: GetNotificationByUUID :one
SELECT *
FROM notification
WHERE notification_uuid = $1;


-- name: UpdateNotification :exec
UPDATE notification
SET account_uuid          = $1,
    received_at           = $2,
    params                = $3,
    data                  = $4,
    data_status           = $5,
    fetched_at            = $6,
    raw_notification_uuid = $7,
    source                = $8
    WHERE notification_uuid = $9;