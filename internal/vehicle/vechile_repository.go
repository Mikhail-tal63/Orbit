package vehicle

import (
	"context"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
)

type VehicleRepository struct {
	db *db.Queries
}

func NewVechileRepository(db *db.Queries) *VehicleRepository {
	return &VehicleRepository{
		db: db,
	}
}

func (r *VehicleRepository) CreateVechile(ctx context.Context, params *db.CreateVehicleParams) (*db.Vehicle, error) {
	vechile, err := r.db.CreateVehicle(ctx, *params)
	if err != nil {
		return nil, err
	}
	return &vechile, nil
}
