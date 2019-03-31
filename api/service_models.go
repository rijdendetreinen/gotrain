package api

import (
	"github.com/rijdendetreinen/gotrain/models"
)

func serviceToJSON(service models.Service, language string, verbose bool) map[string]interface{} {
	response := map[string]interface{}{
		"id":             service.ID,
		"timestamp":      service.Timestamp,
		"service_date":   service.ServiceDate,
		"service_number": service.ServiceNumber,
		"type":           service.ServiceType,
		"type_code":      service.ServiceTypeCode,
		"company":        service.Company,

		"journey_planner":      service.JourneyPlanner,
		"reservation_required": service.ReservationRequired,
		"special_ticket":       service.SpecialTicket,
		"with_supplement":      service.WithSupplement,

		"parts":   []interface{}{},
		"remarks": []interface{}{},
		"tips":    []interface{}{},
	}

	responseParts := []interface{}{}

	for _, part := range service.ServiceParts {
		partResponse := map[string]interface{}{
			"service_number": part.ServiceNumber,
			"remarks":        []interface{}{},
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
		"station_accesible":    stop.StationAccesible,
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

		"remarks": []interface{}{},
		"tips":    []interface{}{},
		// "material": stop.Material,
	}

	return stopResponse
}
