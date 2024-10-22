package pipeobject

import (
	"encoding/json"
	"log/slog"
)

// CloneRepositoryOptions define the options for the clone repository step in the pipeline
type CloneRepositoryOptions struct {
	Cache     bool `yaml:"cache,omitempty"`
	Artifacts bool `yaml:"artifacts,omitempty"`
}

// CloneRepositoryConfig define the configuration for the clone repository step in the pipeline
type CloneRepositoryConfig struct {
	Enable  bool                   `yaml:"enable"`
	Options CloneRepositoryOptions `yaml:"options,omitempty"`
}

// Pipeline is a struct to store the pipeline data
type Pipeline struct {
	Name            string
	Description     string                `yaml:"description,omitempty"`
	Namespace       Namespace             `yaml:"namespace"`
	CloneRepository CloneRepositoryConfig `yaml:"cloneRepository,omitempty"`
	CloneDepth      int                   `yaml:"cloneDepth,omitempty"`
	SshSecretName   string                `yaml:"sshSecretName,omitempty"`
	Launch          Launch                `yaml:"launch,omitempty"`
	Params          map[string]string     `yaml:"params"`
	Workspace       interface{}           `yaml:"workspaceDir,omitempty"`
	Tasks           map[string]Task       `yaml:"tasks"`
	FinishTasks     FinishTasks           `yaml:"finishTasks,omitempty"`
}

// FinishTasks is a struct to store the finish tasks data
type FinishTasks struct {
	Fail    map[string]Task `yaml:"fail,omitempty"`
	Success map[string]Task `yaml:"success,omitempty"`
}

// Task is a struct to store the task data
type Task struct {
	Description     string                `yaml:"description"`
	RunAfter        []string              `yaml:"runAfter,omitempty"`
	Batch           map[string]Batch      `yaml:"batch,omitempty"`
	Params          map[string]string     `yaml:"params,omitempty"`
	CloneRepository CloneRepositoryConfig `yaml:"cloneRepository,omitempty"`
	Paths           Paths                 `yaml:"paths,omitempty"`
	CloneDepth      int                   `yaml:"cloneDepth,omitempty"`
	Steps           []Step                `yaml:"steps"`
	Sidecars        []interface{}         `yaml:"sidecars,omitempty"`
	Volumes         []interface{}         `yaml:"volumes,omitempty"`
}

// Paths define the paths for the artifacts and cache in the task
type Paths struct {
	Artifacts []string `yaml:"artifacts,omitempty"`
	Cache     []string `yaml:"cache,omitempty"`
}

// Step is a struct to store the step data
type Step struct {
	Name         string        `yaml:"name"`
	Image        string        `yaml:"image"`
	Description  string        `yaml:"description,omitempty"`
	Env          []interface{} `yaml:"env,omitempty"`
	VolumeMounts []interface{} `yaml:"volumeMounts,omitempty"`
	Command      []string      `yaml:"command,omitempty"`
	Args         []string      `yaml:"args,omitempty"`
	Script       string        `yaml:"script,omitempty"`
}

// Batch define the batch data like params for the batch task
type Batch map[string]string

// Launch is a struct to store the launch data in the pipeline if it is defined
// It contains the steps to launch the next pipeline in the chain when the current pipeline finishes
// successfully or fails
type Launch struct {
	WhenFail    []string `yaml:"whenFail,omitempty"`
	WhenSuccess []string `yaml:"whenSuccess,omitempty"`
}

// Namespace is a struct to store the namespace data in the pipeline
// It contains the name of the namespace and the labels to apply to the namespace, if any. Also, a boolean to
// indicate if the namespace should be created or not
type Namespace struct {
	Create bool              `yaml:"create,omitempty"`
	Labels map[string]string `yaml:"labels,omitempty"`
	Name   string            `yaml:"name"`
}

// DeepCopy method of Task that creates a deep copy of the task
func (t *Task) DeepCopy() Task {
	var newTask Task
	data, err := json.Marshal(t)
	if err != nil {
		slog.Error("Error copying task", "error", err)
	}
	err = json.Unmarshal(data, &newTask)
	if err != nil {
		slog.Error("Error copying task", "error", err)
	}
	return newTask
}
