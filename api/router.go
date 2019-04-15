package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// ServeAPI serves the REST API on the given address
func ServeAPI(address string, exit chan bool) {
	srv := &http.Server{Addr: address}
	router := mux.NewRouter()

	// Router paths:
	router.HandleFunc("/version", apiVersion).Methods("GET")
	router.HandleFunc("/v1", apiVersion).Methods("GET")
	router.HandleFunc("/v2", apiVersion).Methods("GET")
	router.HandleFunc("/v2/version", apiVersion).Methods("GET")

	router.HandleFunc("/v2/arrivals/stats", arrivalCounters).Methods("GET")
	router.HandleFunc("/v2/arrivals/all", arrivalsAll).Methods("GET")
	router.HandleFunc("/v2/arrivals/station/{station}", arrivalsStation).Methods("GET")
	router.HandleFunc("/v2/arrivals/arrival/{id}/{station}/{date}", arrivalDetails).Methods("GET")

	router.HandleFunc("/v2/departures/stats", departureCounters).Methods("GET")
	router.HandleFunc("/v2/departures/all", departuresAll).Methods("GET")
	router.HandleFunc("/v2/departures/station/{station}", departuresStation).Methods("GET")
	router.HandleFunc("/v2/departures/departure/{id}/{station}/{date}", departureDetails).Methods("GET")

	router.HandleFunc("/v2/services/stats", serviceCounters).Methods("GET")
	router.HandleFunc("/v2/services/all", serviceAll).Methods("GET")
	router.HandleFunc("/v2/services/service/{id}/{date}", serviceDetails).Methods("GET")

	srv.Handler = router

	go listenAndServe(srv, exit)
	log.WithField("address", address).Info("REST API started")

	<-exit
	log.Info("Shutting down REST API")
	srv.Close()
	log.Info("REST API shut down")
	exit <- true
}

func listenAndServe(srv *http.Server, exit chan bool) {
	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.WithError(err).Fatal("REST API fatal error")
		}
	}
}

func apiVersion(w http.ResponseWriter, r *http.Request) {
	version := map[string]int{
		"version": 2,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}
