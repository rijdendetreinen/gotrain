package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rijdendetreinen/gotrain/stores"
)

func serviceCounters(w http.ResponseWriter, r *http.Request) {
	response := Statistics{
		stores.Stores.ServiceStore.Counters,
		stores.Stores.ServiceStore.GetNumberOfServices(),
		stores.Stores.ServiceStore.Status,
		stores.Stores.ServiceStore.LastStatusChange,
		stores.Stores.ServiceStore.MessagesAverage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func serviceAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores.Stores.ServiceStore.GetAllServices())
}

func serviceDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	serviceID := vars["id"]
	serviceDate := vars["date"]
	language := getLanguageVar(r.URL)
	verbose := getBooleanQueryParameter(r.URL, "verbose", false)

	service := stores.Stores.ServiceStore.GetService(serviceID, serviceDate)

	if service == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceToJSON(*service, language, verbose))
}
