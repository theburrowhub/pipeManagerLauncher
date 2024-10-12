package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sergiotejon/pipeManager/internal/app/launcher/buckets"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

var (
	artifactCommit   string
	artifactsProject string
	artifactsPaths   []string
)

// artifactsCmd represents the bucket command for artifacts
var artifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "Manage artifacts",
}

// validateArtifactDownloadFlags checks if the provided flags are valid
func validateArtifactDownloadFlags() error {
	if artifactCommit == "" {
		return errors.New("commit is required")
	}
	if artifactsProject == "" {
		return errors.New("project is required")
	}
	return nil
}

// artifactDownloadCmd represents the download subcommand
var artifactDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Artifact download from bucket",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup the application
		setup()

		if err := validateArtifactDownloadFlags(); err != nil {
			logging.Logger.Error("Invalid flags", "error", err)
			os.Exit(1)
		}

		fileName := fmt.Sprintf("%s.tar.gz", artifactCommit)
		source := filepath.Join("artifacts", artifactsProject, fileName)

		err := buckets.Download(source)
		if err != nil {
			logging.Logger.Error("Artifact download from bucket failed", "error", err)
			os.Exit(ErrCodeBucketDownload)
		}
		logging.Logger.Info("Artifact download from bucket successful")
	},
}

// validateArtifactUploadFlags checks if the provided flags are valid
func validateArtifactUploadFlags() error {
	if artifactCommit == "" {
		return errors.New("commit is required")
	}
	if artifactsProject == "" {
		return errors.New("project is required")
	}
	if len(artifactsPaths) == 0 {
		return errors.New("paths for upload are required")
	}
	return nil
}

// artifactUploadCmd represents the upload subcommand
var artifactUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Artifact upload to bucket",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup the application
		setup()

		if err := validateArtifactUploadFlags(); err != nil {
			logging.Logger.Error("Invalid flags", "error", err)
			os.Exit(1)
		}

		fileName := fmt.Sprintf("%s.tar.gz", artifactCommit)
		destination := filepath.Join("artifacts", artifactsProject, fileName)

		err := buckets.Upload(artifactsPaths, destination)
		if err != nil {
			logging.Logger.Error("Artifacts upload to bucket failed", "error", err)
			os.Exit(ErrCodeBucketUpload)
		}
		logging.Logger.Info("Artifacts upload to bucket successful")
	},
}

func init() {
	var err error

	artifactsCmd.PersistentFlags().StringVar(&artifactCommit, "commit", "", "Commit hash in case of artifact")
	artifactsCmd.PersistentFlags().StringVar(&artifactsProject, "project", "", "Project name")

	err = artifactsCmd.MarkPersistentFlagRequired("commit")
	if err != nil {
		logging.Logger.Error("Error marking flag as required", "error", err)
		return
	}
	err = artifactsCmd.MarkPersistentFlagRequired("project")
	if err != nil {
		logging.Logger.Error("Error marking flag as required", "error", err)
		return
	}

	artifactUploadCmd.PersistentFlags().StringSliceVar(&artifactsPaths, "path", []string{}, "List of directories and files")

	err = artifactUploadCmd.MarkPersistentFlagRequired("path")
	if err != nil {
		logging.Logger.Error("Error marking flag as required", "error", err)
		return
	}

	artifactsCmd.AddCommand(artifactDownloadCmd)
	artifactsCmd.AddCommand(artifactUploadCmd)
}
