package api

import (
	"net/url"
	"time"

	"github.com/rijdendetreinen/gotrain/models"
	"github.com/rijdendetreinen/gotrain/stores"
)

// Statistics includes counters and the inventory
type Statistics struct {
	Counters         stores.Counters `json:"counters"`
	Inventory        int             `json:"inventory"`
	Status           string          `json:"status"`
	LastStatusChange time.Time       `json:"last_status_change"`
	MessagesAverage  float64         `json:"average_messages"`
}

func getLanguageVar(url *url.URL) string {
	language := url.Query().Get("language")

	// Only NL and EN allowed
	if language == "en" {
		return language
	}

	return "nl"
}

func getBooleanQueryParameter(url *url.URL, variable string, defaultValue bool) bool {
	value := url.Query().Get(variable)

	if value != "" {
		return value == "true"
	}

	return defaultValue
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

func materialToJSON(material models.Material, language string, verbose bool) map[string]interface{} {
	materialResponse := map[string]interface{}{
		"type":             material.NaterialType,
		"accessible":       material.Accessible,
		"number":           material.NormalizedNumber(),
		"position":         material.Position,
		"remains_behind":   material.RemainsBehind,
		"closed":           material.Closed,
		"added":            material.Added,
		"destination":      material.DestinationActual.NameLong,
		"destination_code": material.DestinationActual.Code,
	}

	return materialResponse
}

func materialsToJSON(materials []models.Material, language string, verbose bool) []map[string]interface{} {
	materialsResponse := []map[string]interface{}{}

	for _, material := range materials {
		materialsResponse = append(materialsResponse, materialToJSON(material, language, verbose))
	}

	return materialsResponse
}
