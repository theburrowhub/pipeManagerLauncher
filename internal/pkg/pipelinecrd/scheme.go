package pipelinecrd

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

var (
	Scheme        = runtime.NewScheme()
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(schema.GroupVersion{
		Group:   "pipe-manager.org",
		Version: "v1alpha1",
	},
		&Pipeline{},
		&PipelineList{},
	)
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{
		Group:   "pipe-manager.org",
		Version: "v1alpha1",
	})
	return nil
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(runtime.NewScheme()))
	utilruntime.Must(AddToScheme(Scheme))
}
