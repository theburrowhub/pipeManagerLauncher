package cmd

import (
	"errors"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	Run: func(cmd *cobra.Command, args []string) {

		if err := validateRunFlags(); err != nil {
			slog.Error("Invalid flags", "error", err)
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
