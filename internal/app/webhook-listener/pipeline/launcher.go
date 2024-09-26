// Package pipeline contains the implementation of the pipeline launcher that creates a new Kubernetes Job with the given request ID and pipeline data.
package pipeline

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sergiotejon/pipeManager/internal/app/webhook-listener/databuilder"
	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

const containerName = "launcher"

var jobCommand = []string{"sh", "-c", "export && sleep 60"}

// LaunchJob creates a new Kubernetes Job with the given request ID and pipeline data
// It returns an error if the job cannot be created
// The request ID is the unique identifier of the http request coming from the webhook
// The pipeline data contains the pipeline name, path, event, repository, commit, and variables
func LaunchJob(requestID string, pipelineData *databuilder.PipelineData) error {
	jobName := config.Launcher.Data.JobNamePrefix + "-" + requestID
	namespace := config.Launcher.Data.Namespace
	jobTimeout := config.Launcher.Data.Timeout

	// Get the Kubernetes client
	client, err := getKubernetesClient()
	if err != nil {
		return err
	}

	// Convert the environment variables map into an array of corev1.EnvVar objects
	env := getEnvVarsFromPipelineData(pipelineData)

	// Get the current namespace if not provided. "default" if not found
	if namespace == "" {
		namespace, err = getCurrentNamespace()
		if err != nil {
			logging.Logger.Warn("Error getting current namespace", "error", err, "defaultNamespace", namespace)
		}
	}

	// Job definition
	// ** TODO: Create a kubernetes controller to manage a new object type called, for example, "Pipeline". That way, we can manage the pipeline lifecycle
	// ** from the creation to the deletion of the resources. This controller will be responsible for creating the Tekton Pipeline and manage the resources
	// ** created by the pipeline.
	job := createJobObject(jobName, requestID, pipelineData, namespace, jobTimeout, containerName, jobCommand, env)

	// Build the Job
	jobClient := client.BatchV1().Jobs(namespace)
	result, err := jobClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	logging.Logger.Info("Pipeline launcher", "job", result.GetObjectMeta().GetName(), "namespace", namespace)

	return nil
}
