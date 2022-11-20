-- name: GetSession :one
select session_id, data from session
where session_id = $1;

-- name: CreateSession :one
INSERT INTO session (data)
VALUES ($1)
returning session_id;

-- name: UpsertSession :exec
INSERT INTO session (session_id, data)
VALUES ($1, $2)
ON CONFLICT (session_id) DO UPDATE
    SET data = excluded.data;