package namespace

import (
	pipemanagerv1alpha1 "github.com/sergiotejon/pipeManagerController/api/v1alpha1"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/k8s"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

const pipeManagerSA = "pipe-manager-sa"

// Create creates a namespace with the given name and labels and creates the necessary resources inside the namespace
// like the service account and the secrets for the bucket credentials.
func Create(pipeline pipemanagerv1alpha1.PipelineSpec) error {
	ns := pipeline.Namespace
	namespaceName := ns.Name
	labels := ns.Labels

	client, err := k8s.GetKubernetesClient()
	if err != nil {
		return err
	}

	// Check if the namespaceName already namespaceAlreadyExists
	namespaceAlreadyExists, err := checkIfResourceNamespaceExists(client, namespaceName)
	if err != nil {
		return err
	}

	// Create the namespaceName if it does not exist or update the labels if they are different
	if !namespaceAlreadyExists {
		logging.Logger.Info("Creating namespaceName", "namespaceName", namespaceName)
		err := createResourceNamespace(client, namespaceName, labels)
		if err != nil {
			return err
		}
	} else {
		logging.Logger.Info("Updating namespaceName labels", "namespaceName", namespaceName)
		err := updateResourceNamespaceLabels(client, namespaceName, labels)
		if err != nil {
			return err
		}
	}

	// Create or update the service account
	logging.Logger.Info("Creating or updating service account", "namespaceName", pipeManagerSA)
	err = createOrUpdateServiceAccount(client, pipeManagerSA, namespaceName)
	if err != nil {
		return err
	}

	// Retrieve secrets from config for bucket credentials and the SSH secret name from the pipeline, and copy them to the namespace, updating if they already exist
	logging.Logger.Info("Retrieving bucket credentials secret from config")
	secrets := append(getBucketCredentialsSecretFromConfig(), getSshSecretName(pipeline))
	err = CopySecretsToNamespace(client,
		config.Launcher.Data.Namespace,
		namespaceName,
		secrets,
	)
	if err != nil {
		return err
	}

	return nil
}
