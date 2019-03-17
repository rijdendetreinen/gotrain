package main

import (
	"github.com/rijdendetreinen/gotrain/api"
	"github.com/rijdendetreinen/gotrain/receiver"
	"github.com/rijdendetreinen/gotrain/stores"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var exit = make(chan bool)

func main() {
	log.Info("Gotrain starting")
	log.SetLevel(log.DebugLevel)

	loadConfig()

	initStores()

	go receiver.ReceiveData()

	apiAddress := viper.GetString("api.address")
	go api.ServeAPI(apiAddress)

	<-exit
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
