package geocoding

import "context"

type Service interface {
    ReverseGeocode(ctx context.Context, lat, lon float64) (string, error)
}