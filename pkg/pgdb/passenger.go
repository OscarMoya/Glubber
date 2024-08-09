package pgdb

import (
	"context"

	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PassengerCruder interface {
	CreatePassenger(ctx context.Context, passenger *model.Passenger) error
	ListPassengers(ctx context.Context) ([]model.Passenger, error)
	GetPassenger(ctx context.Context, id int) (*model.Passenger, error)
	UpdatePassenger(ctx context.Context, passenger *model.Passenger) error
	DeletePassenger(ctx context.Context, id int) error
}

type PassengerDatabase struct {
	Pool  *pgxpool.Pool
	Table string
}

func NewPassengerDatabase(ctx context.Context, connectionString string, table string) (*PassengerDatabase, error) {
	pool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	db := &PassengerDatabase{
		Pool:  pool,
		Table: table,
	}

	if err := db.createTable(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *PassengerDatabase) createTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS passengers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100) UNIQUE
	);`
	_, err := db.Pool.Exec(ctx, query)
	return err
}

func (db *PassengerDatabase) Close() {
	db.Pool.Close()
}

func (db *PassengerDatabase) CreatePassenger(ctx context.Context, passenger *model.Passenger) error {
	query := `INSERT INTO passengers (name, email) VALUES ($1, $2) RETURNING id;`
	err := db.Pool.QueryRow(ctx, query, passenger.Name, passenger.Email).Scan(&passenger.ID)
	return err
}

func (db *PassengerDatabase) ListPassengers(ctx context.Context) ([]model.Passenger, error) {
	query := `SELECT id, name, email FROM passengers;`
	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	passengers := []model.Passenger{}
	for rows.Next() {
		var passenger model.Passenger
		err := rows.Scan(&passenger.ID, &passenger.Name, &passenger.Email)
		if err != nil {
			return nil, err
		}
		passengers = append(passengers, passenger)
	}

	return passengers, nil
}

func (db *PassengerDatabase) GetPassenger(ctx context.Context, id int) (*model.Passenger, error) {
	query := `SELECT id, name, email FROM passengers WHERE id = $1;`
	row := db.Pool.QueryRow(ctx, query, id)

	var passenger model.Passenger
	err := row.Scan(&passenger.ID, &passenger.Name, &passenger.Email)
	if err != nil {
		return nil, err
	}

	return &passenger, nil
}

func (db *PassengerDatabase) UpdatePassenger(ctx context.Context, passenger *model.Passenger) error {
	query := `UPDATE passengers SET name = $1, email = $2 WHERE id = $3;`
	_, err := db.Pool.Exec(ctx, query, passenger.Name, passenger.Email, passenger.ID)
	return err
}

func (db *PassengerDatabase) DeletePassenger(ctx context.Context, id int) error {
	query := `DELETE FROM passengers WHERE id = $1;`
	_, err := db.Pool.Exec(ctx, query, id)
	return err
}

func (db *PassengerDatabase) DeleteAllPassengers(ctx context.Context) error {
	query := `DELETE FROM passengers;`
	_, err := db.Pool.Exec(ctx, query)
	return err
}
