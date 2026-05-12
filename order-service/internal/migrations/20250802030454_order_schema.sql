-- +goose Up
-- +goose StatementBegin
CREATE TYPE statuses AS ENUM (
    'pending', 'cancelled', 'confirmed', 'shipped'
);

CREATE TABLE IF NOT EXISTS orders(
    id SERIAL PRIMARY KEY,
    customer_id UUID NOT NULL,
    delivery_address VARCHAR(255) NOT NULL,
    status statuses NOT NULL,
    price DOUBLE PRECISION NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS order_goods (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT NOT NULL,
    seller_id UUID NOT NULL,
    quantity INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_goods;
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS statuses;
-- +goose StatementEnd
