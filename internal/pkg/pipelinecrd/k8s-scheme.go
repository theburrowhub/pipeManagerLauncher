package pipelinecrd

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// PipelineStatus defines the observed state of PipelineSpec
type PipelineStatus struct {
	// Aquí puedes definir el estado observado
	// Por ejemplo, condiciones, mensajes, etc.
}

// Pipeline is the Schema for the API
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec,omitempty"`
	Status PipelineStatus `json:"status,omitempty"`
}

// PipelineList contains a list of PipelineMain
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

func init() {
	// TODO: Ver información en chatgpt para implementar esto
	scheme := runtime.NewScheme()
	err := addToScheme(scheme)
	if err != nil {
		panic(err)
	}
}

func addToScheme(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(schema.GroupVersion{
		Group:   "pipe-manager.sergiotejon.github.com",
		Version: "v1beta1",
	},
		&Pipeline{},
		&PipelineList{},
	)
	return nil
}

func (in *Pipeline) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Pipeline)
	in.DeepCopyInto(&out.ObjectMeta)
	return out
}

func (in *PipelineList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(PipelineList)
	in.DeepCopyInto(&out.ListMeta)
	return out
}
