package driver

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type DriverDTO struct {
	ID               uuid.UUID `json:"id"`
	IsOnline         bool      `json:"is_online"`
	IsAvailable      bool      `json:"is_available"`
	CurrentLatitude  pgtype.Numeric   `json:"current_latitude"`
	CurrentLongitude pgtype.Numeric   `json:"current_longitude"`
	Rating           pgtype.Numeric   `json:"rating"`
	CompletedRides   pgtype.Int4     `json:"completed_rides"`
	CreatedAt        pgtype.Timestamp    `json:"created_at"`
	UpdatedAt        pgtype.Timestamp    `json:"updated_at"`
}

type CreateDriverRequest struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
	Make        string `json:"make"`
	Model       string `json:"model"`
	Year        pgtype.Int4  `json:"year"`
	Color       pgtype.Text `json:"color"`
	PlateNumber string `json:"plate_number"`
}
