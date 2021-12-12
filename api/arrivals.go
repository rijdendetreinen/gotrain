package api

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/rijdendetreinen/gotrain/models"
	"github.com/rijdendetreinen/gotrain/stores"
)

func arrivalCounters(w http.ResponseWriter, r *http.Request) {
	response := Statistics{
		stores.Stores.ArrivalStore.Counters,
		stores.Stores.ArrivalStore.GetNumberOfArrivals(),
		stores.Stores.ArrivalStore.Status,
		stores.Stores.ArrivalStore.LastStatusChange,
		stores.Stores.ArrivalStore.MessagesAverage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func arrivalsStation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	station := vars["station"]
	language := getLanguageVar(r.URL)

	arrivals := stores.Stores.ArrivalStore.GetStationArrivals(station, false)

	// Sort arrivals on arrival time, or on planned origin when arrival times are equal
	sort.Slice(arrivals, func(i, j int) bool {
		if arrivals[i].ArrivalTime.Equal(arrivals[j].ArrivalTime) {
			return arrivals[i].PlannedOriginString() < arrivals[j].PlannedOriginString()
		}

		return arrivals[i].ArrivalTime.Before(arrivals[j].ArrivalTime)
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wrapArrivalsStatus("arrivals", arrivalsToJSON(arrivals, language)))
}

func arrivalDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	serviceID := vars["id"]
	serviceDate := vars["date"]
	station := vars["station"]
	language := getLanguageVar(r.URL)

	arrival := stores.Stores.ArrivalStore.GetArrival(serviceID, serviceDate, station)

	if arrival == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wrapArrivalsStatus("arrival", arrivalToJSON(*arrival, language)))
}

func arrivalsToJSON(arrivals []models.Arrival, language string) []map[string]interface{} {
	response := make([]map[string]interface{}, 0)

	for _, arrival := range arrivals {
		response = append(response, arrivalToJSON(arrival, language))
	}

	return response
}

func arrivalToJSON(arrival models.Arrival, language string) map[string]interface{} {
	response := map[string]interface{}{
		"service_id":          arrival.ServiceID,
		"name":                nullString(arrival.ServiceName),
		"line_number":         nullString(arrival.LineNumber),
		"timestamp":           arrival.Timestamp,
		"status":              arrival.Status,
		"service_date":        arrival.ServiceDate,
		"service_number":      arrival.ServiceNumber,
		"station":             arrival.Station.Code,
		"type":                arrival.ServiceType,
		"type_code":           arrival.ServiceTypeCode,
		"company":             arrival.Company,
		"origin_actual":       nullString(arrival.ActualOriginString()),
		"origin_planned":      nullString(arrival.PlannedOriginString()),
		"origin_actual_codes": arrival.ActualOriginCodes(),
		"via":                 nullString(arrival.ViaStationsString()),
		"arrival_time":        localTimeString(arrival.ArrivalTime),
		"platform_actual":     nullString(arrival.PlatformActual),
		"platform_planned":    nullString(arrival.PlatformPlanned),
		"delay":               arrival.Delay,
		"cancelled":           arrival.Cancelled,
		"platform_changed":    arrival.PlatformChanged(),

		"remarks": []interface{}{},
	}

	response["remarks"] = models.GetRemarks(arrival.Modifications, language)

	return response
}

func wrapArrivalsStatus(key string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status": stores.Stores.ArrivalStore.Status,
		key:      data,
	}
}
