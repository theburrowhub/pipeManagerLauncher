package cmd

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sergiotejon/pipeManager/internal/app/launcher/artifacts"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

var (
	artifactCommit      string
	artifactsProject    string
	artifactsPaths      []string
	artifactDestination string
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
	if artifactDestination == "" {
		return errors.New("destination is required")
	}
	return nil
}

// artifactDownloadCmd represents the download subcommand
var artifactDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Artifact download from bucket",
	Run: func(cmd *cobra.Command, args []string) {
		// Set up the application
		setup()

		if err := validateArtifactDownloadFlags(); err != nil {
			logging.Logger.Error("Invalid flags", "error", err)
			os.Exit(1)
		}

		bucketFolder := filepath.Join("artifacts", getMD5Hash(artifactsProject), artifactCommit)

		err := artifacts.Download(artifactsPaths, bucketFolder, artifactDestination)
		if err != nil {
			logging.Logger.Error("Artifact download from bucket failed", "error", err)
			os.Exit(ErrCodeBucketDownload)
		}
		logging.Logger.Info("Artifact download from bucket successfully")
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
		// Set up the application
		setup()

		if err := validateArtifactUploadFlags(); err != nil {
			logging.Logger.Error("Invalid flags", "error", err)
			os.Exit(1)
		}

		destination := filepath.Join("artifacts", getMD5Hash(artifactsProject), artifactCommit)

		err := artifacts.Upload(artifactsPaths, destination)
		if err != nil {
			logging.Logger.Error("Artifacts upload to bucket failed", "error", err)
			os.Exit(ErrCodeBucketUpload)
		}
		logging.Logger.Info("Artifacts upload to bucket successfully")
	},
}

func init() {
	var err error

	// artifact global flags
	artifactsCmd.PersistentFlags().StringVar(&artifactCommit, "commit", "", "Commit hash in case of artifact")
	artifactsCmd.PersistentFlags().StringVar(&artifactsProject, "project", "", "Project name")
	artifactsCmd.PersistentFlags().StringSliceVar(&artifactsPaths, "path", []string{}, "List of directories and files")

	err = artifactsCmd.MarkPersistentFlagRequired("commit")
	if err != nil {
		slog.Error("Error marking flag as required", "error", err)
		return
	}
	err = artifactsCmd.MarkPersistentFlagRequired("project")
	if err != nil {
		slog.Error("Error marking flag as required", "error", err)
		return
	}
	err = artifactsCmd.MarkPersistentFlagRequired("path")
	if err != nil {
		slog.Error("Error marking flag as required", "error", err)
		return
	}

	// artifact download flags
	artifactDownloadCmd.PersistentFlags().StringVar(&artifactDestination, "destination", "", "Destination to extract the artifact")

	err = artifactDownloadCmd.MarkPersistentFlagRequired("destination")
	if err != nil {
		slog.Error("Error marking flag as required", "error", err)
		return
	}

	// add commands
	artifactsCmd.AddCommand(artifactDownloadCmd)
	artifactsCmd.AddCommand(artifactUploadCmd)
}
