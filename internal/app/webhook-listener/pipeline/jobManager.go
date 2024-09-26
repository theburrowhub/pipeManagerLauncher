package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/sergiotejon/pipeManager/internal/app/webhook-listener/databuilder"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
	"github.com/sergiotejon/pipeManager/internal/pkg/version"
)

// getKubernetesClient returns a Kubernetes clientset configured for either in-cluster or local access
func getKubernetesClient() (*kubernetes.Clientset, error) {
	// Try to get the in-cluster config
	cfg, err := rest.InClusterConfig()
	if err != nil {
		logging.Logger.Warn("Failed to get in-cluster config, trying local config", "error", err)
		// Fallback to local kubeconfig
		configPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			return nil, err
		}
	}

	// Create a clientset for interacting with the Kubernetes API
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

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

// getCurrentNamespace returns the current namespace where the pod is running
func getCurrentNamespace() (string, error) {
	namespaceFile := "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	namespace, err := os.ReadFile(namespaceFile)
	if err != nil {
		return "default", err
	}
	return string(namespace), nil
}

// getEnvVarsFromPipelineData converts the pipeline data into a slice of corev1.EnvVar
func getEnvVarsFromPipelineData(pipelineData *databuilder.PipelineData) []corev1.EnvVar {
	var env []corev1.EnvVar
	for key, value := range pipelineData.Variables {
		env = append(env, corev1.EnvVar{
			Name:  fmt.Sprintf("PIPELINE_VARIABLE_%s", strings.ToUpper(key)),
			Value: value,
		})
	}

	env = append(env, corev1.EnvVar{
		Name:  "PIPELINE_COMMIT",
		Value: pipelineData.Commit,
	})
	env = append(env, corev1.EnvVar{
		Name:  "PIPELINE_DIFF_COMMIT",
		Value: pipelineData.DiffCommit,
	})
	env = append(env, corev1.EnvVar{
		Name:  "PIPELINE_REPOSITORY",
		Value: pipelineData.Repository,
	})

	return env
}

// createJobObject creates a Kubernetes Job object with the given parameters
func createJobObject(
	jobName,
	requestID string,
	pipelineData *databuilder.PipelineData,
	namespace string,
	jobTimeout int64,
	containerName string,
	jobCommand []string,
	env []corev1.EnvVar) *batchv1.Job {

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Labels:    getLabels(requestID, pipelineData),
			Namespace: namespace,
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
}
