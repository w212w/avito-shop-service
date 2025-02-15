CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    item TEXT NOT NULL,
    price INT NOT NULL CHECK (price > 0),
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);
