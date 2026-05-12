-- +goose Up
-- +goose StatementBegin

CREATE TYPE roles AS ENUM (
    'seller', 'customer', 'admin'
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    verified BOOLEAN NOT NULL DEFAULT FALSE, 
    role roles NOT NULL, 
    hashpassword TEXT NOT NULL 
);

CREATE TABLE IF NOT EXISTS refreshSessions (
    id SERIAL PRIMARY KEY,
    userId uuid REFERENCES users(id) ON DELETE CASCADE,
    refreshToken uuid NOT NULL,
    fingerprint character varying(200) NOT NULL,
    ip character varying(15) NOT NULL,
    expiresIn bigint NOT NULL,
    createdAt timestamp with time zone NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refreshSessions;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS roles;
-- +goose StatementEnd
