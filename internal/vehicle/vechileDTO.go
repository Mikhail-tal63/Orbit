package vehicle

import "github.com/jackc/pgx/v5/pgtype"

type VehicleDTO struct {
	ID          string `json:"id"`
	Make        string `json:"make"`
	Model       string `json:"model"`
	Year        int32  `json:"year"`
	Color       string `json:"color"`
	PlateNumber string `json:"plate_number"`
	CreatedAt   string `json:"created_at"`
	CpdatedAt   string `json:"updated_at"`
}

type CreatVechileRequest struct {
	Make        string `json:"make"`
	Model       string `json:"model"`
	Year        pgtype.Int4  `json:"year"`
	Color       pgtype.Text `json:"color"`
	PlateNumber string `json:"plate_number"`

}
