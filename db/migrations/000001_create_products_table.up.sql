CREATE TABLE IF NOT EXISTS product
(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    status INTEGER NOT NULL DEFAULT 0
);