package ride

import (
	"context"
	"strconv"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/Mikhail-Tal63/Orbit/internal/geocoding"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func FuckingNumericFormat(f float64)(pgtype.Numeric,error){
	var n pgtype.Numeric
 
	err := n.Scan(strconv.FormatFloat(f,'f',6,64))
	if err != nil {
		return pgtype.Numeric{},err
	}
    return n,nil
}

func (s *RideService) CreateRideReq(ctx context.Context,payload RideRequestDto,passengerID uuid.UUID)(*db.RideRequest,error){

}