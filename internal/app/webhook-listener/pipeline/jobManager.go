package pipeline

import (
	"github.com/sergiotejon/pipeManager/internal/app/webhook-listener/databuilder"
	"github.com/sergiotejon/pipeManager/internal/pkg/version"
)

// getLabels returns a map of labels to be used in Kubernetes objects
func getLabels(requestID string, pipelineData *databuilder.PipelineData) map[string]string {
	return map[string]string{
		"handleBy":               "pipeManager",
		"pipe-manager/Version":   version.GetVersion(),
		"pipe-manager/RequestID": requestID,
		"pipe-manager/Route":     pipelineData.Name,
		"pipe-manager/Event":     pipelineData.Event,
	}
}
