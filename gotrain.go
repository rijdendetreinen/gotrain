package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rijdendetreinen/gotrain/api"
	"github.com/rijdendetreinen/gotrain/receiver"
	"github.com/rijdendetreinen/gotrain/stores"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var exitReceiverChannel = make(chan bool)
var exitRestAPI = make(chan bool)

func main() {
	log.Info("Gotrain starting")

	// TODO: setup logger:
	log.SetLevel(log.DebugLevel)

	signalChan := make(chan os.Signal, 1)
	shutdownFinished := make(chan struct{})

	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)

	loadConfig()

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

func shutdown() {
	log.Warn("Shutting down")

	exitRestAPI <- true
	exitReceiverChannel <- true

	<-exitRestAPI
	<-exitReceiverChannel
}

func loadConfig() {
	viper.SetConfigName("config")

	viper.AddConfigPath("/etc/gotrain/")
	viper.AddConfigPath("$HOME/.gotrain")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()

	if err != nil {
		log.WithField("error", err).Panic("Could not load config file")
	}

	log.Debug("Configuration loaded")
}

func initStores() {
	stores.InitializeStores()
}
