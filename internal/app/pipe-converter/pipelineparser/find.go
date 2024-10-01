package pipelineparser

import (
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
	"regexp"
	"strings"
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
		}
	}

	return pipelines
}
