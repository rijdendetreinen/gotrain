package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// ServeAPI serves the REST API on the given address
func ServeAPI(address string) {
	router := mux.NewRouter()

	// Router paths:
	router.HandleFunc("/version", apiVersion).Methods("GET")
	router.HandleFunc("/v1", apiVersion).Methods("GET")
	router.HandleFunc("/v2", apiVersion).Methods("GET")
	router.HandleFunc("/v2/version", apiVersion).Methods("GET")
	router.HandleFunc("/v2/services/stats", serviceCounters).Methods("GET")
	router.HandleFunc("/v2/services/all", serviceAll).Methods("GET")

	log.WithField("address", address).Info("REST API started")

	serverError := http.ListenAndServe(address, router)
	log.WithField("error", serverError).Fatal("REST API fatal error")
}

func apiVersion(w http.ResponseWriter, r *http.Request) {
	version := map[string]int{
		"version": 2,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}
