-- name: CreateRideReq :one
INSERT INTO ride_requests (
        id,
    passenger_id,
    driver_id,

    pickup_latitude,
    pickup_longitude,

    destination_latitude,
    destination_longitude,

    pickup_address,
    destination_address,

    estimated_distance_km,
    estimated_duration_minutes,
    estimated_price,

    expires_at
)VALUES(
     $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13
)
RETURNING *;