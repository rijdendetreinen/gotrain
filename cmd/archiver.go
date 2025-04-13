package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rijdendetreinen/gotrain/archiver"
	"github.com/rijdendetreinen/gotrain/receiver"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var archiverCommand = &cobra.Command{
	Use:   "archiver",
	Short: "Start archiver",
	Long:  `Start the GoTrain archiver. It receives data and pushes processed data to the archive queue.`,
	Run: func(cmd *cobra.Command, args []string) {
		startArchiver(cmd)
	},
}

func init() {
	RootCmd.AddCommand(archiverCommand)
}

var exitArchiverReceiverChannel = make(chan bool)

func startArchiver(cmd *cobra.Command) {
	initLogger(cmd)

	log.Info().Msgf("GoTrain archiver %v starting", Version.VersionStringLong())

	signalChan := make(chan os.Signal, 1)
	shutdownArchiverFinished := make(chan struct{})

	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		log.Warn().Msgf("Received signal: %+v, shutting down", sig)
		signal.Reset()
		shutdownArchiver()
		close(shutdownArchiverFinished)
	}()

	connectionError := archiver.Connect()

	if connectionError != nil {
		log.Error().Err(connectionError).Msg("Error while connecting to archive queue")
		return
	}

	receiver.ProcessStores = false
	receiver.ArchiveServices = true

	go receiver.ReceiveData(exitArchiverReceiverChannel)

	<-shutdownArchiverFinished
	log.Warn().Msg("Exiting")
}

func shutdownArchiver() {
	log.Warn().Msg("Shutting down")

	exitArchiverReceiverChannel <- true

	<-exitArchiverReceiverChannel
}
