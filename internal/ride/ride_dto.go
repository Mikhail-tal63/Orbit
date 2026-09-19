package ride

import "github.com/google/uuid"

type RideRequestDto struct {
    PickupLatitude        float64 `json:"pickup_latitude"`
    PickupLongitude       float64 `json:"pickup_longitude"`
    DestinationLatitude   float64 `json:"destination_latitude"`
    DestinationLongitude  float64 `json:"destination_longitude"`
}

type RideOffer struct {
	Type string `json:"type"`

	RideRequestID uuid.UUID `json:"ride_request_id"`

	PickupAddress      string `json:"pickup_address"`
	DestinationAddress string `json:"destination_address"`

	EstimatedDistanceKm      float64 `json:"estimated_distance_km"`
	EstimatedDurationMinutes int32   `json:"estimated_duration_minutes"`
	EstimatedPrice           string  `json:"estimated_price"`
}