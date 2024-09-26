// Package config contains the configuration data structures and loading functions.
package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

// LogConfig defines the logging configuration.
// It captures the logging level, file, and format.
type LogConfig struct {
	Level  string `yaml:"level"`  // Level is the logging level (e.g., "info")
	File   string `yaml:"file"`   // File is the file to write the logs to
	Format string `yaml:"format"` // Format is the log format (e.g., "json")
}

// Route defines a single webhook route configuration.
// It captures the name, path, event type, and a list of event handlers.
type Route struct {
	Name      string  `yaml:"name"`      // Name of the route (e.g., "github")
	Path      string  `yaml:"path"`      // Path endpoint (e.g., "/github")
	EventType string  `yaml:"eventType"` // EventType is a CEL expression to determine the event
	Events    []Event `yaml:"events"`    // Events is a list of event handlers for this route
}

// Event represents a single event handler within a route.
// It captures the event type, repository, commit, and variables.
type Event struct {
	Type       string            `yaml:"type"`                 // Type of the event (e.g., "push")
	Repository string            `yaml:"repository"`           // Repository name
	Commit     string            `yaml:"commit,omitempty"`     // Commit hash (optional)
	DiffCommit string            `yaml:"diffCommit,omitempty"` // Commit hash to compare with the current commit (optional)
	Variables  map[string]string `yaml:"variables,omitempty"`  // Variables to be used in the event handler with its associated CEL expression (optional)
}

// WebhookStruct defines the webhook configuration.
// It captures the number of workers, logging configuration, and a list of routes.
type WebhookStruct struct {
	Workers int       `yaml:"workers"` // Workers is the number of workers to process the incoming requests
	Log     LogConfig `yaml:"log"`     // Log is the logging configuration
	Routes  []Route   `yaml:"routes"`  // Routes is a list of webhook routes
}

// WebhookConfig defines the webhook configuration.
// It captures the webhook configuration data from the root of the configuration file.
type WebhookConfig struct {
	Data WebhookStruct `yaml:"webhook"` // Data is the webhook configuration
}

// LauncherStruct defines the launcher configuration.
// It captures the image name, pull policy, tag, namespace, job name prefix, and timeout.
type LauncherStruct struct {
	ImageName     string   `yaml:"imageName"`     // ImageName is the name of the Docker image to be used
	PullPolicy    string   `yaml:"pullPolicy"`    // PullPolicy is the policy to use when pulling the image
	Tag           string   `yaml:"tag"`           // Tag is the tag of the Docker image to be used
	Namespace     string   `yaml:"namespace"`     // Namespace is the Kubernetes namespace to deploy the job
	JobNamePrefix string   `yaml:"jobNamePrefix"` // JobNamePrefix is the prefix to use for the job name
	Timeout       int64    `yaml:"timeout"`       // Timeout is the maximum time in seconds to wait for the job to complete
	GitSecret     []string `yaml:"gitSecret"`     // GitSecret is the list of names of secrets to use for the Git credentials
}

// LauncherConfig defines the launcher configuration.
// It captures the launcher configuration data from the root of the configuration file.
type LauncherConfig struct {
	Data LauncherStruct `yaml:"launcher"`
}

var (
	Webhook  WebhookConfig  // Webhook is the global webhook configuration
	Launcher LauncherConfig // Launcher is the global launcher configuration
)

// LoadWebhookConfig loads the webhook configuration from the given file.
// It returns an error if the configuration cannot be loaded.
// The configuration file is expected to be in YAML format.
// The configuration is loaded into the global Webhook variable.
func LoadWebhookConfig(configFile string) error {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yamlFile, &Webhook); err != nil {
		return err
	}

	if Webhook.Data.Workers == 0 {
		return errors.New("the configuration had not been loaded correctly")
	}

	return nil
}

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
