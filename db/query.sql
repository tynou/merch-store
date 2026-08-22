-- name: GetUserByUsername :one
SELECT id, username, password_hash, balance
FROM users
WHERE username = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password_hash)
VALUES ($1, $2)
RETURNING id;

-- name: GetUserForUpdate :one
SELECT id, username, balance
FROM users
WHERE id = $1
LIMIT 1
FOR UPDATE;

-- name: UpdateUserBalance :exec
UPDATE users
SET balance = $1
WHERE id = $2;

-- name: GetMerchByName :one
SELECT id, name, price
FROM merch
WHERE name = $1
LIMIT 1;

-- name: CreatePurchase :exec
INSERT INTO purchases (user_id, merch_id)
VALUES ($1, $2);