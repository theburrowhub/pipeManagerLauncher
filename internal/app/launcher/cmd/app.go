// Package cmd
// cmd package contains the Cobra commands for the launcher application.
package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"

	"github.com/sergiotejon/pipeManager/internal/app/launcher/deploy"
	"github.com/sergiotejon/pipeManager/internal/app/launcher/pipelineprocessor"
	"github.com/sergiotejon/pipeManager/internal/app/launcher/repository"
	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/envvars"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
	"github.com/sergiotejon/pipeManager/internal/pkg/pipelinecrd"
)

const (
	ErrCodeOK                 = 0
	ErrCodeLoadConfig         = 1
	ErrCodeCloneRepo          = 2
	ErrCodeMixFiles           = 3
	ErrCodeConvertingPipeline = 4
	ErrCodeBucketDownload     = 6
	ErrCodeBucketUpload       = 7
	ErrCodeDeploy             = 8
)

const (
	templateFolder = "/etc/pipe-manager/templates" // templateFolder is the folder where the templates are stored
	repoDir        = "/tmp/repo"                   // repoDir is the directory where the repository is cloned
)

var (
	configFile string // configFile is the path to the configuration file
)

// setup initializes the application by loading the configuration, setting up the logger, and getting environment variables
func setup() {
	var err error

	// Load configuration
	err = config.LoadLauncherConfig(configFile)
	if err != nil {
		log.Printf("Error loading launcher config: %v", err)
		os.Exit(ErrCodeLoadConfig)
	}
	err = config.LoadCommonConfig(configFile)
	if err != nil {
		log.Printf("Error loading common config: %v", err)
		os.Exit(ErrCodeLoadConfig)
	}

	// Setup Logger
	err = logging.SetupLogger(config.Common.Data.Log.Level, config.Common.Data.Log.Format, config.Common.Data.Log.File)
	if err != nil {
		log.Printf("Error configuring the logger: %v", err)
		os.Exit(ErrCodeLoadConfig)
	}

	logging.Logger.Info("Pipe Manager starting up...")
	logging.Logger.Debug("Setup", "configFile", configFile,
		"logLevel", config.Common.Data.Log.Level,
		"logFormat", config.Common.Data.Log.Format,
		"logFile", config.Common.Data.Log.File)

	// Get environment variables and log them
	envvars.GetEnvVars()
	for key, value := range envvars.Variables {
		logging.Logger.Debug("Environment variable", key, value)
	}
}

// app is the main application function
// It loads the configuration, sets up the logger and starts the launcher
func app() {
	var err error

	// Clone the repository
	err = repository.Clone(envvars.Variables["REPOSITORY"],
		config.Launcher.Data.CloneDepth,
		envvars.Variables["COMMIT"],
		repoDir)
	if err != nil {
		logging.Logger.Error("Error cloning repository", "msg", err,
			"repository", envvars.Variables["REPOSITORY"],
			"commit", envvars.Variables["COMMIT"],
			"depth", config.Launcher.Data.CloneDepth)
		os.Exit(ErrCodeCloneRepo)
	}

	logging.Logger.Info("Repository cloned successfully", "repository", envvars.Variables["REPOSITORY"], "commit", envvars.Variables["COMMIT"])

	// Mix all the pipeline files
	const pipelineDir = ".pipelines"
	pipelineFolder := filepath.Join(repoDir, pipelineDir)
	err, combinedData := pipelineprocessor.MixPipelineFiles(pipelineFolder)
	if err != nil {
		logging.Logger.Error("Error mixing pipeline files", "msg", err, "folder", pipelineFolder)
		os.Exit(ErrCodeMixFiles)
	}
	logging.Logger.Info("Pipeline files mixed successfully", "folder", pipelineFolder)
	for key, _ := range combinedData {
		if key == "global" {
			continue
		}
		logging.Logger.Debug("Pipeline found", "pipeline", key)
	}

	// Find the pipeline to launch
	var rawPipelines map[string]interface{}
	if envvars.Variables["NAME"] == "" { // If no pipeline name is provided, launch all pipelines that match the triggers
		logging.Logger.Info("Looking for pipelines using triggers")
		rawPipelines = pipelineprocessor.FindPipelineByRegex(combinedData, envvars.Variables)
	} else { // If a pipeline name is provided, launch the pipeline with that name
		logging.Logger.Info("Looking for pipeline using name", "name", envvars.Variables["NAME"])
		rawPipelines = pipelineprocessor.FindPipelineByName(combinedData, envvars.Variables, envvars.Variables["NAME"])
	}
	if len(rawPipelines) == 0 {
		logging.Logger.Warn("No pipelines found")
		os.Exit(ErrCodeOK)
	}

	// WIP: Launch the pipelines
	for name, pipeline := range rawPipelines {
		logging.Logger.Info("WIP: Launching pipeline", "name", name)

		// TODO: Validate raw pipeline

		// Convert pipeline to PipelineSpec
		spec, err := convertToPipelines(pipeline)
		if err != nil {
			logging.Logger.Error("Error converting pipeline to PipelineSpec", "error", err)
			os.Exit(ErrCodeConvertingPipeline)
		}

		// TODO:
		// Create namespace
		namespace := spec.Namespace.Name
		// ...

		// Remove namespace from spec once it's created
		spec.Namespace = pipelinecrd.Namespace{}

		// Deploy the pipeline
		err = deploy.Pipeline(name, namespace, spec)
		if err != nil {
			logging.Logger.Error("Error deploying pipeline", "error", err)
			os.Exit(ErrCodeDeploy)
		}

		logging.Logger.Info("Pipeline deployed successfully", "name", name, "namespace", namespace)
	}

	return
}

// getMD5Hash returns the MD5 hash of the text
func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// convertToPipelines converts the raw data to a PipelineSpec struct
// TODO: Duplicated code from internal/app/pipeline-controller/normalize/normalize.go
// This function convert only one raw pipeline to a PipelineSpec struct
// This would be the future version for normalize.go when it's refactored
func convertToPipelines(data interface{}) (pipelinecrd.PipelineSpec, error) {
	// Convert each item to YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return pipelinecrd.PipelineSpec{}, err
	}

	// Unmarshal YAML to PipelineSpec struct
	var pipeline pipelinecrd.PipelineSpec
	err = yaml.Unmarshal(yamlData, &pipeline)
	if err != nil {
		return pipelinecrd.PipelineSpec{}, err
	}

	return pipeline, nil
}
