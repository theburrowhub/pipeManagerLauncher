package artifacts

import (
	"context"
	"io"
	"os"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"

	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

// DownloadFromBucket downloads a file from the bucket
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

// UploadToBucket uploads a file to the bucket
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
