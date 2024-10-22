// Package normalize provides functionality to normalize the pipeline data.
// It adds the necessary steps to the tasks to:
// - clone the repository
// - download and upload the artifacts and cache
// - expand each batch task in the pipeline
// It also adds the necessary finish tasks to:
// - launch the next pipeline in the chain
package normalize

import (
	"gopkg.in/yaml.v3"

	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
	"github.com/sergiotejon/pipeManager/internal/pkg/pipeobject"
)

const (
	defaultShell     = "/bin/sh"       // Default shell for the automated steps
	defaultShellSets = "set -e"        // Default shell sets for the automated steps
	launcherBinary   = "/app/launcher" // Default launcher binary path
)

// Normalize normalizes the pipeline data
// It adds the necessary steps to the tasks to:
// - clone the repository
// - download and upload the artifacts and cache
// - expand each batch task in the pipeline
// It also adds the necessary finish tasks to:
// - launch the next pipeline in the chain
// TODO: Retrieve just one pipeline object to normalize and launch a Tekton pipeline
func Normalize(data map[string]interface{}) ([]pipeobject.Pipeline, error) {
	var rawPipelines []pipeobject.Pipeline
	var err error

	// Convert pipeline raw data to Pipelines struct
	rawPipelines, err = convertToPipelines(data)
	if err != nil {
		return nil, err
	}

	// Loop rawPipelines
	for item := range rawPipelines {
		rawPipeline := rawPipelines[item]

		repository := rawPipeline.Params["REPOSITORY"]
		commit := rawPipeline.Params["COMMIT"]

		// Loop tasks for:
		// - cloning the repository
		// - download and upload the artifacts and cache
		// - expand each batch task in the rawPipeline
		for taskName, taskData := range rawPipeline.Tasks {
			// Process the task data to add the necessary steps
			taskData = processTask(rawPipeline, taskName, taskData, repository, commit)

			// Explode batch tasks if they are defined
			if taskData.Batch != nil {
				// Explode the batch task
				tasks := processBatchTask(taskName, taskData)
				if tasks != nil {
					// Add all the new tasks to the rawPipeline
					for name, tData := range tasks {
						rawPipeline.Tasks[name] = tData
					}
					// Remove the original task from the rawPipeline
					delete(rawPipeline.Tasks, taskName)
				}
			} else {
				// Replace the task with the new data in the rawPipeline
				rawPipeline.Tasks[taskName] = taskData
			}
		}

		// Normalize Fail tasks
		for taskName, taskData := range rawPipeline.FinishTasks.Fail {
			// Process the task data to add the necessary steps
			taskData = processTask(rawPipeline, taskName, taskData, repository, commit)

			// Explode batch tasks if they are defined
			if taskData.Batch != nil {
				// Explode the batch task
				tasks := processBatchTask(taskName, taskData)
				if tasks != nil {
					// Add all the new tasks to the rawPipeline
					for name, tData := range tasks {
						rawPipeline.FinishTasks.Fail[name] = tData
					}
					// Remove the original task from the rawPipeline
					delete(rawPipeline.Tasks, taskName)
				}
			} else {
				// Replace the task with the new data in the rawPipeline
				rawPipeline.FinishTasks.Fail[taskName] = taskData
			}

		}

		// Loop through the list of rawPipelines to launch when the current rawPipeline fails
		for _, launchPipelineName := range rawPipeline.Launch.WhenFail {
			taskName := k8sObjectName("launch", launchPipelineName)
			rawPipeline.FinishTasks.Fail[taskName] = defineLaunchPipelineTask(rawPipeline, repository, commit, launchPipelineName)
		}

		// Normalize Success tasks
		for taskName, taskData := range rawPipeline.FinishTasks.Success {
			// Process the task data to add the necessary steps
			taskData = processTask(rawPipeline, taskName, taskData, repository, commit)

			// Explode batch tasks if they are defined
			if taskData.Batch != nil {
				// Explode the batch task
				tasks := processBatchTask(taskName, taskData)
				if tasks != nil {
					// Add all the new tasks to the rawPipeline
					for name, tData := range tasks {
						rawPipeline.FinishTasks.Success[name] = tData
					}
					// Remove the original task from the rawPipeline
					delete(rawPipeline.Tasks, taskName)
				}
			} else {
				// Replace the task with the new data in the rawPipeline
				rawPipeline.FinishTasks.Success[taskName] = taskData
			}

		}

		// Loop through the list of rawPipelines to launch when the current rawPipeline finishes successfully
		for _, launchPipelineName := range rawPipeline.Launch.WhenSuccess {
			taskName := k8sObjectName("launch", launchPipelineName)
			rawPipeline.FinishTasks.Success[taskName] = defineLaunchPipelineTask(rawPipeline, repository, commit, launchPipelineName)
		}

		// Clean
		// Remove unnecessary cloneRepository and launchPipeline from the rawPipeline
		rawPipelines[item].CloneRepository = pipeobject.CloneRepositoryConfig{}
		rawPipelines[item].Launch = pipeobject.Launch{}
	}

	return rawPipelines, nil
}

