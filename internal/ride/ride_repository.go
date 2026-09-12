package ride

import (
	"context"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
)

type RideRepository struct {
	reqQuerys db.Querier
}

func NewRideRepository(reqQuerys db.Querier) *RideRepository {
	return &RideRepository{
		reqQuerys: reqQuerys,
	}
}

func (r *RideRepository) CreateRideReq(ctx context.Context, params db.CreateRideReqParams) (*db.RideRequest, error) {
	ride, err := r.reqQuerys.CreateRideReq(ctx, params)
	if err != nil {
		return nil, err
	}
	return &ride, nil
}
