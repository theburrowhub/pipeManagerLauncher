package cmd

import (
	"log"
	"os"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/envvars"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

// setup sets up the application by loading the configuration and setting up the logger.
func setup() {
	// Load configuration
	err := config.LoadLauncherConfig(configFile)
	if err != nil {
		log.Printf("Error loading launcher config: %v", err)
		os.Exit(ErrCodeLoadConfig)
	}
	err = config.LoadCommonConfig(configFile)
	if err != nil {
		log.Printf("Error loading common config: %v", err)
		os.Exit(ErrCodeLoadConfig)
	}

	// Get environment variables
	envvars.GetEnvVars()
	for key, value := range envvars.Variables {
		logging.Logger.Debug("Environment variable", key, value)
	}
}
