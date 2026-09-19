package ride

import (
	"context"
	"strconv"
	"time"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/Mikhail-Tal63/Orbit/internal/geocoding"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RideService struct {
	repository *RideRepository
	Geothing   *geocoding.GeoapifyService
}

func NewRideService(repository *RideRepository, Geothing *geocoding.GeoapifyService) *RideService {
	return &RideService{
		repository: repository,
		Geothing:   Geothing,
	}
}

func (s *RideService) CreateRideReq(ctx context.Context, payload RideRequestDto, passengerID uuid.UUID) (*db.RideRequest, error) {
	pickupAdress, err := s.Geothing.ReverseGeocode(ctx, payload.PickupLatitude, payload.PickupLongitude)
	if err != nil {
		return nil, err
	}
	destinationAdress, err := s.Geothing.ReverseGeocode(ctx, payload.DestinationLatitude, payload.DestinationLongitude)
	if err != nil {
		return nil, err
	}

	pickupLat, err := FuckingPgNumericFormat(payload.PickupLatitude)
	if err != nil {
		return nil, err
	}

	pickupLon, err := FuckingPgNumericFormat(payload.PickupLongitude)
	if err != nil {
		return nil, err
	}

	destinationLat, err := FuckingPgNumericFormat(payload.DestinationLatitude)
	if err != nil {
		return nil, err
	}

	destinationLon, err := FuckingPgNumericFormat(payload.DestinationLongitude)
	if err != nil {
		return nil, err
	}

	rout, err := s.Geothing.GetRoute(ctx,
		payload.PickupLatitude,
		payload.PickupLongitude,
		payload.DestinationLatitude,
		payload.DestinationLongitude,
	)

	DistanceKmIN_FUKING_PG_NUMERICS, err := FuckingPgNumericFormat(rout.DistanceKm)
	if err != nil {
		return nil, err
	}

	price := 5.0 + rout.DistanceKm*2.0

	priceIN_FUKING_PG_NUMERICS, err := FuckingPgNumericFormat(price)
	if err != nil {
		return nil, err
	}

	exp, err := FuckingPgTimetamp(time.Now().Add(2 * time.Minute))

	ride := db.CreateRideReqParams{
		PassengerID:              passengerID,
		PickupLatitude:           pickupLat,
		PickupLongitude:          pickupLon,
		PickupAddress:            pickupAdress,
		DestinationLatitude:      destinationLat,
		DestinationLongitude:     destinationLon,
		DestinationAddress:       destinationAdress,
		EstimatedDistanceKm:      DistanceKmIN_FUKING_PG_NUMERICS,
		EstimatedDurationMinutes: int32(rout.DurationMin),
		EstimatedPrice:           priceIN_FUKING_PG_NUMERICS,
		ExpiresAt:                exp,
	}

	newRide, err := s.repository.CreateRideReq(ctx, ride)
	return newRide, nil
}

// helpers
func FuckingPgNumericFormat(f float64) (pgtype.Numeric, error) {
	var n pgtype.Numeric

	err := n.Scan(strconv.FormatFloat(f, 'f', 6, 64))
	if err != nil {
		return pgtype.Numeric{}, err
	}
	return n, nil
}
func FuckingPgNumericFormatFromIntger(n int64) (pgtype.Numeric, error) {
	var num pgtype.Numeric

	err := num.Scan(strconv.FormatInt(int64(n), 10))
	if err != nil {
		return pgtype.Numeric{}, err
	}
	return num, err
}
func FuckingPgTimetamp(n time.Time) (pgtype.Timestamp, error) {
	var ts pgtype.Timestamp

	err := ts.Scan(n)
	if err != nil {
		return pgtype.Timestamp{}, err
	}
	return ts, err
}
