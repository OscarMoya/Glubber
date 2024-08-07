package location

import (
	"context"
	"testing"
	"time"

	"github.com/OscarMoya/Glubber/pkg/util"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This function sets up a Redis client for testing purposes and flushes the database
// before each test so there is no problem in multiple runs.
// The main idea is that each test is isolated from the others and thus each DB is different
func setupTestRedis(ctx context.Context, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       db, // use a separate DB for testing
	})

	rdb.FlushDB(ctx) // Clear the database before each test
	return rdb
}

// TestShouldSaveDrivers tests the SaveDriverLocation method of the RedisLocationService
func TestShouldSaveDrivers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rdb := setupTestRedis(ctx, 1) // Use DB 1 for testing
	service := &RedisLocationService{redisClient: rdb}

	// Initial conditions, no interface used so we can test the method in isolation
	_, err := rdb.GeoAdd(ctx, driverLocationKey, &redis.GeoLocation{
		Name:      "driver1",
		Longitude: -74.0060,
		Latitude:  40.7128,
	}).Result()
	require.NoError(t, err) // if we can't add the initial driver, the test is useless

	// Test cases, the main idea is to test that We can update the location of a driver
	// and that we can keep adding drivers to the database
	tests := []struct {
		name            string
		driverID        string
		latitude        float64
		longitude       float64
		expectedDrivers int
	}{
		{
			name:            "Update initial driver",
			driverID:        "driver1",
			latitude:        40.7128,
			longitude:       -74.0060,
			expectedDrivers: 1,
		},
		{
			name:            "Add new driver location",
			driverID:        "driver2",
			latitude:        40.7128,
			longitude:       -74.0060,
			expectedDrivers: 2,
		},
		{
			name:            "Add driver 3",
			driverID:        "driver3",
			latitude:        41.7128,
			longitude:       -75.0060,
			expectedDrivers: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.SaveDriverLocation(ctx, tt.driverID, tt.latitude, tt.longitude)
			assert.NoError(t, err)

			result, err := rdb.ZRange(ctx, driverLocationKey, 0, -1).Result()
			require.NoError(t, err)
			assert.Len(t, result, tt.expectedDrivers)
		})
	}
}

// TestShouldRemoveDriver tests the RemoveDriverLocation method of the RedisLocationService
func TestShouldRemoveDriver(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	rdb := setupTestRedis(ctx, 2)
	service := &RedisLocationService{redisClient: rdb}
	// We add a driver to the database so we can remove it later
	_, err := rdb.GeoAdd(ctx, driverLocationKey, &redis.GeoLocation{
		Name:      "driver1",
		Longitude: -74.0060,
		Latitude:  40.7128,
	}).Result()
	require.NoError(t, err)

	// Initial condition, we have one driver in the database
	result, err := rdb.ZRange(ctx, driverLocationKey, 0, -1).Result()
	require.NoError(t, err)
	assert.Len(t, result, 1)
	// Check that we can remove a driver from the database

	// Test cases, the main idea is to test that we can remove a driver from the database
	// if exists and if it does not exist, the database remains the same and no error is thrown
	tests := []struct {
		name            string
		driverID        string
		expectedDrivers int
	}{
		{
			name:            "Remove non existing driver",
			driverID:        "driver2",
			expectedDrivers: 1,
		},
		{
			name:            "Remove existing driver",
			driverID:        "driver1",
			expectedDrivers: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.RemoveDriverLocation(ctx, tt.driverID)
			assert.NoError(t, err)

			result, err := rdb.ZRange(ctx, driverLocationKey, 0, -1).Result()
			require.NoError(t, err)
			assert.Len(t, result, tt.expectedDrivers)
		})
	}
}

// TestShouldGetNearbyDrivers tests the GetNearbyDrivers method of the RedisLocationService
func TestShouldGetNearbyDrivers(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	rdb := setupTestRedis(ctx, 3)
	service := &RedisLocationService{redisClient: rdb}

	baseLatitude := 40.7128
	baseLongitude := -74.0

	// Add drivers to the database each one at a different distance from the base location
	driver1Latitude, driver1Longitude := util.AddKM(baseLatitude, baseLongitude, 1.5, 90)
	driver2Latitude, driver2Longitude := util.AddKM(baseLatitude, baseLongitude, 3.5, 90)
	driver3Latitude, driver3Longitude := util.AddKM(baseLatitude, baseLongitude, 5.5, 90)

	_, err := rdb.GeoAdd(ctx, driverLocationKey, &redis.GeoLocation{
		Name:      "driver1",
		Longitude: driver1Longitude,
		Latitude:  driver1Latitude,
	}).Result()
	require.NoError(t, err)

	_, err = rdb.GeoAdd(ctx, driverLocationKey, &redis.GeoLocation{
		Name:      "driver2",
		Longitude: driver2Longitude,
		Latitude:  driver2Latitude,
	}).Result()
	require.NoError(t, err)

	_, err = rdb.GeoAdd(ctx, driverLocationKey, &redis.GeoLocation{
		Name:      "driver3",
		Longitude: driver3Longitude,
		Latitude:  driver3Latitude,
	}).Result()
	require.NoError(t, err)
	// Added all the drivers to the database

	// Test cases, the main idea is to test that we can get the drivers that are near a given location
	// The radius will be different in each test case and increasing so we go from 0 to 3 drivers found
	tests := []struct {
		name               string
		passengerLatitude  float64
		passengerLongitude float64
		radius             float64
		expectedDrivers    int
	}{
		{
			name:            "Should not match any driver",
			radius:          1,
			expectedDrivers: 0,
		},
		{
			name:            "Should match driver 1",
			radius:          2,
			expectedDrivers: 1,
		},
		{
			name:            "should match driver 1 and 2",
			radius:          4,
			expectedDrivers: 2,
		},
		{
			name:            "should match all drivers",
			radius:          10,
			expectedDrivers: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drivers, err := service.GetNearbyDrivers(ctx, baseLatitude, baseLongitude, tt.radius)
			require.NoError(t, err)
			assert.Len(t, drivers, tt.expectedDrivers)
		})
	}
}
