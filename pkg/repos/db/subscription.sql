-- name: GetSubscriptionByID :one
SELECT *
FROM subscription
WHERE subscription_id = $1;

-- name: GetSubscriptionsByAccountID :many
SELECT *
FROM subscription
WHERE account_id = $1;

-- name: ListSubscriptions :many
SELECT *
FROM subscription
ORDER BY account_id;

-- name: CreateSubscription :exec
INSERT INTO subscription (account_id, appli, callbackurl, comment)
VALUES ($1, $2, $3, $4);

-- name: DeleteSubscription :exec
DELETE
FROM subscription
WHERE subscription_id = $1;