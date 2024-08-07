package pgdb

import (
	"context"
	"testing"

	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/stretchr/testify/require"
)

func DeleteDB(db *Database) func() {
	return func() {
		db.DeleteAllDrivers(context.Background())
		db.Close()
	}

}

func TestShouldCreateDriver(t *testing.T) {
	// Create a new database
	db, err := NewDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "drivers1")
	require.NoError(t, err)
	defer DeleteDB(db)()

	// Create a driver
	driver := &model.Driver{
		Name:          "John Doe",
		Email:         "john@doe.com",
		LicenseNumber: "123456",
		Region:        "US",
		Status:        "active",
	}

	err = db.CreateDriver(context.Background(), driver)
	require.NoError(t, err)
	require.NotZero(t, driver.ID)

	err = db.CreateDriver(context.Background(), driver)
	require.Error(t, err)
}

func TestShouldListDrivers(t *testing.T) {
	// Create a new database
	db, err := NewDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "drivers2")
	require.NoError(t, err)
	defer DeleteDB(db)()

	// Create drivers
	driver := &model.Driver{
		Name:          "John Doe",
		Email:         "john@doe.com",
		LicenseNumber: "123456",
		Region:        "US",
		Status:        "active",
	}

	// Create drivers
	driver2 := &model.Driver{
		Name:          "John Doe2",
		Email:         "john2@doe.com",
		LicenseNumber: "123456",
		Region:        "US",
		Status:        "active",
	}

	err = db.CreateDriver(context.Background(), driver)
	require.NoError(t, err)
	err = db.CreateDriver(context.Background(), driver2)
	require.NoError(t, err)

	drivers, err := db.ListDrivers(context.Background())
	require.NoError(t, err)
	require.Len(t, drivers, 2)
}

func TestShouldGetDriver(t *testing.T) {
	// Create a new database
	db, err := NewDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "drivers3")
	require.NoError(t, err)
	defer DeleteDB(db)()

	// Create a driver
	driver := &model.Driver{
		Name:          "John Doe",
		Email:         "john@doe.com",
		LicenseNumber: "123456",
		Region:        "US",
		Status:        "active",
	}

	err = db.CreateDriver(context.Background(), driver)

	require.NoError(t, err)
	require.NotZero(t, driver.ID)

	driver2, err := db.GetDriver(context.Background(), driver.ID)
	require.NoError(t, err)
	require.Equal(t, driver.ID, driver2.ID)
	require.Equal(t, driver.Name, driver2.Name)
	require.Equal(t, driver.Email, driver2.Email)
	require.Equal(t, driver.LicenseNumber, driver2.LicenseNumber)
	require.Equal(t, driver.Region, driver2.Region)
	require.Equal(t, driver.Status, driver2.Status)
}

func TestShouldUpdateDriver(t *testing.T) {
	// Create a new database
	db, err := NewDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "drivers4")
	require.NoError(t, err)
	defer DeleteDB(db)()

	// Create a driver
	driver := &model.Driver{
		Name:          "John Doe",
		Email:         "john@doe.com",
		LicenseNumber: "123456",
		Region:        "US",
		Status:        "active",
	}

	err = db.CreateDriver(context.Background(), driver)
	require.NoError(t, err)

	driver.Name = "Jane Doe"

	err = db.UpdateDriver(context.Background(), driver)
	require.NoError(t, err)

	driver2, err := db.GetDriver(context.Background(), driver.ID)
	require.NoError(t, err)
	require.Equal(t, driver.ID, driver2.ID)
	require.Equal(t, driver.Name, driver2.Name)

}

func TestShouldDeleteDriver(t *testing.T) {
	// Create a new database
	db, err := NewDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "drivers5")
	require.NoError(t, err)
	defer DeleteDB(db)()

	// Create a driver
	driver := &model.Driver{
		Name:          "John Doe",
		Email:         "john@doe.com",
		LicenseNumber: "123456",
		Region:        "US",
		Status:        "active",
	}

	err = db.CreateDriver(context.Background(), driver)
	require.NoError(t, err)

	err = db.DeleteDriver(context.Background(), driver.ID)
	require.NoError(t, err)

	_, err = db.GetDriver(context.Background(), driver.ID)
	require.Error(t, err)

}
