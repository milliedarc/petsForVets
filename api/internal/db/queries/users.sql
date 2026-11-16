-- Each query has a "-- name: FunctionName :kind" comment. sqlc turns it into
-- a Go function. :one returns one row, :many a slice, :exec nothing.
-- $1, $2 ... are parameters, sent separately from the SQL (like PDO prepared
-- statements), so SQL injection isn't possible.

-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, surname)
VALUES ($1, $2, $3, $4)
RETURNING id, email, name, surname, app_mode, veterinary;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, name, surname, app_mode, veterinary
FROM users
WHERE email = $1;
