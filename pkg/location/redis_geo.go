package location

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const (
	// driverLocationKey is the key used to store the driver locations in Redis
	// for now there will be only one key devoted for drivers, we will see if there
	// is a need to split this into multiple keys in the future
	driverLocationKey = "driver_location"
)

// RedisLocationService is a struct that implements the LocationManager interface
type RedisLocationService struct {
	redisClient *redis.Client
}

// NewRedisLocationService creates a new RedisLocationService
func NewRedisLocationService(redisAddr string) *RedisLocationService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisLocationService{redisClient: rdb}
}

// SaveDriverLocation saves the location of a driver given its ID and coordinates
func (r *RedisLocationService) SaveDriverLocation(
	ctx context.Context,
	driverID string,
	latitude,
	longitude float64,
) error {

	key := driverLocationKey // placeholder in case we need to split the keys in the future
	_, err := r.redisClient.GeoAdd(ctx, key, &redis.GeoLocation{
		Name:      driverID,
		Longitude: longitude,
		Latitude:  latitude,
	}).Result()
	if err != nil {
		return err
	}

	// Here we could add a TTL to the key, so that the driver location expires after a certain time
	// Additionally we could cache in a separate key the driver location, so that we can retrieve it
	// faster when needed, if we need to.

	// The reason for not doing it now is for consistency, we will need to create a pipeline,
	// Add an expiration to the set and then add the new Key to the cache. The happy path is very easy
	// to implement, but the error handling is a bit more complex, so we will leave it for a future improvement.

	return nil
}

// GetNearbyDrivers returns the IDs of the drivers that are near a given location and a radius in kilometers
func (r *RedisLocationService) GetNearbyDrivers(
	ctx context.Context,
	passengerLatitude,
	passengerLongitude,
	radius float64, // radius in kilometers
) ([]string, error) {

	key := driverLocationKey // placeholder in case we need to split the keys in the future
	res, err := r.redisClient.GeoRadius(
		ctx,
		key,
		passengerLongitude,
		passengerLatitude,
		// TODO: Add proper parametrization of this Query
		&redis.GeoRadiusQuery{
			Radius:      radius,
			Unit:        "km",
			WithCoord:   false,
			WithDist:    false,
			WithGeoHash: false,
			Count:       0,
			Sort:        "ASC",
		}).Result()
	if err != nil {
		return nil, err
	}

	// This method only returns the driver IDs, but we could return the driver location as well, will set this
	// as a future improvement if there is a need for enrich the response with the driver location.
	// For the basic cases what is going to happen is that this result is going to be enqueued and locked.
	// In theory, the radius set in the query should not impact the end charge of the passenger as long as the drivers
	// are near enough between them.
	var drivers []string
	for _, loc := range res {
		drivers = append(drivers, loc.Name)
	}

	return drivers, nil
}

// RemoveDriverLocation removes the location of a driver given its ID
// This method is useful when a driver logs out of the system so it's location is not accounted anymore
func (r *RedisLocationService) RemoveDriverLocation(ctx context.Context, driverID string) error {
	key := driverLocationKey
	_, err := r.redisClient.ZRem(ctx, key, driverID).Result()
	if err != nil {
		return err
	}

	// If there are additional cache entries for the driver location, we should remove them here as well

	return nil
}
