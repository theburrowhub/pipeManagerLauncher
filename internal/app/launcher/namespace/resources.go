package namespace

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
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
	// TODO:
	// Add specific labels to the namespace for the application (see webhook-listener for an example)
	// Send that labels to common packaged for reuse
	// Like helm charts with its labels

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
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
	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get namespace: %w", err)
	}

	// Check if the labels are different
	if !mapsOfStringAreDifferent(ns.Labels, labels) {
		return nil
	}

	// Update the labels
	ns.Labels = labels
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
