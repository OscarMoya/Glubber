package model

import (
	"database/sql"

	"github.com/jackc/pgx/v4"
)

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
	RideStatusDeleted            RideStatus = "deleted"
)

// Ride represents a ride in the system
// Rides are created by passengers and are the main object that will contain the workflow to match a driver with a passenger
// and to calculate the price of the ride-
type Ride struct {
	ID          int        `json:"id" db:"id"`
	PassengerID int        `json:"passenger_id" db:"passenger_id"`
	DriverID    *int       `json:"driver_id" db:"driver_id"`
	Price       float64    `json:"price" db:"price"`
	Status      RideStatus `json:"status" db:"status"`
	SrcLat      float64    `json:"src_lat" db:"src_lat"`
	SrcLon      float64    `json:"src_lon" db:"src_lon"`
	DstLat      float64    `json:"dst_lat" db:"dst_lat"`
	DstLon      float64    `json:"dst_lon" db:"dst_lon"`
}

// Scan is a method that allows us to convert a row from the database into a Ride struct
func (r *Ride) Scan(row pgx.Row) error {
	var driverID sql.NullInt64
	err := row.Scan(
		&r.ID,
		&r.PassengerID,
		&driverID,
		&r.Price,
		&r.Status,
		&r.SrcLat,
		&r.SrcLon,
		&r.DstLat,
		&r.DstLon,
	)
	if err != nil {
		return err
	}

	if driverID.Valid {
		driver := int(driverID.Int64)
		r.DriverID = &driver
	} else {
		r.DriverID = nil
	}

	return nil
}

type RideOutbox struct {
	ID     int        `json:"id" db:"id"`
	RideID int        `json:"ride_id" db:"ride_id"`
	Status RideStatus `json:"status" db:"status"`
}

// Scan is a method that allows us to convert a row from the database into a RideOutbox struct
func (r *RideOutbox) Scan(row pgx.Row) error {
	return row.Scan(&r.ID, &r.RideID, &r.Status)
}

// NewRideOutbox creates a new RideOutbox struct
func NewRideOutbox(ride *Ride) *RideOutbox {
	return &RideOutbox{
		RideID: ride.ID,
		Status: ride.Status,
	}
}
