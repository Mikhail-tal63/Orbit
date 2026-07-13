-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY,

    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,

    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    phone VARCHAR(20) NOT NULL,

    role TEXT NOT NULL,

    image_id UUID,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    last_login_at TIMESTAMP
);

-- +goose Down

DROP TABLE users;