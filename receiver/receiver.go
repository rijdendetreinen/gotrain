package receiver

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"time"

	"github.com/rijdendetreinen/gotrain/archiver"

	"github.com/pebbe/zmq4"
	"github.com/rijdendetreinen/gotrain/parsers"
	"github.com/rijdendetreinen/gotrain/stores"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ArchiveServices bool
var ProcessStores bool

// ReceiveData connects to the ZMQ server and starts receiving data
func ReceiveData(exit chan bool) {
	subscriber, _ := zmq4.NewSocket(zmq4.SUB)

	defer subscriber.Close()

	subscriber.SetLinger(0)
	subscriber.SetRcvtimeo(1 * time.Second)

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	zmqHost := viper.GetString("source.server")
	envelopes := viper.GetStringMapString("source.envelopes")

	subscriber.Connect(zmqHost)
	log.WithField("host", zmqHost).Info("Connect to server")

	// Subscribe to all envelopes:
	if !ProcessStores && ArchiveServices {
		log.Info("Archiver enabled, not processing departures and arrivals. Only subscribing to services")
	}

	for key, envelope := range envelopes {
		if !ProcessStores && ArchiveServices {
			if key != "services" {
				continue
			}
		}
		log.WithFields(log.Fields{
			"system":   key,
			"envelope": envelope,
		}).Info("Subscribed to envelope")
		subscriber.SetSubscribe(envelope)
	}

	listen(subscriber, envelopes, exit)
}

// Listen for messages
func listen(subscriber *zmq4.Socket, envelopes map[string]string, exit chan bool) {
	log.Info("Receiving data...")

	for {
		select {
		case <-exit:
			log.Info("Shutting down receiver")

			subscriber.Close()
			log.Info("Receiver shut down")

			exit <- true

			return
		default:
			msg, err := subscriber.RecvMessageBytes(0)

			if err != nil {
				continue
			}

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
					departure, err := parsers.ParseDvsMessage(message)

					if err != nil {
						log.WithError(err).Error("Could not parse departure message")
						stores.Stores.DepartureStore.Counters.Error++
					} else {
						if ProcessStores {
							stores.Stores.DepartureStore.ProcessDeparture(departure)
						}

						log.WithFields(log.Fields{
							"ProductID":   departure.ProductID,
							"DepartureID": departure.ID,
						}).Debug("Departure received")
					}

				case strings.HasPrefix(envelope, envelopes["arrivals"]):
					arrival, err := parsers.ParseDasMessage(message)

					if err != nil {
						log.WithError(err).Error("Could not parse arrival message")
						stores.Stores.ArrivalStore.Counters.Error++
					} else {
						if ProcessStores {
							stores.Stores.ArrivalStore.ProcessArrival(arrival)
						}

						log.WithFields(log.Fields{
							"ProductID": arrival.ProductID,
							"ArrivalID": arrival.ID,
						}).Debug("Arrival received")
					}

				case strings.HasPrefix(envelope, envelopes["services"]):
					service, err := parsers.ParseRitMessage(message)

					if err != nil {
						log.WithError(err).Error("Could not parse service message")
						stores.Stores.ServiceStore.Counters.Error++
					} else {
						if ProcessStores {
							stores.Stores.ServiceStore.ProcessService(service)
						}
						if ArchiveServices {
							archiver.ProcessService(service)
						}

						log.WithFields(log.Fields{
							"ProductID": service.ProductID,
							"ServiceID": service.ID,
						}).Debug("Service received")
					}

				default:
					log.WithFields(log.Fields{
						"envelope": envelope,
					}).Warning("Unknown envelope")
				}
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
