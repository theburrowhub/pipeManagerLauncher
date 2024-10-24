package namespace

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
