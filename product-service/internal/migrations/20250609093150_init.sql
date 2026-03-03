-- +goose Up
-- +goose StatementBegin
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    category VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE product_list (
    id BIGSERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL UNIQUE,
    price DOUBLE PRECISION NOT NULL,
    seller_id INT NOT NULL,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE stock (
    seller_id INT NOT NULL,
    product_id INT REFERENCES product_list(id) ON DELETE CASCADE,
    stock INT default 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stock;
DROP TABLE product_list;
DROP TABLE categories;
-- +goose StatementEnd
