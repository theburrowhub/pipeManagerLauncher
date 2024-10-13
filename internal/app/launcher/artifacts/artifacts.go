// Package artifacts
// Manage artifacts and cache in the bucket by downloading and uploading content.
// It uses the Go Cloud Blob API to interact with the cloud storage bucket.
// The content is downloaded to a local directory and uploaded to the bucket.
// A tar file is used to store the content to be uploaded to the bucket.
package artifacts

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

var (
	bucketURL string // URL of the bucket
	basePath  string // Base path of the bucket
)

// setup initializes the bucket configuration
func setup() {
	// Compose the bucket URL and parameters
	bucketURL = fmt.Sprintf("%s?", config.Launcher.Data.ArtifactsBucket.URL)
	for key, value := range config.Launcher.Data.ArtifactsBucket.Parameters {
		separator := "&"
		if bucketURL[len(bucketURL)-1] == '?' {
			separator = ""
		}
		bucketURL = fmt.Sprintf("%s%s%s=%s", bucketURL, separator, key, value)
	}

	basePath = config.Launcher.Data.ArtifactsBucket.BasePath
}

// Download downloads the source file from the bucket and extracts it to the destination folder, keeping the tar file
// in the destinationTarFile to use it later in the upload process.
func Download(paths []string, bucketPath, destinationFolder string) error {
	// Set up the bucket configuration
	setup()

	for _, path := range paths {
		logging.Logger.Info("Downloading data from the bucket", "bucket", bucketURL, "path", path, "destinationFolder", destinationFolder)

		// Create a temporary directory to store the tar file
		tempDir, err := os.MkdirTemp("", "temp_tar_")
		if err != nil {
			return err
		}

		tarFileName := fmt.Sprintf("%s.tar.gz", getMD5Hash(path))
		bucketFile := filepath.Join(basePath, bucketPath, tarFileName)
		tarFullPath := filepath.Join(tempDir, tarFileName)

		// Download the source file from the bucket
		err = downloadFromBucket(bucketFile, tarFullPath)
		if err != nil {
			return err
		}

		// Extract the tar file
		err = extract(tarFullPath, destinationFolder)
		if err != nil {
			return err
		}

		logging.Logger.Info("Data downloaded from the bucket", "bucket", bucketURL, "path", path, "destinationFolder", destinationFolder)
	}

	return nil
}

// Upload uploads the paths to the bucket
func Upload(paths []string, bucketPath string) error {
	// Set up the bucket configuration
	setup()

	for _, path := range paths {
		logging.Logger.Info("Uploading path to the bucket", "buket", bucketURL, "path", path, "bucketPath", bucketPath)

		// Create a temporary directory to store the tar file
		tempDir, err := os.MkdirTemp("", "temp_tar_")
		if err != nil {
			return err
		}

		// Create the tar file
		tarFileName := fmt.Sprintf("%s.tar.gz", getMD5Hash(path))
		tarFullPath := filepath.Join(tempDir, tarFileName)
		err = packageTarFile(path, tarFullPath)
		if err != nil {
			return err
		}

		// Upload the tar file to the bucket
		bucketFile := filepath.Join(basePath, bucketPath, tarFileName)
		err = uploadToBucket(tarFullPath, bucketFile)
		if err != nil {
			return err
		}

		logging.Logger.Info("File uploaded to the bucket", "buket", bucketURL, "path", path, "tarFile", tarFullPath, "bucketFile", bucketFile)
	}

	return nil
}

// getMD5Hash returns the MD5 hash of the text
func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
