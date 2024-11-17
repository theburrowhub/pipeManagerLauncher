package namespace

import (
	pipemanagerv1alpha1 "github.com/sergiotejon/pipeManagerController/api/v1alpha1"

	"github.com/sergiotejon/pipeManagerLauncher/pkg/config"
)

// getBucketCredentialsSecretFromConfig returns the names of the secrets that are used in the namespace for the artifacts bucket and the
// credentials secret.
func getBucketCredentialsSecretFromConfig() []string {
	var secretNames []string

	bc := config.K8sCredentials

	for _, env := range bc.Env {
		if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
			secretName := env.ValueFrom.SecretKeyRef.Name
			if !containsString(secretNames, secretName) && secretName != "" {
				secretNames = append(secretNames, secretName)
			}
		}
	}

	for _, vol := range bc.Volumes {
		if vol.Secret != nil {
			secretName := vol.Secret.SecretName
			if !containsString(secretNames, secretName) && secretName != "" {
				secretNames = append(secretNames, secretName)
			}
		}
	}

	return secretNames
}

// getSshSecretName returns the name of the secret that is used in the namespace for the SSH key to clone the repository.
func getSshSecretName(pipeline pipemanagerv1alpha1.PipelineSpec) string {
	return pipeline.SshSecretName
}
