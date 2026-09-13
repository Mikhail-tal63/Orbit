package ride

type RideRequestDto struct {
    PickupLatitude        float64 `json:"pickup_latitude"`
    PickupLongitude       float64 `json:"pickup_longitude"`
    DestinationLatitude   float64 `json:"destination_latitude"`
    DestinationLongitude  float64 `json:"destination_longitude"`
}