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