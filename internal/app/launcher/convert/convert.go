package convert

import (
	"encoding/json"

	pipemanagerv1alpha1 "github.com/sergiotejon/pipeManagerController/api/v1alpha1"
)

// ConvertToPipelines converts the raw data to a PipelineSpec struct
// This function convert only one raw pipeline to a PipelineSpec struct
// This would be the future version for normalize.go when it's refactored
func ConvertToPipelines(data interface{}) (pipemanagerv1alpha1.PipelineSpec, error) {
	// Convert each item to YAML
	jsonData, err := json.Marshal(data)
	if err != nil {
		return pipemanagerv1alpha1.PipelineSpec{}, err
	}

	// Unmarshal YAML to PipelineSpec struct
	var pipeline pipemanagerv1alpha1.PipelineSpec
	err = json.Unmarshal(jsonData, &pipeline)
	if err != nil {
		return pipemanagerv1alpha1.PipelineSpec{}, err
	}

	return pipeline, nil
}
