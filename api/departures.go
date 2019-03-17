package api

import (
	"encoding/json"
	"net/http"

	"github.com/rijdendetreinen/gotrain/stores"
)

func departureCounters(w http.ResponseWriter, r *http.Request) {
	response := Statistics{stores.Stores.DepartureStore.Counters, stores.Stores.DepartureStore.GetNumberOfDepartures()}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func departuresAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores.Stores.DepartureStore.GetAllDepartures())
}
