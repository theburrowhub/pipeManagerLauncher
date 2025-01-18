package config

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"

	"github.com/sergiotejon/pipeManagerLauncher/internal/pkg/version"
)

// LauncherStruct defines the launcher configuration.
// It captures the image name, pull policy, tag, namespace, job name prefix, and timeout.
type LauncherStruct struct {
	ImageName       string       `json:"imageName"`       // ImageName is the name of the Docker image to be used
	PullPolicy      string       `json:"pullPolicy"`      // PullPolicy is the policy to use when pulling the image
	Tag             string       `json:"tag"`             // Tag is the tag of the Docker image to be used
	Namespace       string       `json:"namespace"`       // Namespace is the Kubernetes namespace to deploy the job
	JobNamePrefix   string       `json:"jobNamePrefix"`   // JobNamePrefix is the prefix to use for the job name
	Timeout         int64        `json:"timeout"`         // Timeout is the maximum time in seconds to wait for the job to complete
	BackoffLimit    int32        `json:"backoffLimit"`    // BackoffLimit is the number of retries before considering the job as failed
	ConfigmapName   string       `json:"configmapName"`   // ConfigmapName is the name of the ConfigMap to use
	CloneDepth      int          `json:"cloneDepth"`      // CloneDepth is the depth to use when cloning the Git repository
	RolesBinding    []string     `json:"rolesBinding"`    // RolesBinding is the list of roles to bind to the Service Account
	ArtifactsBucket BucketConfig `json:"artifactsBucket"` // ArtifactsBucket is the bucket configuration for storing the artifacts
}

// BucketConfig defines the bucket configuration.
type BucketConfig struct {
	URL         string            `json:"url"`                   // URL is the URL of the bucket
	BasePath    string            `json:"basePath"`              // BasePath is the name of the bucket
	SecretName  string            `json:"secretName"`            // SecretName is the name of the secret to use when accessing the bucket
	Parameters  map[string]string `json:"parameters,omitempty"`  // Parameters is a map of additional parameters for the bucket
	Credentials BucketCredentials `json:"credentials,omitempty"` // Credentials is the credentials to use when accessing the bucket
}

// BucketCredentials defines the bucket credentials.
// It captures the environment variables, volumes, and volume mounts to use.
// Usually, the credentials are stored in a secret and mounted as a volume or as environment variables.
// Mutually exclusive, but both can be used.
type BucketCredentials struct {
	Env          []corev1.EnvVar      `json:"env,omitempty"`          // Env is the environment variables to use
	Volumes      []corev1.Volume      `json:"volumes,omitempty"`      // Volumes is the volumes to use
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"` // VolumeMounts is the volume mounts to use
}

// LauncherConfig defines the launcher configuration.
// It captures the launcher configuration data from the root of the configuration file.
type LauncherConfig struct {
	Data LauncherStruct `json:"launcher"`
}

var Launcher LauncherConfig // Launcher is the global launcher configuration

// LoadLauncherConfig loads the launcher configuration from the given file.
// It returns an error if the configuration cannot be loaded.
// The configuration file is expected to be in YAML format if a file is provided.
// If the configuration file is not provided, the configuration is loaded from the environment variables.
// The configuration loads into the global Launcher variable.
func LoadLauncherConfig(configFile string) error {
	if configFile == "" { // Load the configuration from the environment variables
		loadConfigFromEnv(&Launcher, "LAUNCHER")
	} else { // Load the configuration from the file
		err := loadConfigFromFile(configFile, &Launcher)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetLauncherImage returns the image name and tag for the launcher image if format "name:tag"
func (l *LauncherStruct) GetLauncherImage() string {
	return fmt.Sprintf("%s:%s", l.ImageName, func() string {
		if l.Tag == "" {
			return version.GetVersion()
		}
		return l.Tag
	}())
}
