CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(255),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR(255) REFERENCES orders(order_uid),
    name VARCHAR(255),
    phone VARCHAR(255),
    zip VARCHAR(255),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS payment (
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(255),
    provider VARCHAR(255),
    amount NUMERIC(10, 2),
    payment_dt BIGINT,
    bank VARCHAR(255),
    delivery_cost NUMERIC(10, 2),
    goods_total NUMERIC(10, 2),
    custom_fee NUMERIC(10, 2),
    order_uid VARCHAR(255) REFERENCES orders(order_uid)
);

CREATE TABLE IF NOT EXISTS order_items (
    chrt_id INT PRIMARY KEY,
    track_number VARCHAR(255),
    price NUMERIC(10, 2),
    rid VARCHAR(255),
    name VARCHAR(255),
    sale NUMERIC(5, 2),
    size VARCHAR(255),
    total_price NUMERIC(10, 2),
    nm_id INT,
    brand VARCHAR(255),
    status INT,
    order_uid VARCHAR(255) REFERENCES orders(order_uid)
);



