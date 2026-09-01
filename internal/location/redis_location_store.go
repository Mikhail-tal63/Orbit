package location

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const driverLocationKey = "drivers:locations"

type RedisLocationStore struct {
	client *redis.Client
}

func NewRedisLocationStore(client *redis.Client) *RedisLocationStore {
	return &RedisLocationStore{
		client: client,
	}
}

func (r *RedisLocationStore) UpdateDriverLocation(ctx context.Context, driverID string, latitude float64, longitude float64) error {
	err := r.client.GeoAdd(
		ctx,
		driverLocationKey,
		&redis.GeoLocation{
			Name:      driverID,
			Longitude: longitude,
			Latitude:  latitude,
		},
	).Err()
	if err != nil {
		return fmt.Errorf("failed to update driver location: %w", err)
	}

	return nil
}

func (r *RedisLocationStore) RemoveDriverLocation(ctx context.Context, driverID string) error {
	err := r.client.ZRem(ctx, driverLocationKey, driverID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove driver location: %w", err)
	}
	return nil
}

func (r *RedisLocationStore) FindNearbyDrivers(ctx context.Context, latitude float64, longitude float64, radiuskm float64) ([]string, error) {
	drivers, err := r.client.GeoSearch(ctx, driverLocationKey, &redis.GeoSearchQuery{
		Latitude:   latitude,
		Longitude:  longitude,
		Radius:     radiuskm,
		RadiusUnit: "km",
		Sort:       "ASC",
	},
	).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to find nearby drivers: %w", err)
	}
	return drivers, nil

}
