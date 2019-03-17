package receiver

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"

	"github.com/pebbe/zmq4"
	"github.com/rijdendetreinen/gotrain/parsers"
	"github.com/rijdendetreinen/gotrain/stores"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ReceiveData connects to the ZMQ server and starts receiving data
func ReceiveData() {
	subscriber, _ := zmq4.NewSocket(zmq4.SUB)

	subscriber.SetLinger(0)
	defer subscriber.Close()

	zmqHost := viper.GetString("source.server")
	envelopes := viper.GetStringMapString("source.envelopes")

	subscriber.Connect(zmqHost)
	log.WithField("host", zmqHost).Info("Connect to server")

	// Subscribe to all envelopes:
	for key, envelope := range envelopes {
		log.WithFields(log.Fields{
			"system":   key,
			"envelope": envelope,
		}).Info("Subscribed to envelope")
		subscriber.SetSubscribe(envelope)
	}

	listen(subscriber, envelopes)
}

// Listen for messages
func listen(subscriber *zmq4.Socket, envelopes map[string]string) {
	log.Info("Receiving data...")

	for {
		msg, err := subscriber.RecvMessageBytes(0)

		envelope := string(msg[0])

		// Decompress message:

		message, _ := gunzip(msg[1])

		if err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"envelope": envelope,
				"message":  string(msg[1]),
			}).Error("Error decompressing message. Message ignored")
		} else {
			switch {
			case strings.HasPrefix(envelope, envelopes["departures"]) == true:
				departure := parsers.ParseDvsMessage(message)
				stores.Stores.DepartureStore.ProcessDeparture(departure)

				log.WithFields(log.Fields{
					"ProductID":   departure.ProductID,
					"DepartureID": departure.ID,
				}).Debug("Departure received")

			case strings.HasPrefix(envelope, envelopes["arrivals"]):
				// TODO: process arrival
				log.Debug("Arrival received")

			case strings.HasPrefix(envelope, envelopes["services"]):
				service := parsers.ParseRitMessage(message)
				stores.Stores.ServiceStore.ProcessService(service)

				log.WithFields(log.Fields{
					"ProductID": service.ProductID,
					"ServiceID": service.ID,
				}).Debug("Service received")

			default:
				log.WithFields(log.Fields{
					"envelope": envelope,
				}).Warning("Unknown envelope")
			}
		}
	}
}

func gunzip(data []byte) (io.Reader, error) {
	buf := bytes.NewBuffer(data)
	reader, err := gzip.NewReader(buf)
	defer reader.Close()

	if err != nil {
		return nil, err
	}

	buf3 := new(bytes.Buffer)
	buf3.ReadFrom(reader)

	return buf3, nil
}
