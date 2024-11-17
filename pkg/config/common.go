// Package config contains the configuration data structures and loading functions.
package config

// LogConfig defines the logging configuration.
// It captures the logging level, file, and format.
type LogConfig struct {
	Level  string `yaml:"level"`  // Level is the logging level (e.g., "info")
	File   string `yaml:"file"`   // File is the file to write the logs to
	Format string `yaml:"format"` // Format is the log format (e.g., "json")
}

// CommonStruct defines the common configuration.
type CommonStruct struct {
	Log LogConfig `yaml:"log"` // Log is the logging configuration
}

// CommonConfig defines the common configuration.
type CommonConfig struct {
	Data CommonStruct `yaml:"common"` // Data is the common configuration data
}

var Common CommonConfig // Common is the common configuration for all components

// LoadCommonConfig loads the common configuration from the given file.
// It returns an error if the configuration cannot be loaded.
// The configuration file is expected to be in YAML format.
// It allows to load from environment variables also.
// The configuration is loaded into the global Common variable.
func LoadCommonConfig(configFile string) error {
	if configFile == "" { // Load the configuration from the environment variables
		loadConfigFromEnv(&Common, "COMMON")
	} else { // Load the configuration from the file
		err := loadConfigFromFile(configFile, &Common)
		if err != nil {
			return err
		}
	}
	return nil
}
