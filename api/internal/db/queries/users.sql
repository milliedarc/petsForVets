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

-- name: GetUserPasswordHash :one
SELECT password_hash
FROM users
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = @password_hash
WHERE id = @id;

-- name: UpdateUserAppMode :one
UPDATE users
SET app_mode = @app_mode
WHERE id = @id
RETURNING id, email, name, surname, app_mode, veterinary;
