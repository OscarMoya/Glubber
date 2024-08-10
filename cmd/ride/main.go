package main

import (
	"context"
	"log"
	"net/http"

	"github.com/OscarMoya/Glubber/pkg/billing"
	"github.com/OscarMoya/Glubber/pkg/queue"
	"github.com/OscarMoya/Glubber/pkg/repository"
	"github.com/OscarMoya/Glubber/pkg/service"
	"github.com/gorilla/mux"
)

const (

	// serveURL is the URL where the server will be listening
	serveURL    = "localhost:8083"
	rideWSURI   = "ws/v1/ride"
	rideHTTPUri = "v1/rides"
)

// ServiceData is the struct that holds the database connection
type ServiceData struct {
	PGDB   *service.RideService
	Biller billing.Biller
}

func main() {
	repository, err := repository.NewDBRepository("postgresql://admin:admin123@localhost:5432/glubber?sslmode=disable", "rides_events")
	if err != nil {
		log.Fatal(err)
	}
	defer repository.CloseListener(context.Background())
	producer, err := queue.NewSaramaKafkaProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	biller := billing.NewSimpleBiller(2.0, 1.0)

	// Create a new service
	riderOpts := service.RideServiceOpts{
		Repository:  repository,
		Producer:    producer,
		Biller:      biller,
		Table:       "rides",
		DriverTopic: "drivers",
		DriverKey:   "driver",
	}

	pgdb, err := service.NewRideService(context.Background(), riderOpts)
	if err != nil {
		log.Fatal(err)
	}
	serviceData := &ServiceData{}
	serviceData.PGDB = pgdb

	r := mux.NewRouter()

	// HTTP Handlers
	r.HandleFunc(rideHTTPUri, createRideHandler(serviceData)).Methods("POST")
	r.HandleFunc(rideHTTPUri, listRidesHandler(serviceData)).Methods("GET")
	r.HandleFunc(rideHTTPUri+"/{id}", getRideHandler(serviceData)).Methods("GET")
	r.HandleFunc(rideHTTPUri+"/{id}", updateRideHandler(serviceData)).Methods("PUT")
	r.HandleFunc(rideHTTPUri+"/{id}", deleteRideHandler(serviceData)).Methods("DELETE")

	log.Printf("HTTP server started on %s\n", serveURL)
	err = http.ListenAndServe(serveURL, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
