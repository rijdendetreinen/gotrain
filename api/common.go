package api

import (
	"net/url"
	"time"

	"github.com/rijdendetreinen/gotrain/stores"
)

// Statistics includes counters and the inventory
type Statistics struct {
	Counters  stores.Counters `json:"counters"`
	Inventory int             `json:"inventory"`
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
