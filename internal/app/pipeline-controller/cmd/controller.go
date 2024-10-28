package cmd

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/sergiotejon/pipeManager/internal/pkg/pipelinecrd"
)

// PipelineReconcile reconciles a PushMain object
type PipelineReconcile struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is the main reconciliation loop for the controller
func (r *PipelineReconcile) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var pipeline pipelinecrd.Pipeline
	if err := r.Get(ctx, req.NamespacedName, &pipeline); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Pipeline resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get PushMain.")
		return ctrl.Result{}, err
	}

	// Implementa tu lógica de reconciliación aquí
	// Por ejemplo, crear namespaces, gestionar tareas, etc.

	logger.Info("Pipeline spec", "spec", pipeline.Spec)

	// TODO: Normalize
	fmt.Println("Normalize... Under construction.")
	// -- Read spec from k8s object (only one pipeline)
	// -- Refactor Normalize to work only with one pipeline
	// Normalize the pipelines
	//pipelines, err := normalize.Normalize(rawPipelines)
	//if err != nil {
	//	logging.Logger.Error("Error normalizing pipelines", "msg", err)
	//	os.Exit(ErrCodeNormalize)
	//}

	// Aquí puedes iterar sobre las tareas y gestionar cada una
	for taskName, task := range pipeline.Spec.Tasks {
		logger.Info("Processing Task", "TaskName", taskName, "Description", task.Description)
		for _, step := range task.Steps {
			logger.Info("  Step", "Name", step.Name, "Image", step.Image)
			// Implementa la lógica para manejar cada step
			// Por ejemplo, crear Pods o Jobs según los steps
		}
	}

	// Actualizar el estado si es necesario
	// pushMain.Status.SomeField = "SomeValue"
	// if err := r.Status().Update(ctx, &pushMain); err != nil {
	//     logger.Error(err, "Failed to update PushMain status")
	//     return ctrl.Result{}, err
	// }

	return ctrl.Result{RequeueAfter: 10 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *PipelineReconcile) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinecrd.Pipeline{}).
		Owns(&corev1.Namespace{}). // Si creas recursos que deseas observar
		Complete(r)

	//
}
