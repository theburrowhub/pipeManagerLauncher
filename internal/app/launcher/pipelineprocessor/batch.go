package pipelineprocessor

// processBatchTask processes a batch task and returns a map of tasks to run in parallel
func processBatchTask(taskName string, taskData Task) map[string]Task {
	tasks := make(map[string]Task)

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
