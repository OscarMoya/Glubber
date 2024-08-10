package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/gorilla/mux"
)

func createDriverHandler(serviceData *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var driver model.Driver
		if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		err := serviceData.PGDB.CreateDriver(ctx, &driver)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(driver)
	}

}

func listDriversHandler(serviceStatus *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		drivers, err := serviceStatus.PGDB.ListDrivers(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(drivers)
	}
}

func getDriverHandler(serviceData *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		driver, err := serviceData.PGDB.GetDriver(ctx, idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(driver)

	}
}

func updateDriverHandler(serviceData *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var driver model.Driver
		if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		err := serviceData.PGDB.UpdateDriver(ctx, &driver)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(driver)

	}
}

func deleteDriverHandler(serviceStatus *ServiceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		err = serviceStatus.PGDB.DeleteDriver(ctx, idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
