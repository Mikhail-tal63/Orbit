package driver

import (
	"context"

	"github.com/Mikhail-Tal63/Orbit/internal/auth"
	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/Mikhail-Tal63/Orbit/internal/vehicle"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DriverSevrice struct {
	pool       *pgxpool.Pool
	repository DriverRepository
	user       auth.AuthRepository
	vechail    vehicle.VehicleRepository
}

func NewDriverService(repository DriverRepository, user auth.AuthRepository, vechail vehicle.VehicleRepository, pool *pgxpool.Pool) *DriverSevrice {
	return &DriverSevrice{
		repository: repository,
		user:       user,
		vechail:    vechail,
		pool:       pool,
	}
}

func (s *DriverSevrice) CreateDriver(ctx context.Context, userID uuid.UUID, driverPayload CreateDriverRequest) (*DriverDTO, error) {

	existedUser, err := s.user.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existedUser == nil {
		return nil, auth.ErrUserNotFound
	}

	if existedUser.Role == "Driver" {
		return nil, ErrDriverAlreadyExists
	}
	newDriver := db.CreateDriverParams{
		ID:     uuid.New(),
		UserID: userID,
	}
	

	changedRole := db.UpdateUserRoleParams{
		ID:   userID,
		Role: "Driver",
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := db.New(tx)

	driverRepo := NewDriverRepository(q)
	vehicleRepo := vehicle.NewVechileRepository(q)
	userRepo := auth.NewAuthRepository(q)



	createdDriver, err := driverRepo.CreateDriver(ctx, newDriver)
	if err != nil {
		return nil, err
	}
	NewVechile := db.CreateVehicleParams{
		ID:          uuid.New(),
		DriverID:    createdDriver.ID,
		Make:        driverPayload.Make,
		Model:       driverPayload.Model,
		Year:        driverPayload.Year,
		Color:       driverPayload.Color,
		PlateNumber: driverPayload.PlateNumber,
	}
	_, err = vehicleRepo.CreateVechile(ctx, &NewVechile)
	if err != nil {
		return nil, err
	}
	_, err = userRepo.UpdateUserRole(ctx, changedRole)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &DriverDTO{
		ID:               newDriver.ID,
		IsOnline:         createdDriver.IsOnline,
		IsAvailable:      createdDriver.IsAvailable,
		CurrentLatitude:  createdDriver.CurrentLatitude,
		CurrentLongitude: createdDriver.CurrentLongitude,
		Rating:           createdDriver.Rating,
		CompletedRides:   createdDriver.CompletedRides,
		CreatedAt:        createdDriver.CreatedAt,
		UpdatedAt:        createdDriver.UpdatedAt,
	}, nil
}
