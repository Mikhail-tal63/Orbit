package location

import "context"

type LocationStore interface {
	UpdateDriverLocation(
		ctx context.Context,
		driverID string,
		latitude float64,
		longitude float64,
	) error

	RemoveDriverLocation(
		ctx context.Context,
		driverID string,
	) error

	FindNearbyDrivers(
		ctx context.Context,
		latitude float64,
		longitude float64,
		radiusKm float64,
	) ([]string, error)
}
