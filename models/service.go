package models

import "time"

// Service is the train service containing all parts, stops etc.
type Service struct {
	StoreItem

	ValidUntil      time.Time `json:"-"`
	ServiceID       string    `json:"service_id"`
	ServiceNumber   string    `json:"service_number"`
	ServiceDate     string    `json:"service_date"`
	ServiceType     string    `json:"type"`
	ServiceTypeCode string    `json:"type_code"`
	Company         string    `json:"company"`

	ServiceParts []ServicePart `json:"parts"`

	ReservationRequired bool `json:"reservation_required"`
	WithSupplement      bool `json:"with_supplement"`
	SpecialTicket       bool `json:"special_ticket"`
	JourneyPlanner      bool `json:"in_journey_planner"`

	Modifications []Modification `json:"modifications"`

	Hidden bool `json:"-"`
}

// ServicePart is a single part of a train service (a service usually contains just one part, but may contain more)
type ServicePart struct {
	ServiceNumber string         `json:"service_number"`
	Stops         []ServiceStop  `json:"stops"`
	Modifications []Modification `json:"modifications"`
}

// ServiceStop is a stops which is called by a train service.
type ServiceStop struct {
	Station Station `json:"station"`

	StationAccesible    bool `json:"station_accesible"`
	AssistenceAvailable bool `json:"assistance_available"`

	DestinationActual        string `json:"-"`
	DestinationPlanned       string `json:"-"`
	ArrivalPlatformActual    string `json:"arrival_platform_actual"`
	ArrivalPlatformPlanned   string `json:"arrival_platform_planned"`
	DeparturePlatformActual  string `json:"departure_platform_actual"`
	DeparturePlatformPlanned string `json:"departure_platform_planned"`

	StoppingActual  bool   `json:"stopping_actual"`
	StoppingPlanned bool   `json:"stopping_planned"`
	StopType        string `json:"stop_type"`
	DoNotBoard      bool   `json:"do_not_board"`

	ArrivalTime    time.Time `json:"arrival_time"`
	ArrivalDelay   int       `json:"arrival_delay"`
	DepartureTime  time.Time `json:"departure_time"`
	DepartureDelay int       `json:"departure_delay"`

	ArrivalCancelled   bool `json:"arrival_cancelled"`
	DepartureCancelled bool `json:"departure_cancelled"`

	Modifications []Modification `json:"modifications"`
	Material      []Material     `json:"material"`
}
