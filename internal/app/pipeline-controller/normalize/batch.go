package normalize

import "github.com/sergiotejon/pipeManager/internal/pkg/pipeobject"

// processBatchTask processes a batch task and returns a map of tasks to run in parallel
func processBatchTask(taskName string, taskData pipeobject.Task) map[string]pipeobject.Task {
	tasks := make(map[string]pipeobject.Task)

	// If no batch is defined, return nil
	if taskData.Batch == nil {
		return nil
	}

	for batchName, batchParams := range taskData.Batch {
		// Add the new task to the tasks map
		name := k8sObjectName(taskName, batchName)

		// Copy taskData to newTask removing the batch field
		newTask := taskData.DeepCopy()
		newTask.Batch = nil
		tasks[name] = newTask

		// Copy the params from the batch task to the new task
		for key, value := range batchParams {
			tasks[name].Params[key] = value
		}
	}

	return tasks
}
