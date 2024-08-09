package main

import (
	"log"
	"net/http"

	"github.com/OscarMoya/Glubber/pkg/pgdb"
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
	PGDB *pgdb.RideDatabase
}

func main() {

	serviceData := &ServiceData{}
	serviceData.PGDB = &pgdb.RideDatabase{}
	r := mux.NewRouter()

	// HTTP Handlers
	r.HandleFunc(rideHTTPUri, createRideHandler(serviceData)).Methods("POST")
	r.HandleFunc(rideHTTPUri, listRidesHandler(serviceData)).Methods("GET")
	r.HandleFunc(rideHTTPUri+"/{id}", getRideHandler(serviceData)).Methods("GET")
	r.HandleFunc(rideHTTPUri+"/{id}", updateRideHandler(serviceData)).Methods("PUT")
	r.HandleFunc(rideHTTPUri+"/{id}", deleteRideHandler(serviceData)).Methods("DELETE")

	log.Printf("HTTP server started on %s\n", serveURL)
	err := http.ListenAndServe(serveURL, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
