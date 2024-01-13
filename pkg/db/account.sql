-- name: GetAccountByWithingsUserID :one
SELECT *
FROM account
WHERE withings_user_id = $1
LIMIT 1;

-- name: GetAccountByAccountUUID :one
SELECT
    account_id,
    withings_user_id,
    withings_access_token,
    withings_refresh_token,
    withings_access_token_expiry,
    withings_scopes,
    account_uuid
FROM account
WHERE account_uuid = $1
LIMIT 1;

-- name: ListAccounts :many
SELECT *
FROM account
ORDER BY account_id;

-- name: CreateAccount :exec
INSERT INTO account (account_uuid, withings_user_id, withings_access_token, withings_refresh_token,
                     withings_access_token_expiry, withings_scopes)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateAccount :exec
UPDATE account
SET withings_access_token=$1,
    withings_refresh_token=$2,
    withings_access_token_expiry=$3,
    withings_scopes=$4
WHERE withings_user_id = $5;

-- name: DeleteAccount :exec
DELETE
FROM account
WHERE withings_user_id = $1;