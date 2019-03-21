package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rijdendetreinen/gotrain/api"
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

func startServer(cmd *cobra.Command) {
	initLogger(cmd)

	log.Info("GoTrain starting")

	signalChan := make(chan os.Signal, 1)
	shutdownFinished := make(chan struct{})

	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)

	initStores()

	go func() {
		sig := <-signalChan
		log.Errorf("Received signal: %+v, shutting down", sig)
		signal.Reset()
		shutdown()
		close(shutdownFinished)
	}()

	go receiver.ReceiveData(exitReceiverChannel)

	apiAddress := viper.GetString("api.address")
	go api.ServeAPI(apiAddress, exitRestAPI)

	<-shutdownFinished
	log.Error("Exiting")
}

func initLogger(cmd *cobra.Command) {
	// TODO: setup logger

	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Verbose logging enabled")
	}
}

func initStores() {
	stores.InitializeStores()

	log.Info("Reading saved store contents...")
	err := stores.LoadStores()

	if err != nil {
		log.WithError(err).Warn("Error while loading stores")
	}
}

func shutdown() {
	log.Warn("Shutting down")

	exitRestAPI <- true
	exitReceiverChannel <- true

	<-exitRestAPI
	<-exitReceiverChannel

	log.Info("Saving store contents...")
	err := stores.SaveStores()

	if err != nil {
		log.WithError(err).Error("Error while saving stores")
	}
}
