package artifacts

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/sergiotejon/pipeManagerLauncher/internal/pkg/logging"
)

// packageTarFile creates a tar file with the given path
func packageTarFile(path string, destinationTarFile string) error {
	logging.Logger.Debug("Creating tar file", "tarFile", destinationTarFile, "path", path)

	// Overwrite the destination tar file
	tarFile, err := os.OpenFile(destinationTarFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer func(tarFile *os.File) {
		err := tarFile.Close()
		if err != nil {
			logging.Logger.Error("Error closing the tar file:", err)
		}
	}(tarFile)

	// Make a new gzip writer
	gzipWriter := gzip.NewWriter(tarFile)
	defer func(gzipWriter *gzip.Writer) {
		err := gzipWriter.Close()
		if err != nil {
			logging.Logger.Error("Error closing the gzip writer:", err)
		}
	}(gzipWriter)

	// Make a new tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer func(tarWriter *tar.Writer) {
		err := tarWriter.Close()
		if err != nil {
			logging.Logger.Error("Error closing the tar writer:", err)
		}
	}(tarWriter)

	// Add the files to the tar
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		logging.Logger.Debug("Walking the path for the tar file", "path", filePath)

		if err != nil {
			return err
		}

		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				logging.Logger.Error("Error closing the file:", err)
			}
		}(file)

		// Get file info
		info, err = file.Stat()
		if err != nil {
			return err
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		// Update header name to maintain directory structure
		header.Name = filePath

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// If not a directory, write file content
		if !info.IsDir() {
			if _, err := io.Copy(tarWriter, file); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	logging.Logger.Info("Files added to the tar file", "tarFile", destinationTarFile, "path", path)
	return nil
}

// extract extracts the files from a tar file to a destination folder
func extract(sourceTarFile, destinationFolder string) error {
	logging.Logger.Debug("Extracting tar file", "tarFile", sourceTarFile, "destination", destinationFolder)

	// Open the tar file
	tarFile, err := os.Open(sourceTarFile)
	if err != nil {
		return err
	}

	// Crea un lector gzip
	gzipReader, err := gzip.NewReader(tarFile)
	if err != nil {
		return err
	}
	defer func(gzipReader *gzip.Reader) {
		err := gzipReader.Close()
		if err != nil {
			logging.Logger.Error("Error closing the gzip reader:", err)
		}
	}(gzipReader)

	// Make a new tar reader
	tarReader := tar.NewReader(gzipReader)

	// Extract the files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(destinationFolder, header.Name)
		logging.Logger.Debug("Extracting file", "file", path)
		info := header.FileInfo()

		if info.IsDir() {
			err = os.MkdirAll(path, info.Mode())
			if err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, info.Mode())
		if err != nil {
			return err
		}

		if _, err := io.Copy(file, tarReader); err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}
	}

	logging.Logger.Info("Tar file extracted", "tarFile", sourceTarFile, "destination", destinationFolder)
	return nil
}
