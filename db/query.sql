-- name: GetUserByUsername :one
SELECT id, username, password_hash, balance
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserBalance :one
SELECT balance
FROM users
WHERE id = $1
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

-- name: CreateTransfer :exec
INSERT INTO coin_transfers (from_id, to_id, amount)
VALUES ($1, $2, $3);

-- name: GetUserInventory :many
SELECT m.name AS type, COUNT(*)::int AS quantity
FROM purchases AS p
JOIN merch AS m ON p.merch_id = m.id
WHERE p.user_id = $1
GROUP BY m.name
ORDER BY m.name;

-- name: GetReceivedTransfers :many
SELECT u.username AS from_user, SUM(t.amount)::int AS amount
FROM coin_transfers AS t
JOIN users AS u ON t.from_id = u.id
WHERE t.to_id = $1
GROUP BY u.username
ORDER BY u.username;

-- name: GetSentTransfers :many
SELECT u.username AS to_user, SUM(t.amount)::int AS amount
FROM coin_transfers AS t
JOIN users AS u ON t.to_id = u.id
WHERE t.from_id = $1
GROUP BY u.username
ORDER BY u.username;

-- name: CleanTestData :exec
TRUNCATE TABLE users, purchases, coin_transfers RESTART IDENTITY CASCADE;