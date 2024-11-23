package deploy

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pipemanagerv1alpha1 "github.com/sergiotejon/pipeManagerController/api/v1alpha1"

	"github.com/sergiotejon/pipeManagerLauncher/internal/pkg/k8s"
	"github.com/sergiotejon/pipeManagerLauncher/pkg/envvars"
)

const (
	Kind       = "Pipeline"
	APIVersion = "pipemanager.sergiotejon.github.io/v1alpha1"
)

// Pipeline deploys a pipeline object to the Kubernetes cluster
func Pipeline(name, namespace string, spec pipemanagerv1alpha1.PipelineSpec) (string, string, error) {
	spec.Name = name

	// Check if the params map is not nil and create it if it is
	if spec.Params == nil {
		spec.Params = make(map[string]string)
	}
	// Add the extra params to the params map
	spec.Params["COMMIT"] = envvars.Variables["COMMIT"]
	spec.Params["REPOSITORY"] = envvars.Variables["REPOSITORY"]

	// Generate the pipeline object
	pipeline := generatePipelineObject(name, namespace, spec)

	// Deploy the pipeline object to the Kubernetes cluster
	err := deployPipelineObject(pipeline)
	if err != nil {
		return "", "", err
	}

	resourceName := pipeline.Name
	resourceNamespace := pipeline.Namespace

	return resourceName, resourceNamespace, nil
}

// generatePipelineObject generates a pipeline object for the given name, namespace and spec to use in the deployment process
func generatePipelineObject(name, namespace string, spec pipemanagerv1alpha1.PipelineSpec) *pipemanagerv1alpha1.Pipeline {
	return &pipemanagerv1alpha1.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:       Kind,
			APIVersion: APIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%.58s-%s", name, generateRandomString(63-len(name)-1)),
			Namespace: namespace,
		},
		Spec: spec,
	}
}

// deployPipelineObject deploys the given pipeline object to the Kubernetes cluster
func deployPipelineObject(pipeline *pipemanagerv1alpha1.Pipeline) error {
	config, err := k8s.GetKubernetesConfig()
	if err != nil {
		return err
	}

	k8sClient, err := client.New(config, client.Options{Scheme: pipemanagerv1alpha1.Scheme})
	if err != nil {
		return err
	}

	err = k8sClient.Create(context.Background(), pipeline)
	if err != nil {
		return err
	}

	return nil
}

// generateRandomString generates a random string of the given length and returns its SHA-256 hash
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	randomString := string(b)
	hash := sha256.Sum256([]byte(randomString))
	return fmt.Sprintf("%x", hash)
}
