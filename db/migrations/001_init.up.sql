CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    balance INT NOT NULL DEFAULT 1000,
    CONSTRAINT check_balance_non_negative CHECK (balance >= 0)
);

CREATE TABLE IF NOT EXISTS merch (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    price INT NOT NULL,
    CONSTRAINT check_price_positive CHECK (price > 0)
);

CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    merch_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (merch_id) REFERENCES merch(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS coin_transfers (
    id SERIAL PRIMARY KEY,
    from_id INT NOT NULL,
    to_id INT NOT NULL,
    amount INT NOT NULL,
    FOREIGN KEY (from_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (to_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_amount_positive CHECK (amount > 0),
    CONSTRAINT check_self_transfer CHECK (from_id <> to_id)
);

INSERT INTO merch (name, price) VALUES
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 500)
ON CONFLICT (name) DO NOTHING;