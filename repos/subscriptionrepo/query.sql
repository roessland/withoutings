-- name: GetSubscription :one
SELECT *
FROM subscription
WHERE account_id = $1
LIMIT 1;

-- name: ListSubscription :many
SELECT *
FROM subscription
ORDER BY account_id;

-- name: CreateSubscription :one
INSERT INTO subscription (account_id, appli, callbackurl, comment)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteSubscription :exec
DELETE
FROM subscription
WHERE subscription_id = $1;