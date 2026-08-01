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

