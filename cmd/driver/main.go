package main

import (
	"log"
	"net/http"

	"github.com/OscarMoya/Glubber/pkg/authentication"
	"github.com/OscarMoya/Glubber/pkg/location"
	"github.com/OscarMoya/Glubber/pkg/service"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	serveURL     = "localhost:8081"
	driverWSURI  = "ws/v1/driver"
	drverHTTPUri = "v1/drivers"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	locationURI   = driverWSURI
	mainURI       = drverHTTPUri
	mainUriWithID = drverHTTPUri + "/{id}"
)

type ServiceData struct {
	Authenticator authentication.DriverAuthenticator
	GeoService    location.LocationManager
	PGDB          service.DriverCruder
}

func main() {
	serviceStatus := &ServiceData{}
	serviceStatus.GeoService = &location.RedisLocationService{}
	serviceStatus.Authenticator = &authentication.JWTDriverAuthenticationService{}
	r := mux.NewRouter()

	// HTTP Handlers
	r.HandleFunc(mainURI, createDriverHandler(serviceStatus)).Methods("POST")
	r.HandleFunc(mainURI, listDriversHandler(serviceStatus)).Methods("GET")
	r.HandleFunc(mainUriWithID, getDriverHandler(serviceStatus)).Methods("GET")
	r.HandleFunc(mainUriWithID, updateDriverHandler(serviceStatus)).Methods("PUT")
	r.HandleFunc(mainUriWithID, deleteDriverHandler(serviceStatus)).Methods("DELETE")

	// WebSocket Handlers
	r.HandleFunc(locationURI, func(w http.ResponseWriter, r *http.Request) {
		handleDriverConnections(w, r, serviceStatus)
	})

	log.Printf("HTTP server started on %s\n", serveURL)
	err := http.ListenAndServe(serveURL, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
