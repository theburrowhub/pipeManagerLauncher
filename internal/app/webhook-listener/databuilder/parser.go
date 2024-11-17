// Package databuilder contains the implementation of the data builder that parses the webhook payload and routes configuration.
package databuilder

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sergiotejon/pipeManagerLauncher/internal/pkg/logging"
	"github.com/sergiotejon/pipeManagerLauncher/pkg/config"
)

// PipelineData represents a pipeline to be executed
// It contains the name, path, event, repository, commit, and variables (a map of variable names and their values)
type PipelineData struct {
	Name          string
	Path          string
	Event         string
	Repository    string
	GitSecretName string
	Commit        string
	DiffCommit    string
	Variables     map[string]string
}

// Run executes the parser with the given payload and routes configuration returning a Pipeline
// It returns an error if the payload cannot be unmarshalled, the route cannot be found, the event route cannot be found,
// or the CEL expression cannot be evaluated
// It returns the PipelineData object if the parser is successful
// The payload is the JSON data to be parsed from the webhook
// The path is the path of the route to be executed from the request
// The routes are the list of routes to be executed with their associated events and variables from the configuration
func Run(payload json.RawMessage, path string, routes []config.Route) (*PipelineData, error) {
	var err error

	// Unmarshal JSON data into a map
	var jsonData map[string]interface{}
	if err = json.Unmarshal(payload, &jsonData); err != nil {
		return nil, err
	}

	// Retrieve the route from the routes configuration
	var route *config.Route
	route, err = getConfiguredRouteByPath(path, routes)
	if err != nil {
		return nil, err
	}

	// Retrieve the event value from the route
	var eventType string
	eventType, err = evaluateCELExpression(route.EventType, jsonData)
	if err != nil {
		return nil, err
	}

	// Evaluate the git secret name from the route
	gitSecretName := ""
	if route.GitSecretName != "" {
		gitSecretName, err = evaluateCELExpression(route.GitSecretName, jsonData)
		if err != nil {
			return nil, err
		}
	}

	// Retrieve the event route information from the events configuration
	var event *config.Event
	event, err = getEventRouteByEventType(eventType, route.Events)
	if err != nil {
		return nil, err
	}

	// Evaluate the repository from the event route
	var repository string
	repository, err = evaluateCELExpression(event.Repository, jsonData)
	if err != nil {
		return nil, err
	}

	// Evaluate the commits from the event route. If the commit is not defined, it will be an empty string
	commit := ""
	if event.Commit != "" {
		commit, err = evaluateCELExpression(event.Commit, jsonData)
		if err != nil {
			return nil, err
		}
	}
	diffCommit := ""
	if event.DiffCommit != "" {
		diffCommit, err = evaluateCELExpression(event.DiffCommit, jsonData)
		if err != nil {
			return nil, err
		}
	}

	logging.Logger.Info("Data Builder", "route", route.Name, "eventType", eventType, "repository", repository, "commit", commit)

	// Create a PipelineData object
	pipelineData := PipelineData{
		Name:          route.Name,
		Path:          route.Path,
		Event:         eventType,
		Repository:    repository,
		GitSecretName: gitSecretName,
		Commit:        commit,
		DiffCommit:    diffCommit,
		Variables:     make(map[string]string),
	}

	// Evaluate the variables and store them in the PipelineData object
	for key, value := range event.Variables {
		pipelineData.Variables[key], err = evaluateCELExpression(value, jsonData)
		if err != nil {
			return nil, err
		}

		logging.Logger.Info("Data Builder", "variableName", key,
			"celExpression", value,
			"valueRetrieved", pipelineData.Variables[key],
			"event", eventType,
			"route", route.Name,
		)
	}

	return &pipelineData, nil
}

// Retrieve the route from the routes configuration
func getConfiguredRouteByPath(path string, routes []config.Route) (*config.Route, error) {
	for _, route := range routes {
		if route.Path == path {
			return &route, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("route path '%s' not found", path))
}

// getEventRouteByEventType retrieves the event route from the events configuration
// It returns an error if the event route is not found
// It returns the event route if it is found
// The eventType is the type of the event to be retrieved
// The events are the list of events to be searched. These are the events associated with the route and its path from the configuration
func getEventRouteByEventType(eventType string, events []config.Event) (*config.Event, error) {
	for _, event := range events {
		if event.Type == eventType {
			return &event, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("event route '%s' not found", eventType))
}
