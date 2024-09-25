package httpServer

import (
	"encoding/json"

	"github.com/sergiotejon/pipeManager/internal/app/webhook-listener/databuilder"
	"github.com/sergiotejon/pipeManager/internal/app/webhook-listener/pipeline"
	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

// processJob is the function that processes the incoming HTTP request
// It creates a PipelineData object and launches a job
// It returns an error if the job fails to launch
func processJob(job Job) error {
	var err error

	logging.Logger.Info("Request received", "method", job.Method, "path", job.Path)
	logging.Logger.Debug("Payload", "job", job)

	var jsonData []byte
	jsonData, err = json.MarshalIndent(job, "", "  ")
	if err != nil {
		return err
	}

	var pipelineData *databuilder.PipelineData
	pipelineData, err = databuilder.Run(jsonData, job.Path, config.Webhook.Data.Routes)
	if err != nil {
		return err
	}

	logging.Logger.Debug("Pipeline", "data", pipelineData)

	err = pipeline.LaunchJob(job.RequestID, pipelineData)
	if err != nil {
		return err
	}

	return nil
}
