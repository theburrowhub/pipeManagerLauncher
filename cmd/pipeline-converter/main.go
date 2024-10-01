// Package main contains the main entrypoint for the application.
package main

import (
	"fmt"
	"github.com/sergiotejon/pipeManager/internal/app/pipe-converter/pipelineparser"
	"gopkg.in/yaml.v3"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sergiotejon/pipeManager/internal/app/pipe-converter/repository"
	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/envvars"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
	"github.com/sergiotejon/pipeManager/internal/pkg/version"
)

const (
	defaultConfigFile = "/etc/pipe-manager/config.yaml" // defaultConfigFile is the default configuration file
	templateFolder    = "/etc/pipe-manager/templates"   // templateFolder is the folder where the templates are stored
)

var (
	configFile  string // configFile is the path to the configuration file
	showVersion bool   // showVersion is a flag to show the version
)

// main is the entrypoint for the application
// It sets up the root command and executes the application
func main() {
	rootCmd := &cobra.Command{
		Use:   "pipe-manager-launcher",
		Short: "Pipe Manager Launcher CLI",
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
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print the version")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

// app is the main application function
// It loads the configuration, sets up the logger and starts the launcher
func app() {
	var err error

	// Load configuration
	err = config.LoadLauncherConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading launcher config: %v", err)
	}
	err = config.LoadCommonConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading common config: %v", err)
	}

	// Setup Logger
	err = logging.SetupLogger(config.Common.Data.Log.Level, config.Common.Data.Log.Format, config.Common.Data.Log.File)
	if err != nil {
		log.Fatalf("Error configuring the logger: %v", err)
	}

	logging.Logger.Info("Pipe Manager starting up...")
	logging.Logger.Info("Setup", "configFile", configFile,
		"logLevel", config.Common.Data.Log.Level,
		"logFormat", config.Common.Data.Log.Format,
		"logFile", config.Common.Data.Log.File)

	// Print all environment variables in log
	envvars.GetEnvVars()
	for key, value := range envvars.Variables {
		logging.Logger.Debug("Environment variable", key, value)
	}

	// Clone the repository
	const repoDir = "/tmp/repo"
	err = repository.Clone(envvars.Variables["REPOSITORY"],
		config.Launcher.Data.CloneDepth,
		envvars.Variables["COMMIT"],
		repoDir)
	if err != nil {
		slog.Error("Error cloning repository", "msg", err,
			"repository", envvars.Variables["REPOSITORY"],
			"commit", envvars.Variables["COMMIT"],
			"depth", config.Launcher.Data.CloneDepth)
		os.Exit(1)
	}

	logging.Logger.Info("Repository cloned successfully", "repository", envvars.Variables["REPOSITORY"], "commit", envvars.Variables["COMMIT"])

	// Mix all the pipeline files
	const pipelineDir = ".pipelines"
	pipelineFolder := filepath.Join(repoDir, pipelineDir)
	err, combinedData := pipelineparser.MixPipelineFiles(pipelineFolder)
	if err != nil {
		logging.Logger.Error("Error mixing pipeline files", "msg", err, "folder", pipelineFolder)
		os.Exit(1)
	}

	logging.Logger.Info("Pipeline files mixed successfully", "folder", pipelineFolder)
	for key, _ := range combinedData {
		if key == "global" {
			continue
		}
		logging.Logger.Debug("Pipeline found", "pipeline", key)
	}

	// Temporal
	if config.Common.Data.Log.Level == "debug" {
		data, err := yaml.Marshal(combinedData)
		if err != nil {
			os.Exit(1)
		}
		fmt.Println(string(data))
	}
	// Temporal

	pipelines := pipelineparser.FindPipelineByRegex(combinedData, envvars.Variables)
	for key, _ := range pipelines {
		logging.Logger.Info("Launching pipeline", "pipeline", key)
		// TODO: Launch the pipeline
	}

	// Temporal
	if config.Common.Data.Log.Level == "debug" {
		data, err := yaml.Marshal(pipelines)
		if err != nil {
			os.Exit(1)
		}
		fmt.Println(string(data))
	}
	// Temporal

	// Parse the pipeline

	return
}
