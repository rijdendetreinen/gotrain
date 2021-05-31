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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wrapDeparturesStatus("departures", departuresToJSON(departures, language, verbose)))
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

	var service *models.Service

	if verbose {
		// Look up service
		service = stores.Stores.ServiceStore.GetService(departure.ServiceNumber, departure.ServiceDate)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wrapDeparturesStatus("departure", departureToJSON(*departure, language, verbose, service)))
}

func departuresToJSON(departures []models.Departure, language string, verbose bool) []map[string]interface{} {
	response := make([]map[string]interface{}, 0)

	for _, departure := range departures {
		response = append(response, departureToJSON(departure, language, verbose, nil))
	}

	return response
}

func departureToJSON(departure models.Departure, language string, verbose bool, service *models.Service) map[string]interface{} {
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
		"via":                      nullString(departure.ActualViaStationsString()),
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

	if departure.Cancelled {
		// Override actual destination and via stations with planned destination and via:
		response["destination_actual"] = response["destination_planned"]
		response["via"] = nullString(departure.PlannedViaStationsString())
		response["delay"] = 0
	}

	if verbose {
		var serviceStops map[string]models.ServiceStop

		if service != nil {
			serviceStops = service.GetStops()
		}

		for _, trainWing := range departure.TrainWings {
			wingResponse := map[string]interface{}{
				"destination_actual":  trainWing.DestinationActualString(),
				"destination_planned": trainWing.DestinationPlannedString(),
				"remarks":             models.GetRemarks(trainWing.Modifications, language),
				"stops":               []interface{}{},
			}

			stops := []interface{}{}

			wingStops := trainWing.Stations

			if departure.Cancelled {
				wingStops = trainWing.StationsPlanned
			}

			for _, station := range wingStops {
				stopData := map[string]interface{}{
					"code":                       station.Code,
					"short":                      station.NameShort,
					"medium":                     station.NameMedium,
					"long":                       station.NameLong,
					"assistance_available":       false,
					"accessible":                 false,
					"arrival_time":               nil,
					"arrival_platform":           nil,
					"arrival_cancelled":          false,
					"arrival_delay":              0,
					"arrival_platform_changed":   false,
					"departure_time":             nil,
					"departure_platform":         nil,
					"departure_cancelled":        false,
					"departure_delay":            0,
					"departure_platform_changed": false,
				}

				serviceStop, exists := serviceStops[station.Code]
				if exists {
					stopData["assistance_available"] = serviceStop.AssistanceAvailable
					stopData["accessible"] = serviceStop.StationAccessible

					stopData["arrival_time"] = localTimeString(serviceStop.ArrivalTime)
					stopData["arrival_platform"] = nullString(serviceStop.ArrivalPlatformActual)
					stopData["arrival_cancelled"] = serviceStop.ArrivalCancelled
					stopData["arrival_delay"] = serviceStop.ArrivalDelay
					stopData["arrival_platform_changed"] = serviceStop.ArrivalPlatformChanged()

					stopData["departure_time"] = localTimeString(serviceStop.DepartureTime)
					stopData["departure_platform"] = nullString(serviceStop.DeparturePlatformActual)
					stopData["departure_cancelled"] = serviceStop.DepartureCancelled
					stopData["departure_delay"] = serviceStop.DepartureDelay
					stopData["departure_platform_changed"] = serviceStop.DeparturePlatformChanged()
				}

				stops = append(stops, stopData)
			}

			wingResponse["material"] = materialsToJSON(trainWing.Material, language, verbose)
			wingResponse["stops"] = stops

			responseWings = append(responseWings, wingResponse)
		}
	}

	response["wings"] = responseWings

	return response
}

func wrapDeparturesStatus(key string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status": stores.Stores.DepartureStore.Status,
		key:      data,
	}
}
