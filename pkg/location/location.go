package location

import (
	"context"
)

type LocationManager interface {
	SaveDriverLocation(ctx context.Context, driverID string, latitude, longitude float64) error
	RemoveDriverLocation(ctx context.Context, driverID string) error
	GetNearbyDrivers(ctx context.Context, latitude, longitude, radius float64) ([]string, error)
}
