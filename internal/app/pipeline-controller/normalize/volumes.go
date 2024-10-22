package normalize

import "github.com/sergiotejon/pipeManager/internal/pkg/pipelinecrd"

const workspaceDir = "/workspaceDir" // Default workspace directory for the all steps

// addDefaultVolumes adds the default volumes to the task
func addDefaultVolumes(task pipelinecrd.Task, workspace interface{}, sshSecretName string) pipelinecrd.Task {
	// Volumes for the workspaceDir and the ssh secret if it is defined
	var volumes []interface{}
	volumes = append(volumes, workspaceVolume(workspace))
	if sshSecretName != "" {
		volumes = append(volumes, sshSecretVolume(sshSecretName))
	}

	// Add the volumes to the task
	task.Volumes = append(task.Volumes, volumes...)

	return task
}

// workspaceVolume defines the volume for the workspaceDir
func workspaceVolume(w interface{}) interface{} {
	if w != nil {
		return w
	} else {
		return map[string]interface{}{
			"name":     "workspace",
			"emptyDir": map[string]interface{}{},
		}
	}
}

// sshSecretVolume defines the volume for the ssh secret
func sshSecretVolume(sshSecretName string) interface{} {
	return map[string]interface{}{
		"name": "ssh-credentials",
		"secret": map[string]interface{}{
			"secretName":  sshSecretName,
			"defaultMode": 256,
		},
	}
}

// addDefaultVolumeMounts adds the default volume mounts to the step
func addDefaultVolumeMounts(step pipelinecrd.Step, workspaceDir, sshSecretName string) pipelinecrd.Step {
	// Volume mounts for the workspaceDir and the ssh secret if it is defined
	var volumeMounts []interface{}
	volumeMounts = append(volumeMounts, workspaceVolumeMount(workspaceDir))
	if sshSecretName != "" {
		volumeMounts = append(volumeMounts, sshSecretVolumeMount())
	}

	// Add the volumes to the step
	step.VolumeMounts = append(step.VolumeMounts, volumeMounts...)

	return step
}

// workspaceVolumeMount defines the volume mount for the workspaceDir
func workspaceVolumeMount(mountPath string) map[string]interface{} {
	return map[string]interface{}{
		"name":      "workspace",
		"mountPath": mountPath,
	}
}

// sshSecretVolumeMount defines the volume mount for the ssh secret
func sshSecretVolumeMount() map[string]interface{} {
	return map[string]interface{}{
		"name":      "ssh-credentials",
		"mountPath": "/root/.ssh",
		"readOnly":  true,
	}
}
