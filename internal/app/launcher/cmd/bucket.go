package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/sergiotejon/pipeManager/internal/app/launcher/buckets"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

var (
	bucketType     string
	artifactCommit string
	project        string
)

// bucketCmd represents the bucket command
var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Bucket operations",
}

// validateFlags checks if the provided flags are valid
func validateFlags() error {
	if bucketType != "artifact" && bucketType != "cache" {
		return errors.New("invalid type: must be 'artifact' or 'cache'")
	}
	if bucketType == "artifact" && artifactCommit == "" {
		return errors.New("commit is required when type is 'artifact'")
	}
	return nil
}

// downloadCmd represents the download subcommand
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download from bucket",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup the application
		setup()

		if err := validateFlags(); err != nil {
			logging.Logger.Error("Invalid flags", "error", err)
			os.Exit(1)
		}

		err := buckets.Download()
		if err != nil {
			logging.Logger.Error("Download from bucket failed", "error", err)
			os.Exit(ErrCodeBucketDownload)
		}
		logging.Logger.Info("Download from bucket successful")
	},
}

// uploadCmd represents the upload subcommand
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload to bucket",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup the application
		setup()

		if err := validateFlags(); err != nil {
			logging.Logger.Error("Invalid flags", "error", err)
			os.Exit(1)
		}

		err := buckets.Upload()
		if err != nil {
			logging.Logger.Error("Upload to bucket failed", "error", err)
			os.Exit(ErrCodeBucketUpload)
		}
		logging.Logger.Info("Upload to bucket successful")
	},
}

func init() {
	var err error

	bucketCmd.PersistentFlags().StringVar(&bucketType, "type", "", "Type of bucket operation (artifact, cache)")
	bucketCmd.PersistentFlags().StringVar(&artifactCommit, "commit", "", "Commit hash in case of artifact")
	bucketCmd.PersistentFlags().StringVar(&project, "project", "", "Project name")

	err = bucketCmd.MarkPersistentFlagRequired("type")
	if err != nil {
		logging.Logger.Error("Error marking flag as required", "error", err)
		return
	}
	err = bucketCmd.MarkPersistentFlagRequired("project")
	if err != nil {
		logging.Logger.Error("Error marking flag as required", "error", err)
		return
	}

	bucketCmd.AddCommand(downloadCmd)
	bucketCmd.AddCommand(uploadCmd)
}
