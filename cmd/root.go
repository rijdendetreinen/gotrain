package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/evalphobia/logrus_sentry"
	log "github.com/sirupsen/logrus"
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
		log.WithField("file", viper.ConfigFileUsed()).Info("Using config file:")
	}

	log.Debug("Configuration loaded")
}

// Initialize logger
func initLogger(cmd *cobra.Command) {
	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Verbose logging enabled")
	}

	if viper.GetString("sentry.dsn") != "" {
		// Set log levels. Logging warnings is optional
		logLevels := []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		}

		if viper.GetBool("sentry.warnings") {
			logLevels = append(logLevels, log.WarnLevel)
		}

		hook, err := logrus_sentry.NewSentryHook(viper.GetString("sentry.dsn"), logLevels)

		// 5s timeout seems reasonable
		hook.Timeout = 5 * time.Second

		// Set release version:
		hook.SetRelease(Version.Version)

		// Set environment:
		hook.SetEnvironment(viper.GetString("sentry.environment"))

		// We want these errors in our log but not in Sentry
		hook.SetIgnoreErrors("Shutting down", "Received signal: interrupt, shutting down", "Exiting")

		if err == nil {
			log.AddHook(hook)
			log.WithField("dsn", viper.GetString("sentry.dsn")).Debug("Sentry logging enabled")
		} else {
			log.WithError(err).Error("Failed to initialize Sentry logging")
		}
	}
}
