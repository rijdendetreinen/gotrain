package api

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/rijdendetreinen/gotrain/models"
	"github.com/rijdendetreinen/gotrain/stores"
)

func departureCounters(w http.ResponseWriter, r *http.Request) {
	response := Statistics{
		stores.Stores.DepartureStore.Counters,
		stores.Stores.DepartureStore.GetNumberOfDepartures(),
		stores.Stores.DepartureStore.Status,
		stores.Stores.DepartureStore.LastStatusChange,
		stores.Stores.DepartureStore.MessagesAverage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func departuresAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores.Stores.DepartureStore.GetAllDepartures())
}

func departuresStation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	station := vars["station"]
	language := getLanguageVar(r.URL)
	verbose := getBooleanQueryParameter(r.URL, "verbose", false)

	departures := stores.Stores.DepartureStore.GetStationDepartures(station, false)

	// Sort departures on departure time, or on planned destination when departure times are equal
	sort.Slice(departures, func(i, j int) bool {
		if departures[i].DepartureTime.Equal(departures[j].DepartureTime) {
			return departures[i].PlannedDestinationString() < departures[j].PlannedDestinationString()
		}

		return departures[i].DepartureTime.Before(departures[j].DepartureTime)
	})

	if departures == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(departuresToJSON(departures, language, verbose))
}

func departureDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	serviceID := vars["id"]
	serviceDate := vars["date"]
	station := vars["station"]
	language := getLanguageVar(r.URL)
	verbose := getBooleanQueryParameter(r.URL, "verbose", false)

	departure := stores.Stores.DepartureStore.GetDeparture(serviceID, serviceDate, station)

	if departure == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(departureToJSON(*departure, language, verbose))
}

func departuresToJSON(departures []models.Departure, language string, verbose bool) []map[string]interface{} {
	var response []map[string]interface{}

	for _, departure := range departures {
		response = append(response, departureToJSON(departure, language, verbose))
	}

	return response
}

func departureToJSON(departure models.Departure, language string, verbose bool) map[string]interface{} {
	response := map[string]interface{}{
		"service_id":               departure.ServiceID,
		"name":                     nullString(departure.ServiceName),
		"timestamp":                departure.Timestamp,
		"status":                   departure.Status,
		"service_date":             departure.ServiceDate,
		"service_number":           departure.ServiceNumber,
		"station":                  departure.Station.Code,
		"type":                     departure.ServiceType,
		"type_code":                departure.ServiceTypeCode,
		"company":                  departure.Company,
		"destination_actual":       nullString(departure.ActualDestinationString()),
		"destination_planned":      nullString(departure.PlannedDestinationString()),
		"destination_actual_codes": departure.ActualDestinationCodes(),
		"via":                      nullString(departure.ViaStationsString()),
		"departure_time":           localTimeString(departure.DepartureTime),
		"platform_actual":          nullString(departure.PlatformActual),
		"platform_planned":         nullString(departure.PlatformPlanned),
		"delay":                    departure.Delay,
		"cancelled":                departure.Cancelled,
		"platform_changed":         departure.PlatformChanged(),

		"remarks": []interface{}{},
		"tips":    []interface{}{},

		"wings": []interface{}{},
	}

	response["remarks"], response["tips"] = departure.GetRemarksTips(language)

	responseWings := []interface{}{}

	if verbose {
		for _, trainWing := range departure.TrainWings {
			wingResponse := map[string]interface{}{
				"destination_actual":  trainWing.DestinationActual,
				"destination_planned": trainWing.DestinationPlanned,
				"remarks":             models.GetRemarks(trainWing.Modifications, language),
				"stations":            []interface{}{},
			}

			wingResponse["stops"] = trainWing.Stations
			wingResponse["material"] = materialsToJSON(trainWing.Material, language, verbose)

			responseWings = append(responseWings, wingResponse)
		}
	}

	response["wings"] = responseWings

	return response
}
