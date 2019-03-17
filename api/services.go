package api

import (
	"encoding/json"
	"net/http"

	"github.com/rijdendetreinen/gotrain/stores"
)

func serviceCounters(w http.ResponseWriter, r *http.Request) {
	response := Statistics{stores.Stores.ServiceStore.Counters, stores.Stores.ServiceStore.GetNumberOfServices()}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func serviceAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores.Stores.ServiceStore.GetAllServices())
}
