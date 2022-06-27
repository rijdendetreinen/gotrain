package cmd

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rijdendetreinen/gotrain/api"
	"github.com/rijdendetreinen/gotrain/prometheus_interface"
	"github.com/rijdendetreinen/gotrain/receiver"
	"github.com/rijdendetreinen/gotrain/stores"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  `Start the GoTrain server, which starts receiving new data and starts the REST API.`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer(cmd)
	},
}

func init() {
	RootCmd.AddCommand(serverCommand)
}

var exitReceiverChannel = make(chan bool)
var exitRestAPI = make(chan bool)
var cleanupTicker *time.Ticker
var downtimeDetectorTicker *time.Ticker
var autoSaveTicker *time.Ticker

func startServer(cmd *cobra.Command) {
	initLogger(cmd)

	log.Infof("GoTrain %v starting", Version.VersionStringLong())

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if viper.GetBool("prometheus.enabled") {
		prometheus_interface.SetupPrometheus()
		prometheus_interface.StartPrometheusInterface()
	}

	signalChan := make(chan os.Signal, 1)
	shutdownFinished := make(chan struct{})

	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)

	initStores()

	go func() {
		sig := <-signalChan
		log.Warnf("Received signal: %+v, shutting down", sig)
		signal.Reset()
		shutdown()
		close(shutdownFinished)
	}()

	receiver.ProcessStores = true
	receiver.ArchiveServices = false

	go receiver.ReceiveData(exitReceiverChannel)

	apiAddress := viper.GetString("api.address")
	go api.ServeAPI(apiAddress, exitRestAPI)

	setupCleanupScheduler()
	setupDowntimeDetector()
	setupAutoSave()

	<-shutdownFinished
	log.Warn("Exiting")
}

func setupCleanupScheduler() {
	// Set up our internal "garbage collector" (which cleans up stores):
	cleanupTicker := time.NewTicker(1 * time.Minute)
	log.Debug("Cleanup scheduler set up")

	go func() {
		for {
			<-cleanupTicker.C
			stores.CleanUp()
		}
	}()
}

func setupDowntimeDetector() {
	// Set up the downtime detector, which measures approximately every 20s the number of messages received
	// for each store
	downtimeDetectorTicker := time.NewTicker(20 * time.Second)
	log.Debug("Downtime detector set up")

	go func() {
		for {
			<-downtimeDetectorTicker.C
			stores.TakeMeasurements()
		}
	}()
}

func setupAutoSave() {
	autoSaveTicker := time.NewTicker(12 * time.Hour)
	log.Debug("Autosave set up")

	go func() {
		for {
			<-autoSaveTicker.C
			log.Info("Auto-saving stores")
			log.Infof("Current inventory: %d arrivals, %d departures, %d services",
				stores.Stores.ArrivalStore.GetNumberOfArrivals(),
				stores.Stores.DepartureStore.GetNumberOfDepartures(),
				stores.Stores.ServiceStore.GetNumberOfServices())
			stores.SaveStores()
		}
	}()
}

func initStores() {
	stores.InitializeStores()

	if viper.IsSet("stores.location") {
		stores.StoresDataDirectory = viper.GetString("stores.location")
	}

	if _, err := os.Stat(stores.StoresDataDirectory); os.IsNotExist(err) {
		log.WithField("directory", stores.StoresDataDirectory).Error("Data directory does not exist; not loading stores")
	} else {
		log.WithField("directory", stores.StoresDataDirectory).Info("Data directory initialized")

		log.Info("Reading saved store contents...")
		stores.LoadStores()
	}
}

func shutdown() {
	log.Warn("Shutting down")

	if cleanupTicker != nil {
		cleanupTicker.Stop()
	}

	if downtimeDetectorTicker != nil {
		downtimeDetectorTicker.Stop()
	}

	if autoSaveTicker != nil {
		autoSaveTicker.Stop()
	}

	exitRestAPI <- true
	exitReceiverChannel <- true

	<-exitRestAPI
	<-exitReceiverChannel

	log.Info("Saving store contents...")
	stores.SaveStores()
}
