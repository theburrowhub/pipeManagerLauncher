package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/sergiotejon/pipeManager/internal/app/launcher/repository"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
	"github.com/sergiotejon/pipeManager/pkg/config"
)

var (
	repoURL     string
	cloneDepth  int
	cloneCommit string
	destination string
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone repository",
	Run: func(cmd *cobra.Command, args []string) {
		// Set up the application
		setup()

		// Run the clone command
		err := repository.Clone(repoURL, cloneDepth, cloneCommit, destination)
		if err != nil {
			logging.Logger.Error("Error cloning repository", "error", err)
			os.Exit(ErrCodeCloneRepo)
		}
		logging.Logger.Info("Repository cloned successfully", "destination", destination, "commit", cloneCommit, "depth", cloneDepth, "repository", repoURL)
	},
}

func init() {
	var err error

	cloneCmd.Flags().StringVar(&repoURL, "repository", "", "Repository URL")
	cloneCmd.Flags().IntVar(&cloneDepth, "depth", config.Launcher.Data.CloneDepth, "Depth of the clone")
	cloneCmd.Flags().StringVar(&cloneCommit, "commit", "", "Commit to checkout")
	cloneCmd.Flags().StringVar(&destination, "destination", repoDir, "Destination directory")

	err = cloneCmd.MarkFlagRequired("repository")
	if err != nil {
		logging.Logger.Error("Error marking flag as required", "error", err)
		return
	}
	err = cloneCmd.MarkFlagRequired("commit")
	if err != nil {
		logging.Logger.Error("Error marking flag as required", "error", err)
		return
	}
}
