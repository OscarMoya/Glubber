package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/OscarMoya/Glubber/pkg/queue"
	"github.com/OscarMoya/Glubber/pkg/repository"
	"github.com/stretchr/testify/require"
)

func DeleteRideDB(db *RideService) func() {
	return func() {
		db.DeleteAllRides(context.Background())
		db.Close()
	}

}

var connectionString = "postgresql://admin:admin123@localhost:5432/glubber?sslmode=disable"

func getTestOpts(table string) RideServiceOpts {
	repo, err := repository.NewDBRepository(connectionString, table+"_events")
	if err != nil {
		fmt.Println("Error creating repository")
	}
	producer, err := queue.NewSaramaKafkaProducer([]string{"localhost:9092"})
	if err != nil {
		fmt.Println("Error creating producer")
	}
	return RideServiceOpts{
		Repository:  repo,
		Producer:    producer,
		Table:       table,
		DriverTopic: "drivers" + table,
		DriverKey:   "driver" + table,
	}

}

func TestCreateRide(t *testing.T) {
	// Create a new database
	db, err := NewRideService(context.Background(), getTestOpts("rides1"))
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	not := db.Repository.Notifications(ctx)
	nots := 0
	for n := range not {
		fmt.Printf("Notification: %v\n", n)
		nots++
		if nots == 1 {
			break
		}
	}
	require.Equal(t, 1, nots)
}

func TestListRides(t *testing.T) {
	// Create a new database
	db, err := NewRideService(context.Background(), getTestOpts("rides2"))
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

	rides, err := db.ListRides(context.Background())
	require.NoError(t, err)
	require.Len(t, rides, 2)
	require.Equal(t, ride.ID, rides[0].ID)
	require.Equal(t, ride2.ID, rides[1].ID)

}

func TestGetRide(t *testing.T) {

	db, err := NewRideService(context.Background(), getTestOpts("rides3"))
	require.NoError(t, err)
	defer DeleteRideDB(db)()

	ride := &model.Ride{
		PassengerID: 20,
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
	db, err := NewRideService(context.Background(), getTestOpts("rides4"))
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
	db, err := NewRideService(context.Background(), getTestOpts("rides5"))
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
