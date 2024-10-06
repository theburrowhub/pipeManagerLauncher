package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/sergiotejon/pipeManager/internal/app/launcher/repository"
	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/envvars"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone repository",
	Run: func(cmd *cobra.Command, args []string) {
		initApp()

		err := repository.Clone(envvars.Variables["REPOSITORY"],
			config.Launcher.Data.CloneDepth,
			envvars.Variables["COMMIT"],
			repoDir)
		if err != nil {
			logging.Logger.Error("Error cloning repository", "error", err)
			os.Exit(ErrCodeClone)
		}

	},
}
