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

// JobConfig contains the configuration for a new Kubernetes Job
type JobConfig struct {
	JobName         string
	RequestID       string
	PipelineData    *databuilder.PipelineData
	Namespace       string
	JobTimeout      int64
	ContainerName   string
	Env             []corev1.EnvVar
	ConfigmapName   string
	ImagePullPolicy string
}

// getKubernetesClient returns a Kubernetes clientset configured for either in-cluster or local access
func getKubernetesClient() (*kubernetes.Clientset, error) {
	// Try to get the in-cluster config
	cfg, err := rest.InClusterConfig()
	if err != nil {
		logging.Logger.Warn("Failed to get in-cluster config, trying local config", "error", err)
		// Fallback to local kubeconfig
		var configPath string
		if os.Getenv("KUBECONFIG") != "" { // Get the kubeconfig path from the KUBECONFIG environment variable
			configPath = os.Getenv("KUBECONFIG")
		} else { // If not set, use the default path ~/.kube/config
			configPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}
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

// pointerInt32 returns a pointer to an int32
func pointerInt32(i int32) *int32 {
	return &i
}

// createJobObject creates a Kubernetes Job object with the given parameters
func createJobObject(job *JobConfig) *batchv1.Job {
	// If the GitSecretName is empty, use an emptyDir volume. Otherwise, use a secret volume
	var gitSecretVolume corev1.VolumeSource
	if job.PipelineData.GitSecretName == "" {
		gitSecretVolume = corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		}
	} else {
		gitSecretVolume = corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  job.PipelineData.GitSecretName,
				DefaultMode: pointerInt32(0o600),
			},
		}
	}

	// Create the Job object
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      job.JobName,
			Labels:    getLabels(job.RequestID, job.PipelineData),
			Namespace: job.Namespace,
		},
		Spec: batchv1.JobSpec{
			ActiveDeadlineSeconds: &job.JobTimeout,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getLabels(job.RequestID, job.PipelineData),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            job.ContainerName,
							Image:           GetLauncherImage(),
							ImagePullPolicy: corev1.PullPolicy(job.ImagePullPolicy),
							Env:             job.Env, // Environment variables with the pipeline data
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config-volume",
									MountPath: "/etc/pipe-manager",
								},
								{
									Name:      "git-credentials",
									MountPath: "/root/.ssh",
								},
								{
									Name:      "repo-storage",
									MountPath: "/tmp/repo",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: []corev1.Volume{
						{
							Name: "config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: job.ConfigmapName,
									},
								},
							},
						},
						{
							Name:         "git-credentials",
							VolumeSource: gitSecretVolume,
						},
						{
							Name: "repo-storage",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
}
