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
	cacheProject string
	cachePaths   []string
	cacheTarFile string
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
	if cacheTarFile == "" {
		return errors.New("tar file is required")
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

		fileName := fmt.Sprintf("%s.tar.gz", cacheProject)
		source := filepath.Join("cache", fileName)

		err := buckets.Download(source, cacheTarFile)
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
	if cacheTarFile == "" {
		return errors.New("tar file is required")
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

		fileName := fmt.Sprintf("%s.tar.gz", cacheProject)
		destination := filepath.Join("cache", fileName)

		err := buckets.Upload(cachePaths, cacheTarFile, destination)
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

	err = cacheCmd.MarkPersistentFlagRequired("project")
	if err != nil {
		logging.Logger.Error("Error marking flag as required", "error", err)
		return
	}

	// download flags
	cacheDownloadCmd.PersistentFlags().StringVar(&cacheTarFile, "tar", "", "The name of the tar file where the cache will be downloaded")

	err = cacheDownloadCmd.MarkPersistentFlagRequired("tar")
	if err != nil {
		logging.Logger.Error("Error marking flag as required", "error", err)
		return
	}

	// upload flags
	cacheUploadCmd.PersistentFlags().StringSliceVar(&cachePaths, "path", []string{}, "List of directories and files")
	cacheUploadCmd.PersistentFlags().StringVar(&cacheTarFile, "tar", "", "Tar file to add the path to")

	err = cacheUploadCmd.MarkPersistentFlagRequired("path")
	if err != nil {
		logging.Logger.Error("Error marking flag as required", "error", err)
		return
	}
	err = cacheUploadCmd.MarkPersistentFlagRequired("tar")
	if err != nil {
		logging.Logger.Error("Error marking flag as required", "error", err)
		return
	}

	cacheCmd.AddCommand(cacheDownloadCmd)
	cacheCmd.AddCommand(cacheUploadCmd)
}
