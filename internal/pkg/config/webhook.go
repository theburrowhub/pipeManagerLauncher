package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

// Route defines a single webhook route configuration.
// It captures the name, path, event type, and a list of event handlers.
type Route struct {
	Name          string  `yaml:"name"`                    // Name of the route (e.g., "github")
	Path          string  `yaml:"path"`                    // Path endpoint (e.g., "/github")
	EventType     string  `yaml:"eventType"`               // EventType is a CEL expression to determine the event
	GitSecretName string  `yaml:"gitSecretName,omitempty"` // GitSecretName is the name of the secret containing the Git credentials
	Events        []Event `yaml:"events"`                  // Events is a list of event handlers for this route
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
	Workers int     `yaml:"workers"` // Workers is the number of workers to process the incoming requests
	Routes  []Route `yaml:"routes"`  // Routes is a list of webhook routes
}

// WebhookConfig defines the webhook configuration.
// It captures the webhook configuration data from the root of the configuration file.
type WebhookConfig struct {
	Data WebhookStruct `yaml:"webhook"` // Data is the webhook configuration
}

var Webhook WebhookConfig // Webhook is the global webhook configuration

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
