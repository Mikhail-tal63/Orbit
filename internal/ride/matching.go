package ride

import (
	"context"
	"fmt"

	"github.com/Mikhail-Tal63/Orbit/internal/driver"
	"github.com/Mikhail-Tal63/Orbit/internal/location"
	"github.com/google/uuid"
)

type DriverNotifier interface {
	SendToDriver(driverID uuid.UUID, event any) error
}

type Matching struct {
	driverRepo *driver.DriverRepository
	rideRepo   *RideRepository
	location   *location.LocationService
	notifier   DriverNotifier
}

func NewMatchingService(
	driverRepo *driver.DriverRepository,
	rideRepo *RideRepository,
	locationService *location.LocationService,
	notifier DriverNotifier,
) *Matching {
	return &Matching{
		driverRepo: driverRepo,
		rideRepo:   rideRepo,
		location:   locationService,
		notifier:   notifier,
	}
}

func (m *Matching) MatchRide(
	ctx context.Context,
	ride *Ride,
) error {

	driverIDs, err := m.location.FindNearbyDrivers(
		ctx,
		ride.PickupLatitude,
		ride.PickupLongitude,
		5,
	)
	if err != nil {
		return fmt.Errorf("find nearby drivers: %w", err)
	}

	if len(driverIDs) == 0 {
		return fmt.Errorf("no nearby drivers")
	}

	availableDrivers, err := m.driverRepo.GetAvailableDrivers(
		ctx,
		driverIDs,
	)
	if err != nil {
		return fmt.Errorf("get available drivers: %w", err)
	}

	if len(availableDrivers) == 0 {
		return fmt.Errorf("no available drivers")
	}

	offer := RideOffer{
		Type:                     "ride_offer",
		RideRequestID:            ride.ID,
		PickupAddress:            ride.PickupAddress,
		DestinationAddress:       ride.DestinationAddress,
		EstimatedDurationMinutes: ride.EstimatedDurationMinutes,
	}

	maxOffers := 3

	sent := 0

	for _, driverID := range availableDrivers {
		if sent >= maxOffers {
			break
		}

		err := m.notifier.SendToDriver(
			driverID,
			offer,
		)

		
		if err != nil {
			continue
		}

		sent++
	}

	if sent == 0 {
		return fmt.Errorf("could not contact any nearby drivers")
	}

	return nil
}