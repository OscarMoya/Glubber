package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/gorilla/mux"
)

func createRideHandler(serviceData *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ride model.Ride
		if err := json.NewDecoder(r.Body).Decode(&ride); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		ride.Status = model.RideStatusPending
		err := serviceData.Biller.EstimateRide(&ride)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = serviceData.PGDB.CreateRide(ctx, &ride)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ride)
	}

}

func listRidesHandler(serviceData *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		Rides, err := serviceData.PGDB.ListRides(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Rides)
	}
}

func getRideHandler(serviceData *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		Ride, err := serviceData.PGDB.GetRide(ctx, idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Ride)

	}
}

func updateRideHandler(serviceData *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Ride model.Ride
		if err := json.NewDecoder(r.Body).Decode(&Ride); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		err := serviceData.PGDB.UpdateRide(ctx, &Ride)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Ride)
	}
}

func deleteRideHandler(serviceData *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		err = serviceData.PGDB.DeleteRide(ctx, idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
