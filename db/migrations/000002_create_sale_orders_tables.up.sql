CREATE TABLE IF NOT EXISTS sale_order
(
    id INTEGER PRIMARY KEY,
    number TEXT UNIQUE NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    status INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS sale_order_product
(
    id INTEGER PRIMARY KEY,
    parent_id INTEGER,
    product_id INTEGER NOT NULL,
    quantity DECIMAL NOT NULL,
    price DECIMAL NOT NULL,
    FOREIGN KEY(parent_id) REFERENCES sale_order(id),
    FOREIGN KEY(product_id) REFERENCES product(id)
);