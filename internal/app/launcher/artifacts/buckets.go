package buckets

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

var (
	bucketURL string // URL of the bucket
	basePath  string // Base path of the bucket
)

// setup initializes the bucket configuration
func setup() {
	bucketURL = config.Launcher.Data.Bucket.URL
	basePath = config.Launcher.Data.Bucket.BasePath
}

// Download downloads the source file from the bucket and extracts it to the destination folder, keeping the tar file
// in the destinationTarFile to use it later in the upload process.
func Download(source, destinationTarFile, destinationFolder string) error {
	var err error

	// Setup the bucket configuration
	setup()

	sourceURL := filepath.Join(bucketURL, basePath, source)
	logging.Logger.Info("Downloading source from the bucket", "source", sourceURL, "destination tar file", destinationTarFile, "destination", destinationFolder)

	// Download the source file from the bucket
	err = downloadFromBucket(source, destinationTarFile)
	if err != nil {
		return err
	}

	// Extract the tar file
	err = untarPath(destinationTarFile, destinationFolder)
	if err != nil {
		return err
	}

	logging.Logger.Info("Source downloaded from the bucket", "source", sourceURL, "destination tar file", destinationTarFile, "destination", destinationFolder)

	return nil
}

// Upload uploads the paths to the bucket in the destination folder. The paths are added to it are added to the tar file
func Upload(paths []string, sourceTarFile, destination string) error {
	var err error

	// Setup the bucket configuration
	setup()

	destinationURL := filepath.Join(bucketURL, basePath, destination)
	logging.Logger.Info("Uploading paths to the bucket", "paths", paths, "source tar file", sourceTarFile, "destination", destinationURL)

	err = tarPaths(paths, sourceTarFile)
	if err != nil {
		return err
	}

	err = uploadToBucket(sourceTarFile, destinationURL)
	if err != nil {
		return err
	}

	logging.Logger.Info("Paths uploaded to the bucket", "paths", paths, "source tar file", sourceTarFile, "destination", destinationURL)

	return nil
}

func downloadFromBucket(source, destination string) error {
	ctx := context.Background()

	// Open a connection to the bucket
	bucket, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		return err
	}
	defer func(bucket *blob.Bucket) {
		err := bucket.Close()
		if err != nil {
			logging.Logger.Error("Error closing bucket", "error", err)
		}
	}(bucket)

	// Create the destination file
	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			logging.Logger.Error("Error closing the destination file", "error", err)
		}
	}(destFile)

	// Download the file from the bucket
	reader, err := bucket.NewReader(ctx, source, nil)
	if err != nil {
		return err
	}
	defer func(reader *blob.Reader) {
		err := reader.Close()
		if err != nil {
			logging.Logger.Error("Error closing the reader", "error", err)
		}
	}(reader)

	// Copy the content to the destination file
	if _, err := io.Copy(destFile, reader); err != nil {
		return err
	}

	logging.Logger.Info("File downloaded from the bucket", "source", source, "destination", destination)
	return nil
}

func uploadToBucket(source, destination string) error {
	ctx := context.Background()

	// Open a connection to the bucket
	bucket, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		return err
	}
	defer func(bucket *blob.Bucket) {
		err := bucket.Close()
		if err != nil {
			logging.Logger.Error("Error closing bucket", "error", err)
		}
	}(bucket)

	// Create the source file
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			logging.Logger.Error("Error closing the source file", "error", err)
		}
	}(sourceFile)

	// Upload the file to the bucket
	writer, err := bucket.NewWriter(ctx, destination, nil)
	if err != nil {
		return err
	}
	defer func(writer *blob.Writer) {
		err := writer.Close()
		if err != nil {
			logging.Logger.Error("Error closing the writer", "error", err)
		}
	}(writer)

	// Copy the content to the destination file
	if _, err := io.Copy(writer, sourceFile); err != nil {
		return err
	}

	logging.Logger.Info("File uploaded to the bucket", "source", source, "destination", destination)
	return nil
}
