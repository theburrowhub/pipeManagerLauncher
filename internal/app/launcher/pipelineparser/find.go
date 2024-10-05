package pipelineparser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sergiotejon/pipeManager/internal/pkg/envvars"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

// FindPipelineByName finds the pipeline to launch based on the name
func FindPipelineByName(data map[string]interface{}, variables map[string]string, pipelineName string) map[string]interface{} {
	// Create a map to store the pipelines that match the triggers
	pipelines := make(map[string]interface{})

	// Find the pipeline to launch
	// Key is the name of the pipeline or "global"
	// Value is the pipeline data
	for key, value := range data {
		// If the key is "global", skip it
		if key == "global" {
			continue
		}

		if key == pipelineName {
			pipelines[key] = createAtomicPipeline(data, value)
		}
	}

	return pipelines
}

// FindPipelineByRegex finds the pipeline to launch based on the variables
func FindPipelineByRegex(data map[string]interface{}, variables map[string]string) map[string]interface{} {
	// Create a map to store the pipelines that match the triggers
	pipelines := make(map[string]interface{})

	// Find the pipeline to launch
	// Key is the name of the pipeline or "global"
	// Value is the pipeline data
	for key, value := range data {
		// If the key is "global", skip it
		if key == "global" {
			continue
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// Get the pipeline triggers
			triggers := value.(map[string]interface{})["pipelineTriggers"]
			if triggers == nil {
				logging.Logger.Warn("No pipeline triggers found", "pipeline", key)
				continue
			}

			// Iterate over the triggers
			triggerList := triggers.([]interface{})
			addPipeline := true
			for _, trigger := range triggerList {
				triggerMap := trigger.(map[string]interface{})
				variableName := strings.ToUpper(triggerMap["variableName"].(string))
				variableRegex := triggerMap["valueRegex"].(string)

				// Check if the variable matches the regex
				matched, err := regexp.MatchString(variableRegex, variables[variableName])
				if err != nil {
					logging.Logger.Debug("Invalid regex",
						"variable", variableName, "regex", variableRegex, "pipeline", key)
					continue
				}
				// It the variable does match, continue with the next trigger
				// If it doesn't match, break the loop and don't add the pipeline
				if matched {
					logging.Logger.Debug("Trigger matched",
						"variable", variableName, "regex", variableRegex, "pipeline", key)
				} else {
					logging.Logger.Debug("Trigger not matched",
						"variable", variableName, "regex", variableRegex, "pipeline", key)
					addPipeline = false
					break
				}
			}

			// Add the pipeline to the list if all triggers matched
			if addPipeline {
				pipelines[key] = createAtomicPipeline(data, value)
			}
		default:
			logging.Logger.Warn("Unexpected type found", "type", fmt.Sprintf("%T", v), "pipeline", key)
		}
	}

	return pipelines
}

// convertEnvVarsIntoParams converts the environment variables into parameters of pipeline
// Return two strings, the first one is the name of the parameter and the second one is the value
// of the parameter
func convertEnvVarsIntoParams() map[string]string {
	params := make(map[string]string)

	for key, value := range envvars.Variables {
		paramName := strings.ReplaceAll(key, "_", "-")
		paramValue := value

		params[paramName] = paramValue
	}

	return params
}

// createAtomicPipeline creates a pipeline with the global variables and the pipeline variables
func createAtomicPipeline(data map[string]interface{}, value interface{}) map[string]interface{} {
	pipeline := make(map[string]interface{})

	// Include the global variables into the pipeline
	mergeMaps(pipeline, data["global"].(map[string]interface{}))

	// Merge the pipeline with global variables. Overwrite global variables with pipeline variables
	mergeMaps(pipeline, value.(map[string]interface{}))

	// Add the environment variables as parameters
	for paramName, paramValue := range convertEnvVarsIntoParams() {
		pipeline["params"].(map[string]interface{})[paramName] = paramValue
	}

	// Remove the pipeline triggers if they exist
	delete(pipeline, "pipelineTriggers")

	return pipeline
}
