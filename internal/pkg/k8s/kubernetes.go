package k8s

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

// GetKubernetesConfig returns a Kubernetes client configuration for either in-cluster or local access
func GetKubernetesConfig() (*rest.Config, error) {
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

	return cfg, nil
}

// GetKubernetesClient returns a Kubernetes clientset configured for either in-cluster or local access
func GetKubernetesClient() (*kubernetes.Clientset, error) {
	cfg, err := GetKubernetesConfig()
	if err != nil {
		return nil, err
	}

	// Create a clientset for interacting with the Kubernetes API
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}
