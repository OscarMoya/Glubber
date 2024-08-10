package main

import (
	"log"
	"net/http"

	"github.com/OscarMoya/Glubber/pkg/service"
	"github.com/gorilla/mux"
)

type ServiceData struct {
	PGDB *service.PassengerDatabase
}

const (
	serveURL         = "localhost:8082"
	passengerWSURI   = "ws/v1/passenger"
	passengerHTTPUri = "v1/passengers"
)

func main() {

	serviceStatus := &ServiceData{}
	serviceStatus.PGDB = &service.PassengerDatabase{}
	r := mux.NewRouter()

	// HTTP Handlers
	r.HandleFunc(passengerHTTPUri, createPassengerHandler(serviceStatus)).Methods("POST")
	r.HandleFunc(passengerHTTPUri, listPassengersHandler(serviceStatus)).Methods("GET")
	r.HandleFunc(passengerHTTPUri+"/{id}", getPassengerHandler(serviceStatus)).Methods("GET")
	r.HandleFunc(passengerHTTPUri+"/{id}", updatePassengerHandler(serviceStatus)).Methods("PUT")
	r.HandleFunc(passengerHTTPUri+"/{id}", deletePassengerHandler(serviceStatus)).Methods("DELETE")

	// WebSocket Handlers TODO: Implement this to exchange messages to the passenger to get updates on the ride or
	// get the driver location
	/*r.HandleFunc(passengerWSURI, func(w http.ResponseWriter, r *http.Request) {
		handlePassengerConnections(w, r, serviceStatus)
	})
	*/

	log.Printf("HTTP server started on %s\n", serveURL)
	err := http.ListenAndServe(serveURL, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
