package namespace

import (
	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/k8s"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
	"github.com/sergiotejon/pipeManager/internal/pkg/pipelinecrd"
)

const pipeManagerSA = "pipe-manager-sa"

// Create creates a namespace with the given name and labels and creates the necessary resources inside the namespace
// like the service account and the secrets for the bucket credentials.
func Create(ns pipelinecrd.Namespace) error {
	namespace := ns.Name
	labels := ns.Labels

	// Get the Kubernetes client
	client, err := k8s.GetKubernetesClient()
	if err != nil {
		return err
	}

	// Check if the namespace already namespaceAlreadyExists
	namespaceAlreadyExists, err := checkIfResourceNamespaceExists(client, namespace)
	if err != nil {
		return err
	}

	// Create the namespace if it does not exist or update the labels if they are different
	if !namespaceAlreadyExists {
		logging.Logger.Info("Creating namespace", "namespace", namespace)
		err := createResourceNamespace(client, namespace, labels)
		if err != nil {
			return err
		}
	} else {
		logging.Logger.Info("Updating namespace labels", "namespace", namespace)
		err := updateResourceNamespaceLabels(client, namespace, labels)
		if err != nil {
			return err
		}
	}

	// Create or update the service account
	logging.Logger.Info("Creating or updating service account", "namespace", pipeManagerSA)
	err = createOrUpdateServiceAccount(client, pipeManagerSA, namespace)
	if err != nil {
		return err
	}

	// Retrieve secrets from config for bucket credentials and copy them to the namespace, update if they already exist
	logging.Logger.Info("Retrieving bucket credentials secret from config")
	err = CopySecretsToNamespace(client,
		config.Launcher.Data.Namespace,
		namespace,
		getBucketCredentialsSecretFromConfig())
	if err != nil {
		return err
	}

	return nil
}
