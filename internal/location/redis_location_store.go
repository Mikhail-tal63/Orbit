package location

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	driverLocationKey = "drivers:locations"

	driverLastSeenPrefix = "drivers:last_seen:"

	driverLocationTTL = 10 * time.Second
)

type RedisLocationStore struct {
	client *redis.Client
}

func NewRedisLocationStore(client *redis.Client) *RedisLocationStore {
	return &RedisLocationStore{
		client: client,
	}
}

func driverLastSeenKey(driverID string) string {
	return driverLastSeenPrefix + driverID
}

func (r *RedisLocationStore) UpdateDriverLocation(
	ctx context.Context,
	driverID string,
	latitude float64,
	longitude float64,
) error {

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

	now := time.Now().Unix()

	err = r.client.Set(
		ctx,
		driverLastSeenKey(driverID),
		strconv.FormatInt(now, 10),
		0,
	).Err()

	if err != nil {
		return fmt.Errorf("failed to update driver last seen: %w", err)
	}

	return nil
}

func (r *RedisLocationStore) RemoveDriverLocation(
	ctx context.Context,
	driverID string,
) error {

	err := r.client.ZRem(
		ctx,
		driverLocationKey,
		driverID,
	).Err()

	if err != nil {
		return fmt.Errorf("failed to remove driver location: %w", err)
	}

	err = r.client.Del(
		ctx,
		driverLastSeenKey(driverID),
	).Err()

	if err != nil {
		return fmt.Errorf("failed to remove driver last seen: %w", err)
	}

	return nil
}

func (r *RedisLocationStore) FindNearbyDrivers(
	ctx context.Context,
	latitude float64,
	longitude float64,
	radiusKm float64,
) ([]string, error) {

	drivers, err := r.client.GeoSearch(
		ctx,
		driverLocationKey,
		&redis.GeoSearchQuery{
			Latitude:   latitude,
			Longitude:  longitude,
			Radius:     radiusKm,
			RadiusUnit: "km",
			Sort:       "ASC",
		},
	).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to find nearby drivers: %w", err)
	}

	if len(drivers) == 0 {
		return []string{}, nil
	}

	activeDrivers := make([]string, 0, len(drivers))

	now := time.Now().Unix()
	staleBefore := now - int64(driverLocationTTL.Seconds())

	for _, driverID := range drivers {

		lastSeenValue, err := r.client.Get(
			ctx,
			driverLastSeenKey(driverID),
		).Result()

		if err != nil {
			if err == redis.Nil {
				continue
			}

			return nil, fmt.Errorf(
				"failed to get last seen for driver %s: %w",
				driverID,
				err,
			)
		}

		lastSeen, err := strconv.ParseInt(lastSeenValue, 10, 64)
		if err != nil {
			continue
		}

		if lastSeen < staleBefore {
			_ = r.client.ZRem(
				ctx,
				driverLocationKey,
				driverID,
			).Err()

			_ = r.client.Del(
				ctx,
				driverLastSeenKey(driverID),
			).Err()

			continue
		}

		activeDrivers = append(activeDrivers, driverID)
	}

	return activeDrivers, nil
}