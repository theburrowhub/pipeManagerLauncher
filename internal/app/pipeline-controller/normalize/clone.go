package normalize

import (
	"fmt"
	"strings"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/pipelinecrd"
)

// defineCloneRepoStep defines the clone repository step in the task
func defineCloneRepoStep(taskData pipelinecrd.Task, repository, commit string) pipelinecrd.Step {
	// Get the git image from the configuration
	gitImage := config.Launcher.Data.GetLauncherImage()

	// Set the clone depth
	var cloneDepth int
	if taskData.CloneDepth == 0 {
		cloneDepth = config.Launcher.Data.CloneDepth
	} else {
		cloneDepth = taskData.CloneDepth
	}

	// Repository step for cloning the repository
	command := fmt.Sprintf("%s %s --depth %d --repository '%s' --commit '%s' --destination '%s'",
		launcherBinary, "clone", cloneDepth, repository, commit, workspaceDir)
	step := pipelinecrd.Step{
		Name:        "clone-repository",
		Description: "Automatically clone the repository",
		Image:       gitImage,
		Script: strings.Join([]string{
			fmt.Sprintf("#!%s", defaultShell),
			defaultShellSets,
			command,
		}, "\n"),
	}

	return step
}
