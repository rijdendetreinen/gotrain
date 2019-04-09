package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rijdendetreinen/gotrain/models"
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

func departuresStation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	station := vars["station"]
	language := getLanguageVar(r.URL)
	verbose := getBooleanQueryParameter(r.URL, "verbose", false)

	departures := stores.Stores.DepartureStore.GetStationDepartures(station)

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
		"platform_actual":          departure.PlatformActual,
		"platform_planned":         departure.PlatformPlanned,
		"delay":                    departure.Delay,
		"cancelled":                departure.Cancelled,
		"platform_changed":         departure.PlatformChanged(),

		"do_not_board":         departure.DoNotBoard,
		"reservation_required": departure.ReservationRequired,
		"special_ticket":       departure.SpecialTicket,
		"with_supplement":      departure.WithSupplement,
		"rear_part_remains":    departure.RearPartRemains,
		"not_real_time":        departure.NotRealTime,

		"remarks": models.GetRemarks(departure.Modifications, language),
		"tips":    []interface{}{},

		"wings": []interface{}{},
	}

	if departure.ServiceName != "" {
		response["tips"] = departure.ServiceName
	}

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
