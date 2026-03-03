-- +goose Up
-- +goose StatementBegin
CREATE TABLE sellers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL, 
    location VARCHAR(255) NOT NULL, 
    license INT NOT NULL
);

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    category VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE product_list (
    id BIGSERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    seller_id INT REFERENCES sellers(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE stock (
    seller_id INT REFERENCES sellers(id) ON DELETE CASCADE,
    product_id INT REFERENCES product_list(id) ON DELETE CASCADE,
    stock INT default 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stock;
DROP TABLE product_list;
DROP TABLE categories;
DROP TABLE sellers;
-- +goose StatementEnd
