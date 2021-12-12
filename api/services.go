package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rijdendetreinen/gotrain/models"
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
	json.NewEncoder(w).Encode(wrapServicesStatus("service", ServiceToJSON(*service, language, verbose)))
}

// ServiceToJSON generates an interface (convertible to JSON) with all service details
func ServiceToJSON(service models.Service, language string, verbose bool) map[string]interface{} {
	response := map[string]interface{}{
		"id":             service.ID,
		"timestamp":      service.Timestamp,
		"service_date":   service.ServiceDate,
		"service_number": service.ServiceNumber,
		"type":           service.ServiceType,
		"type_code":      service.ServiceTypeCode,
		"line_number":    nullString(service.LineNumber),
		"company":        service.Company,

		"journey_planner":      service.JourneyPlanner,
		"reservation_required": service.ReservationRequired,
		"special_ticket":       service.SpecialTicket,
		"with_supplement":      service.WithSupplement,

		"parts":   []interface{}{},
		"remarks": models.GetRemarks(service.Modifications, language),
		"tips":    []interface{}{},
	}

	responseParts := []interface{}{}

	for _, part := range service.ServiceParts {
		partResponse := map[string]interface{}{
			"service_number": part.ServiceNumber,
			"remarks":        models.GetRemarks(part.Modifications, language),
			"tips":           []interface{}{},
			"stops":          []interface{}{},
		}

		var stops []models.ServiceStop

		if verbose {
			stops = part.Stops
		} else {
			stops = part.GetStoppingStations()
		}

		responseStops := []interface{}{}

		for _, stop := range stops {
			responseStops = append(responseStops, serviceStopToJSON(stop, language, verbose))
		}

		partResponse["stops"] = responseStops
		responseParts = append(responseParts, partResponse)
	}

	response["parts"] = responseParts

	return response
}

func serviceStopToJSON(stop models.ServiceStop, language string, verbose bool) map[string]interface{} {
	stopResponse := map[string]interface{}{
		"station":              stop.Station,
		"station_accessible":   stop.StationAccessible,
		"assistance_available": stop.AssistanceAvailable,
		"stopping_actual":      stop.StoppingActual,
		"stopping_planned":     stop.StoppingPlanned,
		"stop_type":            stop.StopType,
		"do_not_board":         stop.DoNotBoard,

		"arrival_time":             localTimeString(stop.ArrivalTime),
		"arrival_platform_actual":  nullString(stop.ArrivalPlatformActual),
		"arrival_platform_planned": nullString(stop.ArrivalPlatformPlanned),
		"arrival_delay":            stop.ArrivalDelay,
		"arrival_cancelled":        stop.ArrivalCancelled,

		"departure_time":             localTimeString(stop.DepartureTime),
		"departure_platform_actual":  nullString(stop.DeparturePlatformActual),
		"departure_platform_planned": nullString(stop.DeparturePlatformPlanned),
		"departure_delay":            stop.DepartureDelay,
		"departure_cancelled":        stop.DepartureCancelled,

		"remarks":  models.GetRemarks(stop.Modifications, language),
		"tips":     []interface{}{},
		"material": materialsToJSON(stop.Material, language, verbose),
	}

	return stopResponse
}

func wrapServicesStatus(key string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status": stores.Stores.ServiceStore.Status,
		key:      data,
	}
}
