package normalize

import (
	"fmt"
	"strings"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/envvars"
)

const envVarPrefix = "PIPELINE_"

// defineLaunchPipelineTask defines the task to launch the next pipeline in the chain
func defineLaunchPipelineTask(currentPipeline Pipeline, repository, commit, pipelineToLaunch string) Task {
	var env []interface{}

	// Set the parameters for the new pipeline through environment variables
	for name, value := range envvars.Variables {
		env = append(env, map[string]interface{}{
			"name":  envVarPrefix + name,
			"value": value,
		})
	}

	// Set environment variable for configuration instead of using the configuration file
	env = append(env, map[string]interface{}{
		"name":  envVarPrefix + "NAME",
		"value": pipelineToLaunch,
	})
	env = append(env, map[string]interface{}{
		"name":  "COMMON_DATA_LOG_LEVEL",
		"value": config.Common.Data.Log.Level,
	})
	env = append(env, map[string]interface{}{
		"name":  "COMMON_DATA_LOG_FORMAT",
		"value": config.Common.Data.Log.Format,
	})
	env = append(env, map[string]interface{}{
		"name":  "COMMON_DATA_LOG_FILE",
		"value": config.Common.Data.Log.File,
	})
	env = append(env, map[string]interface{}{
		"name":  "LAUNCHER_DATA_CLONEDEPTH",
		"value": fmt.Sprintf("%d", config.Launcher.Data.CloneDepth),
	})
	env = append(env, map[string]interface{}{
		"name":  "LAUNCHER_DATA_NAMESPACE",
		"value": config.Launcher.Data.Namespace,
	})

	// Get the image from the configuration to launch the pipeline
	launcherImage := config.Launcher.Data.GetLauncherImage()

	// Define the task to launch the pipeline
	launchPipelineTask := Task{
		Description: fmt.Sprintf("Launch the next pipeline '%s' in the chain", pipelineToLaunch),
		Steps: []Step{
			{
				Name:        "launch-pipeline",
				Description: fmt.Sprintf("Launch the next pipeline '%s' in the chain", pipelineToLaunch),
				Image:       launcherImage,
				Env:         env,
				Script: strings.Join([]string{
					fmt.Sprintf("#!%s", defaultShell),
					defaultShellSets,
					fmt.Sprintf("%s %s", launcherBinary, "run"),
				}, "\n"),
			},
		},
	}

	// Append the clone repository step to the task steps to ensure the repository is cloned before launching the pipeline.
	// It's the only way to look for the pipeline definition of pipelineToLaunch in the repository.
	cloneRepositoryStep := defineCloneRepoStep(launchPipelineTask, repository, commit)
	launchPipelineTask.Steps = append([]Step{cloneRepositoryStep}, launchPipelineTask.Steps...)

	return launchPipelineTask
}
