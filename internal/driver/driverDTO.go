package driver

import "github.com/google/uuid"

type DriverDTO struct {
	ID               uuid.UUID `json:"id"`
	IsOnline         bool      `json:"is_online"`
	IsAvailable      bool      `json:"is_available"`
	CurrentLatitude  float64   `json:"current_latitude"`
	CurrentLongitude float64   `json:"current_longitude"`
	Rating           float64   `json:"rating"`
	CompletedRides   int32     `json:"completed_rides"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
}

type CreateDriverRequest struct {

}