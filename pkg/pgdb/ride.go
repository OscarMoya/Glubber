package pgdb

import (
	"context"

	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RideCruder interface {
	CreateRide(ctx context.Context, ride *model.Ride) error
	ListRides(ctx context.Context) ([]model.Ride, error)
	GetRide(ctx context.Context, id int) (*model.Ride, error)
	UpdateRide(ctx context.Context, ride *model.Ride) error
	DeleteRide(ctx context.Context, id int) error
}

type RideDatabase struct {
	Pool  *pgxpool.Pool
	Table string
}

// NewRideDatabase creates a new RideDatabase
func NewRideDatabase(ctx context.Context, connectionString string, table string) (*RideDatabase, error) {
	pool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	db := &RideDatabase{
		Pool:  pool,
		Table: table,
	}

	if err := db.createTable(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *RideDatabase) createTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS rides (
		id SERIAL PRIMARY KEY,
		passenger_id INTEGER,
		driver_id INTEGER,
		price FLOAT,
		status VARCHAR(50)
	);`
	_, err := db.Pool.Exec(ctx, query)
	return err
}

func (db *RideDatabase) Close() {
	db.Pool.Close()
}

func (db *RideDatabase) CreateRide(ctx context.Context, ride *model.Ride) error {
	query := `INSERT INTO rides (passenger_id, driver_id, price, status) VALUES ($1, $2, $3, $4) RETURNING id;`
	err := db.Pool.QueryRow(ctx, query, ride.PassengerID, ride.DriverID, ride.Price, ride.Status).Scan(&ride.ID)
	return err
}

func (db *RideDatabase) ListRides(ctx context.Context) ([]model.Ride, error) {
	query := `SELECT id, passenger_id, driver_id, price, status FROM rides;`
	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rides := []model.Ride{}
	for rows.Next() {
		ride := model.Ride{}
		err = ride.Scan(rows)
		if err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}

	return rides, nil
}

func (db *RideDatabase) GetRide(ctx context.Context, id int) (*model.Ride, error) {
	query := `SELECT id, passenger_id, driver_id, price, status FROM rides WHERE id = $1;`
	row := db.Pool.QueryRow(ctx, query, id)
	ride := &model.Ride{}
	err := ride.Scan(row)
	if err != nil {
		return nil, err
	}

	return ride, nil
}

func (db *RideDatabase) UpdateRide(ctx context.Context, ride *model.Ride) error {
	query := `UPDATE rides SET passenger_id = $1, driver_id = $2, price = $3, status = $4 WHERE id = $5;`
	_, err := db.Pool.Exec(ctx, query, ride.PassengerID, ride.DriverID, ride.Price, ride.Status, ride.ID)
	return err
}

func (db *RideDatabase) DeleteRide(ctx context.Context, id int) error {
	query := `DELETE FROM rides WHERE id = $1;`
	_, err := db.Pool.Exec(ctx, query, id)
	return err
}

func (db *RideDatabase) DeleteAllRides(ctx context.Context) error {
	query := `DELETE FROM rides;`
	_, err := db.Pool.Exec(ctx, query)
	return err
}
