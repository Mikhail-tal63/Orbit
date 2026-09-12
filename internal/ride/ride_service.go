package ride

type RideService struct{
	repository *RideRepository
}

func NewRideService(repository *RideRepository) *RideService{
	return &RideService{
		repository: repository,
	}
}

