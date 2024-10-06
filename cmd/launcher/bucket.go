package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/sergiotejon/pipeManager/internal/app/launcher/buckets"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

// bucketCmd represents the bucket command
var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Bucket operations",
}

// downloadCmd represents the download subcommand
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download from bucket",
	Run: func(cmd *cobra.Command, args []string) {
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
		initApp()

		err := buckets.Upload()
		if err != nil {
			logging.Logger.Error("Upload to bucket failed", "error", err)
			os.Exit(ErrCodeBucketUpload)
		}
		logging.Logger.Info("Upload to bucket successful")
	},
}

func init() {
	bucketCmd.AddCommand(downloadCmd, uploadCmd)
}
