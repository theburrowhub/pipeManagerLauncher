package deploy

import (
	"fmt"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sergiotejon/pipeManager/internal/pkg/pipelinecrd"
)

// WIP: Pipeline deploys a pipeline object to the Kubernetes cluster
func Pipeline(name, namespace string, spec pipelinecrd.PipelineSpec) error {
	// Generate the pipeline object
	pipeline := generatePipelineObject(name, namespace, spec)

	// Create the pipeline object
	// TODO

	data, err := yaml.Marshal(pipeline)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

// generatePipelineObject generates a pipeline object for the given name, namespace and spec to use in the deployment process
func generatePipelineObject(name, namespace string, spec pipelinecrd.PipelineSpec) *pipelinecrd.Pipeline {
	return &pipelinecrd.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pipeline",
			APIVersion: "pipe-manager.org/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: spec,
	}
}
