package location

import (
	"context"
)

// LocationManager is an interface that defines the methods that a location manager should implementge
// For now there are 3 methods: SaveDriverLocation, RemoveDriverLocation and GetNearbyDrivers
// SaveDriverLocation saves the location of a driver given its ID and coordinates
// RemoveDriverLocation removes the location of a driver given its ID
// GetNearbyDrivers returns the IDs of the drivers that are near a given location and a radius in kilometers
type LocationManager interface {
	SaveDriverLocation(ctx context.Context, driverID string, latitude, longitude float64) error
	RemoveDriverLocation(ctx context.Context, driverID string) error
	GetNearbyDrivers(ctx context.Context, latitude, longitude, radius float64) ([]string, error)
}
