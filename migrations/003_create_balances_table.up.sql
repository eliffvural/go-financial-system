CREATE TABLE balances (
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    last_updated_at TIMESTAMP NOT NULL DEFAULT NOW()
); 