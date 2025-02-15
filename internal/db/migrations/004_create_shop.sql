CREATE TABLE IF NOT EXISTS shop (
    item TEXT PRIMARY KEY,
    price INT NOT NULL CHECK (price > 0)
);
