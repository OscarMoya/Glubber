package model

import "github.com/jackc/pgx/v4"

type RideStatus string

const (
	RideStatusPending            RideStatus = "requested"
	RideStatusPassengerAccepted  RideStatus = "passenger_accepted"
	RideStatusPassengerDenied    RideStatus = "passenger_denied"
	RideStatusDriverAccepted     RideStatus = "matched"
	RideStatusPickingUp          RideStatus = "picking_up"
	RideStatusInTransit          RideStatus = "in_transit"
	RideStatusCompleted          RideStatus = "passenger_dropped"
	RideStatusPassengerCancelled RideStatus = "passenger_cancelled"
	RideStatusDriverCancelled    RideStatus = "driver_cancelled"
	RideStatusErrored            RideStatus = "errored"
)

// Ride represents a ride in the system
// Rides are created by passengers and are the main object that will contain the workflow to match a driver with a passenger
// and to calculate the price of the ride-
type Ride struct {
	ID          int        `json:"id"`
	PassengerID int        `json:"passenger_id"`
	DriverID    *int       `json:"driver_id"`
	Price       float64    `json:"price"`
	Status      RideStatus `json:"status"`
	SrcLat      float64    `json:"src_lat"`
	SrcLon      float64    `json:"src_lon"`
	DstLat      float64    `json:"dst_lat"`
	DstLon      float64    `json:"dst_lon"`
}

// Scan is a method that allows us to convert a row from the database into a Ride struct
func (r *Ride) Scan(row pgx.Row) error {
	return row.Scan(
		&r.ID,
		&r.PassengerID,
		&r.DriverID,
		&r.Price,
		&r.Status,
		&r.SrcLat,
		&r.SrcLon,
		&r.DstLat,
		&r.DstLon,
	)
}
