package model

import "github.com/jackc/pgx/v4"

// Ride represents a ride in the system
// Rides are created by passengers and are the main object that will contain the workflow to match a driver with a passenger
// and to calculate the price of the ride-
type Ride struct {
	ID          int     `json:"id"`
	PassengerID int     `json:"passenger_id"`
	DriverID    *int    `json:"driver_id"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
}

// Scan is a method that allows us to convert a row from the database into a Ride struct
func (r *Ride) Scan(row pgx.Row) error {
	return row.Scan(
		&r.ID,
		&r.PassengerID,
		&r.DriverID,
		&r.Price,
		&r.Status,
	)
}
