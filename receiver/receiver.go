package receiver

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"time"

	"github.com/pebbe/zmq4"
	"github.com/rijdendetreinen/gotrain/archiver"
	"github.com/rijdendetreinen/gotrain/parsers"
	"github.com/rijdendetreinen/gotrain/stores"
	"github.com/rs/zerolog/log"
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

	envelopes := map[string]string{
		"arrivals":   viper.GetString("source.envelopes.arrivals"),
		"departures": viper.GetString("source.envelopes.departures"),
		"services":   viper.GetString("source.envelopes.services"),
	}

	subscriber.Connect(zmqHost)
	log.Info().Str("host", zmqHost).Msg("Connected to server")

	// Subscribe to all envelopes:
	if !ProcessStores && ArchiveServices {
		log.Info().Msg("Archiver enabled, not processing departures and arrivals. Only subscribing to services")
	}

	for key, envelope := range envelopes {
		if !ProcessStores && ArchiveServices {
			if key != "services" {
				continue
			}
		}
		log.Info().
			Str("system", key).
			Str("envelope", envelope).
			Msg("Subscribed to envelope")
		subscriber.SetSubscribe(envelope)
	}

	listen(subscriber, envelopes, exit)
}

// Listen for messages
func listen(subscriber *zmq4.Socket, envelopes map[string]string, exit chan bool) {
	log.Info().Msg("Receiving data...")

	for {
		select {
		case <-exit:
			log.Info().Msg("Shutting down receiver")

			subscriber.Close()
			log.Info().Msg("Receiver shut down")

			exit <- true

			return
		default:
			msg, err := subscriber.RecvMessageBytes(0)

			if err != nil {
				continue
			}

			envelope := string(msg[0])

			// Decompress message:
			message, err := gunzip(msg[1])

			if err != nil {
				log.Error().
					Err(err).
					Str("envelope", envelope).
					Str("message", string(msg[1])).
					Msg("Error decompressing message. Message ignored")
			} else {
				switch {
				case strings.HasPrefix(envelope, envelopes["departures"]):
					departure, err := parsers.ParseDvsMessage(message)

					if err != nil {
						log.Error().Err(err).Msg("Could not parse departure message")
						stores.Stores.DepartureStore.Counters.Error++
					} else {
						if ProcessStores {
							stores.Stores.DepartureStore.ProcessDeparture(departure)
						}

						log.Debug().
							Str("ProductID", departure.ProductID).
							Str("DepartureID", departure.ID).
							Msg("Departure received")
					}

				case strings.HasPrefix(envelope, envelopes["arrivals"]):
					arrival, err := parsers.ParseDasMessage(message)

					if err != nil {
						log.Error().Err(err).Msg("Could not parse arrival message")
						stores.Stores.ArrivalStore.Counters.Error++
					} else {
						if ProcessStores {
							stores.Stores.ArrivalStore.ProcessArrival(arrival)
						}

						log.Debug().
							Str("ProductID", arrival.ProductID).
							Str("ArrivalID", arrival.ID).
							Msg("Arrival received")
					}

				case strings.HasPrefix(envelope, envelopes["services"]):
					service, err := parsers.ParseRitMessage(message)

					if err != nil {
						log.Error().Err(err).Msg("Could not parse service message")
						stores.Stores.ServiceStore.Counters.Error++
					} else {
						if ProcessStores {
							stores.Stores.ServiceStore.ProcessService(service)
						}
						if ArchiveServices {
							archiver.ProcessService(service)
						}

						log.Debug().
							Str("ProductID", service.ProductID).
							Str("ServiceID", service.ID).
							Msg("Service received")
					}

				default:
					log.Warn().
						Str("envelope", envelope).
						Msg("Unknown envelope")
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