// convertToPipelines converts the raw data to a list of Pipeline structs
func convertToPipelines(data map[string]interface{}) ([]pipeobject.Pipeline, error) {
	var pipelines []pipeobject.Pipeline

	for name, item := range data {
		// Convert each item to YAML
		yamlData, err := yaml.Marshal(item)
		if err != nil {
			return nil, err
		}

		// Unmarshal YAML to Pipeline struct
		var pipeline pipeobject.Pipeline
		err = yaml.Unmarshal(yamlData, &pipeline)
		if err != nil {
			return nil, err
		}

		// Set the name of the pipeline
		pipeline.Name = name

		// Add the pipeline to the list
		pipelines = append(pipelines, pipeline)
	}

	return pipelines, nil
}

// processTask processes the task data to add the necessary steps to:
// - clone the repository
// - download and upload the artifacts and cache
// - expand each batch task in the pipeline
func processTask(pipe pipeobject.Pipeline, taskName string, taskData pipeobject.Task, repository, commit string) pipeobject.Task {
	logging.Logger.Debug("Normalizing task", "taskName", taskName)

	// This is the list of steps that will be added to the task at the beginning
	var firstSteps []pipeobject.Step

	// Add the clone repository step if it is defined as true in the pipe or in the task itself
	cloneRepository := pipe.CloneRepository.Enable || taskData.CloneRepository.Enable
	if cloneRepository {
		logging.Logger.Debug("Adding clone repository step", "taskName", taskName)
		cloneRepositoryStep := defineCloneRepoStep(taskData, repository, commit)
		firstSteps = append(firstSteps, cloneRepositoryStep)

		// Download artifacts if it is defined as true in the pipe or in the task itself and the clone repository step is enabled
		artifacts := pipe.CloneRepository.Options.Artifacts || taskData.CloneRepository.Options.Artifacts
		if artifacts {
			logging.Logger.Debug("Adding download artifacts step", "taskName", taskName)
			downloadArtifactsStep := defineDownloadArtifactsStep(taskData)
			firstSteps = append(firstSteps, downloadArtifactsStep)
		}

		// idem for cache
		caches := pipe.CloneRepository.Options.Cache || taskData.CloneRepository.Options.Cache
		if caches {
			logging.Logger.Debug("Adding download cache step", "taskName", taskName)
			downloadCacheStep := defineDownloadCacheStep(taskData)
			firstSteps = append(firstSteps, downloadCacheStep)
		}
	}

	// Add all these automatic steps at the beginning of the task
	taskData.Steps = append(firstSteps, taskData.Steps...)

	// This is the list of steps that will be added at the end of the task
	var lastSteps []pipeobject.Step

	// If the clone repository step is enabled, upload the artifacts and cache
	if cloneRepository {
		// Upload artifacts if it is defined as true in the pipe or in the task itself and the clone repository step is enabled
		artifacts := pipe.CloneRepository.Options.Artifacts || taskData.CloneRepository.Options.Artifacts
		if artifacts {
			logging.Logger.Debug("Adding upload artifacts step", "taskName", taskName)
			uploadArtifactsStep := defineUploadArtifactsStep(taskData)
			lastSteps = append(lastSteps, uploadArtifactsStep)
		}

		// idem for cache
		caches := pipe.CloneRepository.Options.Cache || taskData.CloneRepository.Options.Cache
		if caches {
			logging.Logger.Debug("Adding upload cache step", "taskName", taskName)
			uploadCacheStep := defineUploadCacheStep(taskData)
			lastSteps = append(lastSteps, uploadCacheStep)
		}
	}

	// Add all the upload steps at the end of the task
	taskData.Steps = append(taskData.Steps, lastSteps...)

	// Add the default volumes for the workspace and the ssh credentials secret if it is defined to the task
	taskData = addDefaultVolumes(taskData, pipe.Workspace, pipe.SshSecretName)

	// Add the volumeMounts for the workspaceDir and the ssh secret if it is defined to the steps
	for i := range taskData.Steps {
		taskData.Steps[i] = addDefaultVolumeMounts(taskData.Steps[i], workspaceDir, pipe.SshSecretName)
	}

	// Clean
	// Remove unnecessary cloneRepository and path from the task
	taskData.CloneRepository = pipeobject.CloneRepositoryConfig{}
	taskData.Paths = pipeobject.Paths{}

	return taskData
}
