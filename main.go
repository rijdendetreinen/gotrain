package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"time"

	"github.com/rijdendetreinen/gotrain/api"
	"github.com/rijdendetreinen/gotrain/stores"

	"github.com/rijdendetreinen/gotrain/models"
	"github.com/spf13/viper"

	"github.com/rijdendetreinen/gotrain/parsers"

	"fmt"

	"github.com/pebbe/zmq4"
)

var exit = make(chan bool)

func main() {
	fmt.Println("Gotrain starting")

	loadConfig()

	initStores()

	go receiveData()

	apiAddress := viper.GetString("api.address")
	fmt.Printf("REST API listening on %s\n", apiAddress)
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
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

func initStores() {
	stores.InitializeStores()
}

func receiveData() {
	subscriber, _ := zmq4.NewSocket(zmq4.SUB)

	subscriber.SetLinger(0)
	defer subscriber.Close()

	zmqHost := viper.GetString("source.server")
	envelopes := viper.GetStringMapString("source.envelopes")

	subscriber.Connect(zmqHost)
	fmt.Printf("Connected to %s\n", zmqHost)

	// Subscribe to all envelopes:
	for key, envelope := range envelopes {
		fmt.Printf("Subscribed to %s [%s]\n", envelope, key)
		subscriber.SetSubscribe(envelope)
	}

	fmt.Println("Receiving data...")
	Listen(subscriber, envelopes)
}

// Listen for messages
func Listen(subscriber *zmq4.Socket, envelopes map[string]string) {
	counter := 0
	ritCounter := 0

	var services []models.Service

	for {
		msg, err := subscriber.RecvMessageBytes(0)

		envelope := string(msg[0])

		if strings.HasPrefix(envelope, envelopes["services"]) {
			message, _ := gunzip(msg[1])

			if err != nil {
				fmt.Println("ERROR!", err, msg[0], string(msg[1]))
			} else {
				service := parsers.ParseRitMessage(message)
				fmt.Println(time.Now().Format(time.RFC3339), ritCounter, service.ProductID, service.ServiceDate, service.ServiceID, service.ServiceParts[0].Stops[0].Station.NameLong, "-", service.ServiceParts[0].Stops[len(service.ServiceParts[0].Stops)-1].Station.NameLong)

				services = append(services, service)

				stores.ProcessService(&stores.Stores.ServiceStore, service)

				ritCounter++
			}
		} else {
			fmt.Println(string(msg[0]))
		}

		counter++
	}
}

func gunzip(data []byte) (io.Reader, error) {
	buf := bytes.NewBuffer(data)
	reader, err := gzip.NewReader(buf)
	defer reader.Close()

	if err != nil {
		// panic(err)
		return nil, err
	}

	buf3 := new(bytes.Buffer)
	buf3.ReadFrom(reader)

	return buf3, nil
}
