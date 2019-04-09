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

	remarks := models.GetRemarks(departure.Modifications, language)
	var tips []string

	if !departure.Cancelled {
		if departure.DoNotBoard {
			remarks = append(remarks, models.RemarkTranslation("Niet instappen", "Do not board", language))
		}
		if departure.RearPartRemains {
			remarks = append(remarks, models.RemarkTranslation("Achterste treindeel blijft achter", "Rear train part: do not board", language))
		}
		if departure.ReservationRequired {
			tips = append(tips, models.RemarkTranslation("Reservering verplicht", "Reservation required", language))
		}
		if departure.WithSupplement {
			tips = append(tips, models.RemarkTranslation("Toeslag verplicht", "Supplement required", language))
		}
		if departure.SpecialTicket {
			tips = append(tips, models.RemarkTranslation("Bijzonder ticket", "Special ticket", language))
		}

		// TODO: boardingtips etc.
		// TODO: check for material which does not continue to terminal station
	}

	if departure.ServiceName != "" {
		tips = append(tips, departure.ServiceName)
	}

	response["remarks"] = remarks
	response["tips"] = tips

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
