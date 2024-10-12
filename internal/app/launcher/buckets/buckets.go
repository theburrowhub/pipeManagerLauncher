package buckets

import (
	"fmt"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

// TODO

func Download(source string) error {
	fmt.Println("Source:", source)
	logging.Logger.Info("TODO: Download")
	return nil
}

func Upload(paths []string, destination string) error {
	fmt.Println("Paths:", paths)
	fmt.Println("Destination:", destination)
	logging.Logger.Info("TODO: Upload")
	return nil
}
