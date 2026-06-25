-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    email TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;