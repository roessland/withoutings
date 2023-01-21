-- name: GetAccountByID :one
SELECT *
FROM account
WHERE account_id = $1
LIMIT 1;

-- name: GetAccountByWithingsUserID :one
SELECT *
FROM account
WHERE withings_user_id = $1
LIMIT 1;

-- name: ListAccounts :many
SELECT *
FROM account
ORDER BY account_id;

-- name: CreateAccount :exec
INSERT INTO account (withings_user_id, withings_access_token, withings_refresh_token,
                         withings_access_token_expiry, withings_scopes)
VALUES ($1, $2, $3, $4, $5);

-- name: DeleteAccount :exec
DELETE
FROM account
WHERE account_id = $1;