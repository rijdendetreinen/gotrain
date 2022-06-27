package prometheus_interface

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rijdendetreinen/gotrain/stores"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Init counters enzo.
func SetupPrometheus() {
	registerStoreMetrics()
}

func StartPrometheusInterface() {
	address := viper.GetString("prometheus.address")

	if address == "" {
		address = ":2112"
	}

	log.WithField("address", address).Info("Prometheus interface started")
	http.Handle("/metrics", promhttp.Handler())

	go http.ListenAndServe(address, nil)
}

func registerStoreMetrics() {
	// Departures
	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "departures",
			Name:      "received",
			Help:      "Number of received messages",
		},
		func() float64 { return float64(stores.Stores.DepartureStore.Counters.Received) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "departures",
			Name:      "duplicates",
			Help:      "Number of detected duplicates",
		},
		func() float64 { return float64(stores.Stores.DepartureStore.Counters.Duplicates) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "departures",
			Name:      "error",
			Help:      "Number of messages with an error",
		},
		func() float64 { return float64(stores.Stores.DepartureStore.Counters.Error) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "departures",
			Name:      "processed",
			Help:      "Number of processed messages",
		},
		func() float64 { return float64(stores.Stores.DepartureStore.Counters.Processed) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "departures",
			Name:      "outdated",
			Help:      "Number of outdated messages",
		},
		func() float64 { return float64(stores.Stores.DepartureStore.Counters.Outdated) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "departures",
			Name:      "late",
			Help:      "Number of too late messages",
		},
		func() float64 { return float64(stores.Stores.DepartureStore.Counters.TooLate) },
	))

	prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "gotrain",
			Subsystem: "departures",
			Name:      "inventory",
			Help:      "Number of departures in memory",
		},
		func() float64 { return float64(stores.Stores.DepartureStore.GetNumberOfDepartures()) },
	))

	// Arrivals
	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "arrivals",
			Name:      "received",
			Help:      "Number of received messages",
		},
		func() float64 { return float64(stores.Stores.ArrivalStore.Counters.Received) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "arrivals",
			Name:      "duplicates",
			Help:      "Number of detected duplicates",
		},
		func() float64 { return float64(stores.Stores.ArrivalStore.Counters.Duplicates) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "arrivals",
			Name:      "error",
			Help:      "Number of messages with an error",
		},
		func() float64 { return float64(stores.Stores.ArrivalStore.Counters.Error) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "arrivals",
			Name:      "processed",
			Help:      "Number of processed messages",
		},
		func() float64 { return float64(stores.Stores.ArrivalStore.Counters.Processed) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "arrivals",
			Name:      "outdated",
			Help:      "Number of outdated messages",
		},
		func() float64 { return float64(stores.Stores.ArrivalStore.Counters.Outdated) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "arrivals",
			Name:      "late",
			Help:      "Number of too late messages",
		},
		func() float64 { return float64(stores.Stores.ArrivalStore.Counters.TooLate) },
	))

	prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "gotrain",
			Subsystem: "arrivals",
			Name:      "inventory",
			Help:      "Number of arrivals in memory",
		},
		func() float64 { return float64(stores.Stores.ArrivalStore.GetNumberOfArrivals()) },
	))

	// Services
	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "services",
			Name:      "received",
			Help:      "Number of received messages",
		},
		func() float64 { return float64(stores.Stores.ServiceStore.Counters.Received) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "services",
			Name:      "duplicates",
			Help:      "Number of detected duplicates",
		},
		func() float64 { return float64(stores.Stores.ServiceStore.Counters.Duplicates) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "services",
			Name:      "error",
			Help:      "Number of messages with an error",
		},
		func() float64 { return float64(stores.Stores.ServiceStore.Counters.Error) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "services",
			Name:      "processed",
			Help:      "Number of processed messages",
		},
		func() float64 { return float64(stores.Stores.ServiceStore.Counters.Processed) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "services",
			Name:      "outdated",
			Help:      "Number of outdated messages",
		},
		func() float64 { return float64(stores.Stores.ServiceStore.Counters.Outdated) },
	))

	prometheus.Register(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: "gotrain",
			Subsystem: "services",
			Name:      "late",
			Help:      "Number of too late messages",
		},
		func() float64 { return float64(stores.Stores.ServiceStore.Counters.TooLate) },
	))

	prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "gotrain",
			Subsystem: "services",
			Name:      "inventory",
			Help:      "Number of services in memory",
		},
		func() float64 { return float64(stores.Stores.ServiceStore.GetNumberOfServices()) },
	))
}
