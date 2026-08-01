package vehicle

import (
	"context"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
)

type VehicleRepository interface{
	CreateVechile(ctx context.Context, params *db.CreateVehicleParams) (*db.Vehicle, error)
}

type VehicleRepositoryImpl struct {
	db *db.Queries
}

func NewVechileRepository(db *db.Queries) *VehicleRepositoryImpl {
	return &VehicleRepositoryImpl{
		db: db,
	}
}

func (r *VehicleRepositoryImpl) CreateVechile(ctx context.Context, params *db.CreateVehicleParams) (*db.Vehicle, error) {
	vechile, err := r.db.CreateVehicle(ctx, *params)
	if err != nil {
		return nil, err
	}
	return &vechile, nil
}
