// Package pipeline contains the implementation of the pipeline launcher that creates a new Kubernetes Job with the given request ID and pipeline data.
package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

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

	// Load kubeconfig
	configPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	cfg, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return err
	}

	// Client object for interacting with Kubernetes API
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	// Convert the environment variables map into an array of corev1.EnvVar objects
	var env []corev1.EnvVar
	for key, value := range pipelineData.Variables {
		env = append(env, corev1.EnvVar{
			Name:  fmt.Sprintf("PIPELINE_VARIABLE_%s", strings.ToUpper(key)),
			Value: value,
		})
	}

	// Add the pipeline data to the environment variables (commit and repository)
	env = append(env, corev1.EnvVar{
		Name:  "PIPELINE_COMMIT",
		Value: pipelineData.Commit,
	})
	env = append(env, corev1.EnvVar{
		Name:  "PIPELINE_REPOSITORY",
		Value: pipelineData.Repository,
	})

	// Job definition
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   jobName,
			Labels: getLabels(requestID, pipelineData),
		},
		Spec: batchv1.JobSpec{
			ActiveDeadlineSeconds: &jobTimeout,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getLabels(requestID, pipelineData),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    containerName,
							Image:   GetLauncherImage(),
							Command: jobCommand,
							Env:     env, // Environment variables with the pipeline data
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	// Build the Job
	if namespace == "" {
		namespace = "default"
		if ns, err := client.CoreV1().Namespaces().Get(context.TODO(), "default", metav1.GetOptions{}); err == nil && ns != nil {
			namespace = ns.Name
		}
	}
	jobClient := client.BatchV1().Jobs(namespace)
	result, err := jobClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	logging.Logger.Info("Pipeline launcher", "job", result.GetObjectMeta().GetName(), "namespace", namespace)

	return nil
}
