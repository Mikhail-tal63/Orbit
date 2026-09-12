-- +goose Up

CREATE TYPE ride_status AS ENUM (
    'requested',
    'searching',
    'accepted',
    'driver_arriving',
    'driver_arrived',
    'in_progress',
    'completed',
    'cancelled',
    'expired'
);

CREATE TABLE ride_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    passenger_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE RESTRICT,

    driver_id UUID
        REFERENCES drivers(id)
        ON DELETE SET NULL,

    pickup_latitude DECIMAL(9,6) NOT NULL,
    pickup_longitude DECIMAL(9,6) NOT NULL,

    destination_latitude DECIMAL(9,6) NOT NULL,
    destination_longitude DECIMAL(9,6) NOT NULL,

    pickup_address TEXT NOT NULL,
    destination_address TEXT NOT NULL,

    estimated_distance_km DECIMAL(10,2) NOT NULL,
    estimated_duration_minutes INT NOT NULL,
    estimated_price DECIMAL(10,2) NOT NULL,

    status ride_status NOT NULL DEFAULT 'requested',

    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ride_requests_passenger_id
    ON ride_requests(passenger_id);

CREATE INDEX idx_ride_requests_driver_id
    ON ride_requests(driver_id);

CREATE INDEX idx_ride_requests_status
    ON ride_requests(status);

-- +goose Down

DROP TABLE ride_requests;
DROP TYPE ride_status;