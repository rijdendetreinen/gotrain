package api

import (
	"encoding/json"
	"net/http"

	"github.com/rijdendetreinen/gotrain/stores"
)

func arrivalCounters(w http.ResponseWriter, r *http.Request) {
	response := Statistics{stores.Stores.ArrivalStore.Counters, stores.Stores.ArrivalStore.GetNumberOfArrivals()}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func arrivalsAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores.Stores.ArrivalStore.GetAllArrivals())
}
