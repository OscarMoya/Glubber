package model

import (
	"github.com/jackc/pgx/v4"
)

type Driver struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	LicenseNumber string `json:"license_number"`
	Region        string `json:"location"`
	Status        string `json:"status"`
}

func (d *Driver) Scan(row pgx.Row) error {
	return row.Scan(
		&d.ID,
		&d.Name,
		&d.Email,
		&d.LicenseNumber,
		&d.Region,
		&d.Status,
	)
}
