package deploy

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/sergiotejon/pipeManager/internal/pkg/k8s"
	"github.com/sergiotejon/pipeManager/internal/pkg/pipelinecrd"
)

// Pipeline deploys a pipeline object to the Kubernetes cluster
func Pipeline(name, namespace string, spec pipelinecrd.PipelineSpec) (string, string, error) {
	// Remove unused fields
	removeUnusedFields(&spec)

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

// removeUnusedFields removes the unused fields from the pipeline object
func removeUnusedFields(spec *pipelinecrd.PipelineSpec) {
	spec.Namespace = pipelinecrd.Namespace{}
}

// generatePipelineObject generates a pipeline object for the given name, namespace and spec to use in the deployment process
func generatePipelineObject(name, namespace string, spec pipelinecrd.PipelineSpec) *pipelinecrd.Pipeline {
	return &pipelinecrd.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pipeline",
			APIVersion: "pipe-manager.org/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%.58s-%s", name, generateRandomString(63-len(name)-1)),
			Namespace: namespace,
		},
		Spec: spec,
	}
}

// deployPipelineObject deploys the given pipeline object to the Kubernetes cluster
func deployPipelineObject(pipeline *pipelinecrd.Pipeline) error {
	config, err := k8s.GetKubernetesConfig()
	if err != nil {
		return err
	}

	k8sClient, err := client.New(config, client.Options{Scheme: pipelinecrd.Scheme})
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
