-- name: CreateUser :exec
INSERT INTO users (
    id,
    first_name,
    last_name,
    username,
    email,
    password_hash,
    phone,
    role,
    image_id,
    is_active,
    created_at,
    updated_at,
    last_login_at
)
VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, $12, $13
);

-- name: GetUserByEmail :one
SELECT
    id,
    first_name,
    last_name,
    username,
    email,
    password_hash,
    phone,
    role,
    image_id,
    is_active,
    created_at,
    updated_at,
    last_login_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateLastLogin :exec
UPDATE users
SET last_login_at = $2
WHERE id = $1;


-- name: GetUserByID :one
SELECT
    id,
    first_name,
    last_name,
    username,
    email,
    password_hash,
    phone,
    role,
    image_id,
    is_active,
    created_at,
    updated_at,
    last_login_at
FROM users
WHERE id = $1
LIMIT 1;