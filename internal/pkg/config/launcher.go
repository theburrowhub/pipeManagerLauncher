package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

// LauncherStruct defines the launcher configuration.
// It captures the image name, pull policy, tag, namespace, job name prefix, and timeout.
type LauncherStruct struct {
	ImageName     string `yaml:"imageName"`     // ImageName is the name of the Docker image to be used
	PullPolicy    string `yaml:"pullPolicy"`    // PullPolicy is the policy to use when pulling the image
	Tag           string `yaml:"tag"`           // Tag is the tag of the Docker image to be used
	Namespace     string `yaml:"namespace"`     // Namespace is the Kubernetes namespace to deploy the job
	JobNamePrefix string `yaml:"jobNamePrefix"` // JobNamePrefix is the prefix to use for the job name
	Timeout       int64  `yaml:"timeout"`       // Timeout is the maximum time in seconds to wait for the job to complete
	GitSecretPath string `yaml:"gitSecretPath"` // GitSecretPath is the path to the directory containing the Git credentials
}

// LauncherConfig defines the launcher configuration.
// It captures the launcher configuration data from the root of the configuration file.
type LauncherConfig struct {
	Data LauncherStruct `yaml:"launcher"`
}

var Launcher LauncherConfig // Launcher is the global launcher configuration

// LoadLauncherConfig loads the launcher configuration from the given file.
// It returns an error if the configuration cannot be loaded.
// The configuration file is expected to be in YAML format.
// The configuration is loaded into the global Launcher variable.
func LoadLauncherConfig(configFile string) error {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yamlFile, &Launcher); err != nil {
		return err
	}

	if Launcher.Data.ImageName == "" {
		return errors.New("the configuration had not been loaded correctly")
	}

	return nil
}
