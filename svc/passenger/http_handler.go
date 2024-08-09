package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/gorilla/mux"
)

func createPassengerHandler(serviceData *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Passenger model.Passenger
		if err := json.NewDecoder(r.Body).Decode(&Passenger); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		err := serviceData.PGDB.CreatePassenger(ctx, &Passenger)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Passenger)
	}

}

func listPassengersHandler(serviceStatus *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		Passengers, err := serviceStatus.PGDB.ListPassengers(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Passengers)
	}
}

func getPassengerHandler(serviceData *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		Passenger, err := serviceData.PGDB.GetPassenger(ctx, idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Passenger)

	}
}

func updatePassengerHandler(serviceData *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Passenger model.Passenger
		if err := json.NewDecoder(r.Body).Decode(&Passenger); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		err := serviceData.PGDB.UpdatePassenger(ctx, &Passenger)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Passenger)

	}
}

func deletePassengerHandler(serviceStatus *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		err = serviceStatus.PGDB.DeletePassenger(ctx, idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
