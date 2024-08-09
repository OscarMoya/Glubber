package pgdb

import (
	"context"
	"testing"

	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/stretchr/testify/require"
)

func DeletePassengerDB(db *PassengerDatabase) func() {
	return func() {
		db.DeleteAllPassengers(context.Background())
		db.Close()
	}

}

func TestCreatePassenger(t *testing.T) {
	db, err := NewPassengerDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "passengers1")
	require.NoError(t, err)
	defer DeletePassengerDB(db)()

	passenger := &model.Passenger{
		Name:  "John Doe",
		Email: "john@dow.com",
	}
	passenger2 := &model.Passenger{
		Name:  "John Doe2",
		Email: "john@dow2.com",
	}

	err = db.CreatePassenger(context.Background(), passenger)
	require.NoError(t, err)
	require.NotZero(t, passenger.ID)

	err = db.CreatePassenger(context.Background(), passenger2)
	require.NoError(t, err)
	require.NotZero(t, passenger2.ID)

	err = db.CreatePassenger(context.Background(), passenger)
	require.Error(t, err)

}

func TestListPassengers(t *testing.T) {
	db, err := NewPassengerDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "passengers2")
	require.NoError(t, err)
	defer DeletePassengerDB(db)()

	passenger := &model.Passenger{
		Name:  "John Doe",
		Email: "john@dow",
	}

	passenger2 := &model.Passenger{
		Name:  "John Doe2",
		Email: "john@dow2",
	}

	err = db.CreatePassenger(context.Background(), passenger)
	require.NoError(t, err)

	err = db.CreatePassenger(context.Background(), passenger2)
	require.NoError(t, err)

	passengers, err := db.ListPassengers(context.Background())
	require.NoError(t, err)

	require.Len(t, passengers, 2)

}

func TestGetPassenger(t *testing.T) {

	db, err := NewPassengerDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "passengers3")
	require.NoError(t, err)
	defer DeletePassengerDB(db)()

	passenger := &model.Passenger{
		Name:  "John Doe",
		Email: "john@dow",
	}

	err = db.CreatePassenger(context.Background(), passenger)
	require.NoError(t, err)

	passenger2, err := db.GetPassenger(context.Background(), passenger.ID)
	require.NoError(t, err)
	require.Equal(t, passenger.ID, passenger2.ID)
	require.Equal(t, passenger.Name, passenger2.Name)
	require.Equal(t, passenger.Email, passenger2.Email)

}

func TestUpdatePassenger(t *testing.T) {
	db, err := NewPassengerDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "passengers4")
	require.NoError(t, err)
	defer DeletePassengerDB(db)()

	passenger := &model.Passenger{
		Name:  "John Doe",
		Email: "john@dow",
	}

	err = db.CreatePassenger(context.Background(), passenger)
	require.NoError(t, err)

	passenger.Name = "John Doe Updated"
	passenger.Email = "john@doe.com"

	err = db.UpdatePassenger(context.Background(), passenger)
	require.NoError(t, err)

	passenger2, err := db.GetPassenger(context.Background(), passenger.ID)
	require.NoError(t, err)

	require.Equal(t, passenger.ID, passenger2.ID)
	require.Equal(t, passenger.Name, passenger2.Name)
	require.Equal(t, passenger.Email, passenger2.Email)

}

func TestDeletePassenger(t *testing.T) {
	db, err := NewPassengerDatabase(context.Background(), "postgresql://admin:admin123@localhost:5432/glubber", "passengers5")
	require.NoError(t, err)
	defer DeletePassengerDB(db)()

	passenger := &model.Passenger{
		Name:  "John Doe",
		Email: "john@dow",
	}

	err = db.CreatePassenger(context.Background(), passenger)
	require.NoError(t, err)

	err = db.DeletePassenger(context.Background(), passenger.ID)
	require.NoError(t, err)

	passenger2, err := db.GetPassenger(context.Background(), passenger.ID)
	require.Error(t, err)
	require.Nil(t, passenger2)

}
