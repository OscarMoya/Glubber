package pgdb

import (
	"context"
	"testing"

	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/stretchr/testify/require"
)

func DeleteRideDB(db *RideDatabase) func() {
	return func() {
		db.DeleteAllRides(context.Background())
		db.Close()
	}

}

func TestCreateRide(t *testing.T) {
	// Create a new database
	db, err := NewRideDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "rides1")
	require.NoError(t, err)
	defer DeleteRideDB(db)()

	// Create a ride
	ride := &model.Ride{
		PassengerID: 1,
		Price:       100,
		Status:      "active",
	}

	err = db.CreateRide(context.Background(), ride)
	require.NoError(t, err)
	require.NotZero(t, ride.ID)

	err = db.CreateRide(context.Background(), ride)
	require.Error(t, err)
}

func TestListRides(t *testing.T) {
	// Create a new database
	db, err := NewRideDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "rides2")
	require.NoError(t, err)
	defer DeleteRideDB(db)()

	// Create a ride
	ride := &model.Ride{
		PassengerID: 1,
		Price:       100,
		Status:      "active",
	}

	// Create a ride
	ride2 := &model.Ride{
		PassengerID: 2,
		Price:       200,
		Status:      "active",
	}

	err = db.CreateRide(context.Background(), ride)
	require.NoError(t, err)
	require.NotZero(t, ride.ID)

	err = db.CreateRide(context.Background(), ride2)
	require.NoError(t, err)
	require.NotZero(t, ride2.ID)
}

func TestGetRide(t *testing.T) {

	db, err := NewRideDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "rides3")
	require.NoError(t, err)
	defer DeleteRideDB(db)()

	ride := &model.Ride{
		PassengerID: 1,
		Price:       100,
		Status:      "active",
	}

	err = db.CreateRide(context.Background(), ride)
	require.NoError(t, err)

	ride2, err := db.GetRide(context.Background(), ride.ID)
	require.NoError(t, err)
	require.Equal(t, ride.ID, ride2.ID)
	require.Equal(t, ride.PassengerID, ride2.PassengerID)
	require.Equal(t, ride.DriverID, ride2.DriverID)
	require.Equal(t, ride.Price, ride2.Price)
	require.Equal(t, ride.Status, ride2.Status)
}

func TestUpdateRide(t *testing.T) {
	db, err := NewRideDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "rides4")
	require.NoError(t, err)
	defer DeleteRideDB(db)()

	ride := &model.Ride{
		PassengerID: 1,
		Price:       100,
		Status:      "active",
	}

	err = db.CreateRide(context.Background(), ride)
	require.NoError(t, err)

	ride.Status = "inactive"
	err = db.UpdateRide(context.Background(), ride)
	require.NoError(t, err)

	ride2, err := db.GetRide(context.Background(), ride.ID)
	require.NoError(t, err)
	require.Equal(t, ride.ID, ride2.ID)
	require.Equal(t, ride.PassengerID, ride2.PassengerID)
	require.Equal(t, ride.DriverID, ride2.DriverID)
	require.Equal(t, ride.Price, ride2.Price)
	require.Equal(t, ride.Status, ride2.Status)
}

func TestDeleteRide(t *testing.T) {
	db, err := NewRideDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "rides5")
	require.NoError(t, err)
	defer DeleteRideDB(db)()

	ride := &model.Ride{
		PassengerID: 1,
		Price:       100,
		Status:      "active",
	}

	err = db.CreateRide(context.Background(), ride)
	require.NoError(t, err)

	err = db.DeleteRide(context.Background(), ride.ID)
	require.NoError(t, err)

	_, err = db.GetRide(context.Background(), ride.ID)
	require.Error(t, err)
}
