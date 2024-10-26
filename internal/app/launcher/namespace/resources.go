package namespace

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
)

const (
	applicationLabelKey          = "app.kubernetes.io/name"
	applicationManagedByLabelKey = "app.kubernetes.io/managed-by"
	applicationLabelValue        = "pipe-manager"
)

// checkIfResourceNamespaceExists checks if a namespace with the given name exists
func checkIfResourceNamespaceExists(client *kubernetes.Clientset, name string) (bool, error) {
	_, err := client.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false, nil
	}

	return true, nil
}

// createResourceNamespace creates a namespace with the given name and labels
func createResourceNamespace(client *kubernetes.Clientset, name string, labels map[string]string) error {
	// Add the default labels
	customLabels := map[string]string{
		applicationLabelKey:          applicationLabelValue,
		applicationManagedByLabelKey: applicationLabelValue,
	}
	for k, v := range labels {
		customLabels[k] = v
	}

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: customLabels,
		},
	}

	// Create the namespace
	_, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	return nil
}

// updateResourceNamespaceLabels updates the labels of a namespace with the given name if they are different
func updateResourceNamespaceLabels(client *kubernetes.Clientset, name string, labels map[string]string) error {
	// Add the default labels
	customLabels := map[string]string{
		applicationLabelKey:          applicationLabelValue,
		applicationManagedByLabelKey: applicationLabelValue,
	}
	for k, v := range labels {
		customLabels[k] = v
	}

	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get namespace: %w", err)
	}

	// Check if the labels are different
	if !mapsOfStringAreDifferent(ns.Labels, labels) {
		return nil
	}

	// Update the labels
	ns.Labels = customLabels
	_, err = client.CoreV1().Namespaces().Update(context.TODO(), ns, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update namespace labels: %w", err)
	}

	return nil
}

// createOrUpdateServiceAccount creates or updates a service account with the given name and namespace
func createOrUpdateServiceAccount(client *kubernetes.Clientset, saName string, namespace string) error {
	// Check if the service account already exists
	sa, err := client.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), saName, metav1.GetOptions{})
	if err == nil { // Service account exists, update it
		sa.ObjectMeta.Name = saName
		sa.ObjectMeta.Namespace = namespace
		_, err = client.CoreV1().ServiceAccounts(namespace).Update(context.TODO(), sa, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update service account: %w", err)
		}
	} else { // it does not exist, create it
		sa = &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      saName,
				Namespace: namespace,
			},
		}
		_, err = client.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), sa, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create service account: %w", err)
		}
	}

	// Create roleBinding to bind the given roles the service account
	err = roleBinding(client, namespace, saName, config.Launcher.Data.RolesBinding)
	if err != nil {
		return fmt.Errorf("failed to create roles and bind to service account: %w", err)
	}

	return nil
}

// GetSecretsContent retrieves the content of each secret given by a list and returns it in a key/value list
// where the key is the name of the secret and the value is its content.
func GetSecretsContent(client *kubernetes.Clientset, namespace string, secretNames []string) (map[string]map[string][]byte, error) {
	secretsContent := make(map[string]map[string][]byte)

	for _, secretName := range secretNames {
		secret, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to get secret %s: %w", secretName, err)
		}
		secretsContent[secretName] = secret.Data
	}

	return secretsContent, nil
}

// CopySecretsToNamespace copies the key/value pairs of all retrieved secrets to another given namespace
func CopySecretsToNamespace(client *kubernetes.Clientset, sourceNamespace string, targetNamespace string, secretNames []string) error {
	for _, secretName := range secretNames {
		// Retrieve the secret from the source namespace
		secret, err := client.CoreV1().Secrets(sourceNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get secret %s from namespace %s: %w", secretName, sourceNamespace, err)
		}

		// Create a new secret object for the target namespace
		newSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secret.Name,
				Namespace: targetNamespace,
			},
			Data: secret.Data,
			Type: secret.Type,
		}

		// Check if the secret already exists in the target namespace
		_, err = client.CoreV1().Secrets(targetNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err == nil { // Secret exists, update it
			_, err = client.CoreV1().Secrets(targetNamespace).Update(context.TODO(), newSecret, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("failed to update secret %s in namespace %s: %w", secretName, targetNamespace, err)
			}
		} else { // Secret does not exist, create it
			_, err = client.CoreV1().Secrets(targetNamespace).Create(context.TODO(), newSecret, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create secret %s in namespace %s: %w", secretName, targetNamespace, err)
			}
		}
	}

	return nil
}
