// Package main contains the main entrypoint for the application.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/sergiotejon/pipeManager/internal/app/webhook-listener/httpServer"
	"github.com/sergiotejon/pipeManager/internal/app/webhook-listener/pipeline"
	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
	"github.com/sergiotejon/pipeManager/internal/pkg/version"
)

var (
	defaultConfigFile = "/etc/pipe-manager/config.yaml" // defaultConfigFile is the default configuration file
	defaultListenAddr = ":80"                           // defaultListenAddr is the default listen address
	configFile        string                            // configFile is the path to the configuration file
	listenAddr        string                            // listenAddr is the address to listen on
	showVersion       bool                              // showVersion is a flag to show the version
)

// main is the entrypoint for the application
// It sets up the root command and executes the application
func main() {
	rootCmd := &cobra.Command{
		Use:   "pipe-manager",
		Short: "Pipe Manager CLI",
		Run: func(cmd *cobra.Command, args []string) {
			// Show version
			if showVersion {
				fmt.Println(version.GetVersion())
				os.Exit(0)
			}

			// Run the application
			app()
		},
	}

	rootCmd.Flags().StringVarP(&configFile, "config", "c", defaultConfigFile, "Path to the config file")
	rootCmd.Flags().StringVarP(&listenAddr, "listen", "l", defaultListenAddr, "Listen address")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print the version")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

// app is the main application function
// It loads the configuration, sets up the logger and starts the web server for processing requests of the webhook
func app() {
	var err error

	// Load configuration
	err = config.LoadCommonConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading common config: %v", err)
	}

	err = config.LoadWebhookConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading webhook config: %v", err)
	}

	err = config.LoadLauncherConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading launcher config: %v", err)
	}

	// Setup Logger
	err = logging.SetupLogger(config.Common.Data.Log.Level, config.Common.Data.Log.Format, config.Common.Data.Log.File)
	if err != nil {
		log.Fatalf("Error configuring the logger: %v", err)
	}

	logging.Logger.Info("Pipe Manager starting up...")
	logging.Logger.Info("Setup", "configFile", configFile,
		"workers", config.Webhook.Data.Workers,
		"listenAddr", listenAddr,
		"logLevel", config.Common.Data.Log.Level,
		"logFormat", config.Common.Data.Log.Format,
		"logFile", config.Common.Data.Log.File,
		"launcherImage", pipeline.GetLauncherImage())

	// Launch web server
	err = httpServer.HttpServer(listenAddr)
	if err != nil {
		logging.Logger.Error("Error starting server", "error", fmt.Sprintf("%v", err))
		panic(err)
	}

	return
}
