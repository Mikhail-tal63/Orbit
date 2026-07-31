package driver

import (
	"context"

	"github.com/google/uuid"
)

type DriverSevrice struct{
	repository *DriverRepository
}

func NewDriverService(repository *DriverRepository) *DriverSevrice{
	return &DriverSevrice{
		repository: repository,
	}
}

func (s *DriverSevrice) CreateDriver(ctx context.Context,userID uuid.UUID)(*DriverDTO,error){
	driverID := uuid.New()

print(driverID)
return nil ,nil

}