package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

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
		log.WithField("file", viper.ConfigFileUsed()).Debug("Using config file:")
	}

	log.Debug("Configuration loaded")
}
