-- name: CreateVehicle :one

INSERT INTO vehicles (
    id,
    driver_id,
    make,
    model,
    year,
    color,
    plate_number,
    image_file_id
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;