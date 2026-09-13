package ride

import (
	"context"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/Mikhail-Tal63/Orbit/internal/geocoding"
)

type RideService struct{
	repository *RideRepository
    Geothing *geocoding.GeoapifyService
}

func NewRideService(repository *RideRepository,Geothing *geocoding.GeoapifyService) *RideService{
	return &RideService{
		repository: repository,
		Geothing: Geothing,
	}
}

func (s *RideService) CreateRideReq(ctx context.Context,payload RideRequestDto)(*db.RideRequest,error){
	return nil,nil
}