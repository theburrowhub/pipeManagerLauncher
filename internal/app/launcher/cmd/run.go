package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	Run: func(cmd *cobra.Command, args []string) {

		if err := validateCacheUploadFlags(); err != nil {
			logging.Logger.Error("Invalid flags", "error", err)
			os.Exit(1)
		}

		// Set up the application
		setup()
		// Run the main application
		app()
	},
}

// validateRunFlags checks if the provided flags are valid
func validateRunFlags() error {
	if configFile == "" {
		return errors.New("config file is required")
	}
	return nil
}
