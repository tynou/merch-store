-- name: GetUserByUsername :one
SELECT id, username, password_hash, balance
FROM users
WHERE username = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password_hash)
VALUES ($1, $2)
RETURNING id;