package pipelineparser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sergiotejon/pipeManager/internal/pkg/envvars"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

// FindPipelineByRegex finds the pipeline to launch based on the variables
func FindPipelineByRegex(data map[string]interface{}, variables map[string]string) map[string]interface{} {
	// Create a map to store the pipelines that match the triggers
	pipelines := make(map[string]interface{})

	// Find the pipeline to launch
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
				pipelines[key] = make(map[string]interface{})
				// Include the global variables into the pipeline
				mergeMaps(pipelines[key].(map[string]interface{}), data["global"].(map[string]interface{}))
				// Merge the pipeline with global variables. Overwrite global variables with pipeline variables
				mergeMaps(pipelines[key].(map[string]interface{}), value.(map[string]interface{}))
				// Add the environment variables as parameters
				for paramName, paramValue := range convertEnvVarsIntoParams() {
					pipelines[key].(map[string]interface{})["params"].(map[string]interface{})[paramName] = paramValue
				}
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
