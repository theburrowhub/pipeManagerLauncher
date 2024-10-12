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
	cacheProject     string
	cachePaths       []string
	cacheDestination string
)

// cacheCmd represents the bucket command for cache
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage cache",
}

// validateCacheDownloadFlags checks if the provided flags are valid
func validateCacheDownloadFlags() error {
	if cacheProject == "" {
		return errors.New("project is required")
	}
	return nil
}

// cacheDownloadCmd represents the download subcommand
var cacheDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Cache download from bucket",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup the application
		setup()

		if err := validateCacheDownloadFlags(); err != nil {
			logging.Logger.Error("Invalid flags", "error", err)
			os.Exit(1)
		}

		bucketFolder := filepath.Join("artifacts", getMD5Hash(cacheProject))

		err := artifacts.Download(cachePaths, bucketFolder, cacheDestination)
		if err != nil {
			logging.Logger.Error("Cache download from bucket failed", "error", err)
			os.Exit(ErrCodeBucketDownload)
		}
		logging.Logger.Info("Cache download from bucket successful")
	},
}

// validateCacheUploadFlags checks if the provided flags are valid
func validateCacheUploadFlags() error {
	if cacheProject == "" {
		return errors.New("project is required")
	}
	if len(cachePaths) == 0 {
		return errors.New("paths for upload are required")
	}
	return nil
}

// cacheUploadCmd represents the upload subcommand
var cacheUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Cache upload to bucket",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup the application
		setup()

		if err := validateCacheUploadFlags(); err != nil {
			logging.Logger.Error("Invalid flags", "error", err)
			os.Exit(1)
		}

		destination := filepath.Join("cache", getMD5Hash(cacheProject))

		err := artifacts.Upload(cachePaths, destination)
		if err != nil {
			logging.Logger.Error("Cache upload to bucket failed", "error", err)
			os.Exit(ErrCodeBucketUpload)
		}
		logging.Logger.Info("Cache upload to bucket successful")
	},
}

func init() {
	var err error

	// cache flags
	cacheCmd.PersistentFlags().StringVar(&cacheProject, "project", "", "Project name")
	cacheCmd.PersistentFlags().StringSliceVar(&cachePaths, "path", []string{}, "List of directories and files")

	err = cacheCmd.MarkPersistentFlagRequired("project")
	if err != nil {
		slog.Error("Error marking flag as required", "error", err)
		return
	}
	err = cacheCmd.MarkPersistentFlagRequired("path")
	if err != nil {
		slog.Error("Error marking flag as required", "error", err)
		return
	}

	// download flags
	cacheDownloadCmd.PersistentFlags().StringVar(&cacheDestination, "destination", "", "The destination path to extract the cache")

	err = cacheDownloadCmd.MarkPersistentFlagRequired("destination")
	if err != nil {
		slog.Error("Error marking flag as required", "error", err)
		return
	}

	// commands
	cacheCmd.AddCommand(cacheDownloadCmd)
	cacheCmd.AddCommand(cacheUploadCmd)
}
