package api

import (
	"encoding/json"
	"net/http"

	"github.com/rijdendetreinen/gotrain/stores"
)

func serviceCounters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores.Stores.ServiceStore.Counters)
}

func serviceAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores.Stores.ServiceStore.GetAllServices())
}
