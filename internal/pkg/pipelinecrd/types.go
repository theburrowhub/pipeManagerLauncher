package pipelinecrd

import (
	"encoding/json"
	"log/slog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CloneRepositoryOptions define the options for the clone repository step in the pipeline
type CloneRepositoryOptions struct {
	Cache     bool `json:"cache,omitempty"`
	Artifacts bool `json:"artifacts,omitempty"`
}

// CloneRepositoryConfig define the configuration for the clone repository step in the pipeline
type CloneRepositoryConfig struct {
	Enable  bool                   `json:"enable"`
	Options CloneRepositoryOptions `json:"options,omitempty"`
}

// PipelineSpec is a struct to store the pipeline data
type PipelineSpec struct {
	Name            string
	Description     string                `json:"description,omitempty"`
	Namespace       Namespace             `json:"namespace,omitempty"`
	CloneRepository CloneRepositoryConfig `json:"cloneRepository,omitempty"`
	CloneDepth      int                   `json:"cloneDepth,omitempty"`
	SshSecretName   string                `json:"sshSecretName,omitempty"`
	Launch          Launch                `json:"launch,omitempty"`
	Params          map[string]string     `json:"params"`
	Workspace       interface{}           `json:"workspaceDir,omitempty"`
	Tasks           map[string]Task       `json:"tasks"`
	FinishTasks     FinishTasks           `json:"finishTasks,omitempty"`
}

// FinishTasks is a struct to store the finish tasks data
type FinishTasks struct {
	Fail    map[string]Task `json:"fail,omitempty"`
	Success map[string]Task `json:"success,omitempty"`
}

// Task is a struct to store the task data
type Task struct {
	Description     string                `json:"description"`
	RunAfter        []string              `json:"runAfter,omitempty"`
	Batch           map[string]Batch      `json:"batch,omitempty"`
	Params          map[string]string     `json:"params,omitempty"`
	CloneRepository CloneRepositoryConfig `json:"cloneRepository,omitempty"`
	Paths           Paths                 `json:"paths,omitempty"`
	CloneDepth      int                   `json:"cloneDepth,omitempty"`
	Steps           []Step                `json:"steps"`
	Sidecars        []interface{}         `json:"sidecars,omitempty"`
	Volumes         []interface{}         `json:"volumes,omitempty"`
}

// Paths define the paths for the artifacts and cache in the task
type Paths struct {
	Artifacts []string `json:"artifacts,omitempty"`
	Cache     []string `json:"cache,omitempty"`
}

// Step is a struct to store the step data
type Step struct {
	Name         string        `json:"name"`
	Image        string        `json:"image"`
	Description  string        `json:"description,omitempty"`
	Env          []interface{} `json:"env,omitempty"`
	VolumeMounts []interface{} `json:"volumeMounts,omitempty"`
	Command      []string      `json:"command,omitempty"`
	Args         []string      `json:"args,omitempty"`
	Script       string        `json:"script,omitempty"`
}

// Batch define the batch data like params for the batch task
type Batch map[string]string

// Launch is a struct to store the launch data in the pipeline if it is defined
// It contains the steps to launch the next pipeline in the chain when the current pipeline finishes
// successfully or fails
type Launch struct {
	WhenFail    []string `json:"whenFail,omitempty"`
	WhenSuccess []string `json:"whenSuccess,omitempty"`
}

// Namespace is a struct to store the namespace data in the pipeline
// It contains the name of the namespace and the labels to apply to the namespace, if any. Also, a boolean to
// indicate if the namespace should be created or not
type Namespace struct {
	Create bool              `json:"create,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
	Name   string            `json:"name"`
}

// Pipeline is the Schema for the API
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec `json:"spec,omitempty"`
	Status interface{}  `json:"status,omitempty"`
}

// PipelineList contains a list of PipelineMain
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
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

// DeepCopyObject method of PipelineSpec that creates a deep copy of the pipeline spec
func (in *Pipeline) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Pipeline)
	in.DeepCopyInto(&out.ObjectMeta)
	return out
}

// DeepCopyObject method of PipelineList that creates a deep copy of the pipeline list
func (in *PipelineList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(PipelineList)
	in.DeepCopyInto(&out.ListMeta)
	return out
}
