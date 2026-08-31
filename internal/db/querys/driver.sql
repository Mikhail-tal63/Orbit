-- name: CreateDriver :one
INSERT INTO drivers (
        id ,
    user_id 
)
VALUES(
    $1 ,
    $2
)
RETURNING *;

-- name: GetDriverByUserId :one
SELECT * FROM drivers WHERE user_id = $1;

-- name: GoOnline :exec
UPDATE drivers
SET 
is_online = TRUE,
updated_at = NOW()
WHERE user_id = $1;

-- name: GoOffline :exec
UPDATE drivers
SET
    is_online = FALSE,
    is_available = FALSE,
    updated_at = NOW()
WHERE user_id = $1;