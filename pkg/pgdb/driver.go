package pgdb

import (
	"context"
	"fmt"
	"log"

	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DriverCruder interface {
	CreateDriver(ctx context.Context, driver *model.Driver) error
	ListDrivers(ctx context.Context) ([]model.Driver, error)
	GetDriver(ctx context.Context, id int) (*model.Driver, error)
	UpdateDriver(ctx context.Context, driver *model.Driver) error
	DeleteDriver(ctx context.Context, id int) error
}

type Database struct {
	Pool  *pgxpool.Pool
	Table string
}

func NewDatabase(ctx context.Context, connectionString string, table string) (*Database, error) {
	pool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	db := &Database{
		Pool:  pool,
		Table: table,
	}

	if err := db.createTable(ctx); err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	return db, nil
}

func (db *Database) createTable(ctx context.Context) error {
	query := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        email VARCHAR(100) UNIQUE,
        license_number VARCHAR(50),
        region VARCHAR(100),
		status VARCHAR(50)
    );`, db.Table)
	_, err := db.Pool.Exec(context.Background(), query)
	return err
}

func (db *Database) Close() {
	db.Pool.Close()
}

func (db *Database) CreateDriver(ctx context.Context, driver *model.Driver) error {
	query := fmt.Sprintf(`INSERT INTO %s (name, email, license_number, region, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`, db.Table)
	err := db.Pool.QueryRow(context.Background(), query, driver.Name, driver.Email, driver.LicenseNumber, driver.Region, driver.Status).Scan(&driver.ID)
	return err
}

func (db *Database) ListDrivers(ctx context.Context) ([]model.Driver, error) {
	var drivers []model.Driver
	rows, err := db.Pool.Query(
		context.Background(),
		fmt.Sprintf("SELECT id, name, email, license_number, region, status FROM %s", db.Table),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var driver model.Driver
		if err := rows.Scan(&driver.ID, &driver.Name, &driver.Email, &driver.LicenseNumber, &driver.Region, &driver.Status); err != nil {
			return nil, err
		}
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

func (db *Database) GetDriver(ctx context.Context, id int) (*model.Driver, error) {
	var driver model.Driver
	query := fmt.Sprintf(`
			SELECT 
				id, name, email, license_number, region, status 
			FROM 
				%s 
			WHERE id = $1`,
		db.Table)
	err := db.Pool.QueryRow(context.Background(), query, id).Scan(
		&driver.ID,
		&driver.Name,
		&driver.Email,
		&driver.LicenseNumber,
		&driver.Region,
		&driver.Status)
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

func (db *Database) UpdateDriver(ctx context.Context, driver *model.Driver) error {
	query := fmt.Sprintf(`
		UPDATE 
			%s 
		SET 
			name = $1, email = $2, license_number = $3, region = $4, status=$5 
		WHERE 
			id = $6`,
		db.Table)
	_, err := db.Pool.Exec(context.Background(), query, driver.Name, driver.Email, driver.LicenseNumber, driver.Region, driver.Status, driver.ID)
	return err
}

func (db *Database) DeleteDriver(ctx context.Context, id int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, db.Table)
	_, err := db.Pool.Exec(context.Background(), query, id)
	return err
}

func (db *Database) DeleteAllDrivers(ctx context.Context) error {
	query := fmt.Sprintf(`DELETE FROM %s`, db.Table)
	_, err := db.Pool.Exec(context.Background(), query)
	return err
}
