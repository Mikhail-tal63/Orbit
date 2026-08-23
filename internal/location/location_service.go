package location

import (
	"context"
	"fmt"
)

type LocationService struct {
	store LocationStore
}

func NewLocatingService(store LocationStore) *LocationService {
	return &LocationService{
		store: store,
	}
}

func (s *LocationService) UpdateDriverLocation(ctx context.Context, driverID string, latitude float64, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("invalid latitude")
	}

	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("invalid longitude")
	}

	return s.store.UpdateDriverLocation(ctx, driverID, latitude, longitude)

}
