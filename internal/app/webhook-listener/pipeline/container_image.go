package pipeline

import (
	"fmt"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/version"
)

// GetLauncherImage returns the image name and tag for the launcher image if format "name:tag"
func GetLauncherImage() string {
	return fmt.Sprintf("%s:%s", config.Launcher.Data.ImageName, func() string {
		if config.Launcher.Data.Tag == "" {
			return version.GetVersion()
		}
		return config.Launcher.Data.Tag
	}())
}
