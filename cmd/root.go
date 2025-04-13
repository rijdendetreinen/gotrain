package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryzerolog "github.com/getsentry/sentry-go/zerolog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

// Version contains the version information
var Version VersionInformation

// VersionInformation is a simple struct containing the version information
type VersionInformation struct {
	Version string
	Commit  string
	Date    string
}

// VersionStringLong returns a version string
func (v VersionInformation) VersionStringLong() string {
	return fmt.Sprintf("%v (%v; built %v)", v.Version, v.Commit, v.Date)
}

// VersionStringShort returns a shortened version string
func (v VersionInformation) VersionStringShort() string {
	return fmt.Sprintf("%v (%v)", v.Version, v.Commit)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gotrain",
	Short: "GoTrain tracks trains",
	Long:  "GoTrain is a server to process real-time information about Dutch trains",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	preInitLogger(RootCmd)
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")

		viper.AddConfigPath("/etc/gotrain/")
		viper.AddConfigPath("$HOME/.gotrain")
		viper.AddConfigPath("./config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Str("file", viper.ConfigFileUsed()).Msg("Using config file")
	}

	log.Info().Msg("Configuration loaded")
}

// Pre-initialize function to set up the logger
// This function is called before the command is executed
// It sets up the logger based on the command line flags and configuration
// It also sets the global log level based on the verbose flag
func preInitLogger(cmd *cobra.Command) {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05"}
	log.Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
}

// Initialize logger
func initLogger(cmd *cobra.Command) {
	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("Verbose logging enabled")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if viper.GetString("sentry.dsn") != "" {
		initSentry()
	}
}

func initSentry() {
	// Initialize Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         viper.GetString("sentry.dsn"),
		Environment: viper.GetString("sentry.environment"),
		Release:     Version.Version,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Modify or filter events before sending them to Sentry
			return event
		},
		Debug:            true,
		AttachStacktrace: true,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("sentry initialization failed")
	}

	defer sentry.Flush(2 * time.Second)

	// Configure Zerolog to use Sentry as a writer
	sentryWriter, err := sentryzerolog.New(sentryzerolog.Config{
		ClientOptions: sentry.ClientOptions{
			Dsn:         viper.GetString("sentry.dsn"),
			Environment: viper.GetString("sentry.environment"),
			Release:     Version.Version,
		},
		Options: sentryzerolog.Options{
			Levels:          []zerolog.Level{zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel},
			WithBreadcrumbs: true,
			FlushTimeout:    3 * time.Second,
		},
	})

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create sentry writer")
	}

	defer sentryWriter.Close()

	// Use Sentry writer in Zerolog
	log.Logger = log.Output(zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr}, sentryWriter))
}
