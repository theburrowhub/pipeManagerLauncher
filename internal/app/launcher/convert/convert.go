package pipelinecrd_TODELETE

import (
	"encoding/json"
)

// ConvertToPipelines converts the raw data to a PipelineSpec struct
// This function convert only one raw pipeline to a PipelineSpec struct
// This would be the future version for normalize.go when it's refactored
func ConvertToPipelines(data interface{}) (PipelineSpec, error) {
	// Convert each item to YAML
	jsonData, err := json.Marshal(data)
	if err != nil {
		return PipelineSpec{}, err
	}

	// Unmarshal YAML to PipelineSpec struct
	var pipeline PipelineSpec
	err = json.Unmarshal(jsonData, &pipeline)
	if err != nil {
		return PipelineSpec{}, err
	}

	return pipeline, nil
}
