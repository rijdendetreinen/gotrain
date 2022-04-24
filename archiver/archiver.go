package archiver

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/rijdendetreinen/gotrain/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var redisDb *redis.Client

// Connect initializes the Redis client
func Connect() error {
	redisAddress := viper.GetString("archive.address")
	redisPassword := viper.GetString("archive.password")
	redisDbNumber := viper.GetInt("archive.db")

	log.WithField("address", redisAddress).
		WithField("db", redisDbNumber).
		Info("Connecting to Redis server")

	redisDb = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       redisDbNumber,
	})

	result := redisDb.Ping()

	return result.Err()
}

// ProcessService adds a service object to the queue
func ProcessService(service models.Service) {
	serviceJSON, _ := json.Marshal(serviceToJSON(service))

	if serviceJSON != nil {
		result := redisDb.LPush("services", string(serviceJSON))
		err := result.Err()

		if err != nil {
			log.WithField("error", err).WithField("ServiceID", service.ID).WithField("ProductID", service.ProductID).Error("Archiver: could not add service to queue")
		}
	}
}

func serviceToJSON(service models.Service) map[string]interface{} {
	response := map[string]interface{}{
		"id":             service.ID,
		"product":        service.ProductID,
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

		"parts":      []interface{}{},
		"remarks_nl": models.GetRemarks(service.Modifications, "nl"),
		"remarks_en": models.GetRemarks(service.Modifications, "en"),
		"tips":       []interface{}{},
	}

	responseParts := []interface{}{}

	for _, part := range service.ServiceParts {
		partResponse := map[string]interface{}{
			"service_number": part.ServiceNumber,
			"remarks_nl":     models.GetRemarks(part.Modifications, "nl"),
			"remarks_en":     models.GetRemarks(part.Modifications, "en"),
			"tips":           []interface{}{},
			"stops":          []interface{}{},
		}

		stops := part.GetStoppingStations()

		responseStops := []interface{}{}

		for _, stop := range stops {
			responseStops = append(responseStops, serviceStopToJSON(stop))
		}

		partResponse["stops"] = responseStops
		responseParts = append(responseParts, partResponse)
	}

	response["parts"] = responseParts

	return response
}

func serviceStopToJSON(stop models.ServiceStop) map[string]interface{} {
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

		"remarks_nl": models.GetRemarks(stop.Modifications, "nl"),
		"remarks_en": models.GetRemarks(stop.Modifications, "en"),
		"tips":       []interface{}{},
		"material":   materialsToJSON(stop.Material),
	}

	return stopResponse
}

func localTimeString(originalTime time.Time) *string {
	if !originalTime.IsZero() {
		formattedTime := originalTime.Local().Format(time.RFC3339)
		return &formattedTime
	}

	return nil
}

func nullString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func materialsToJSON(materials []models.Material) []map[string]interface{} {
	materialsResponse := []map[string]interface{}{}

	for _, material := range materials {
		materialsResponse = append(materialsResponse, materialToJSON(material))
	}

	return materialsResponse
}

func materialToJSON(material models.Material) map[string]interface{} {
	materialResponse := map[string]interface{}{
		"type":             material.NaterialType,
		"accessible":       material.Accessible,
		"number":           material.NormalizedNumber(),
		"position":         material.Position,
		"remains_behind":   material.RemainsBehind,
		"destination":      material.DestinationActual.NameLong,
		"destination_code": material.DestinationActual.Code,
	}

	return materialResponse
}
