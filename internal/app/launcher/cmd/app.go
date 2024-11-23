// Package cmd
// cmd package contains the Cobra commands for the launcher application.
package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"

	"github.com/sergiotejon/pipeManagerLauncher/internal/app/launcher/convert"
	"github.com/sergiotejon/pipeManagerLauncher/internal/app/launcher/deploy"
	"github.com/sergiotejon/pipeManagerLauncher/internal/app/launcher/namespace"
	"github.com/sergiotejon/pipeManagerLauncher/internal/app/launcher/pipelineprocessor"
	"github.com/sergiotejon/pipeManagerLauncher/internal/app/launcher/repository"
	"github.com/sergiotejon/pipeManagerLauncher/internal/pkg/logging"
	"github.com/sergiotejon/pipeManagerLauncher/pkg/config"
	"github.com/sergiotejon/pipeManagerLauncher/pkg/envvars"
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
	envvar_prefix  = "PIPELINE_"
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
	envvars.GetEnvVars(envvar_prefix)
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

	// Launch the pipelines
	for name, pipeline := range rawPipelines {
		// --- DEBUG
		//yamlData, err := yaml.Marshal(pipeline)
		//if err != nil {
		//	logging.Logger.Error("Error marshaling pipeline to YAML", "pipeline", name, "error", err)
		//	continue
		//}
		//fmt.Println(string(yamlData))
		// --- DEBUG

		logging.Logger.Info("Launching pipeline", "name", name)
		// TODO: Validate raw pipeline?. A validation is done when converting to PipelineSpec

		// Convert pipeline to PipelineSpec
		spec, err := convert.ConvertToPipelines(pipeline)
		if err != nil {
			logging.Logger.Error("Error converting pipeline to PipelineSpec. Pipeline not deployed",
				"pipeline", name, "error", err)
			continue
		}

		// Create namespace
		namespaceName := spec.Namespace.Name
		err = namespace.Create(spec)
		if err != nil {
			logging.Logger.Error("Error creating namespace. Pipeline not deployed",
				"namespace", namespaceName, "pipeline", name, "error", err)
			continue
		}

		// Deploy the pipeline
		resourceName, resourceNamespace, err := deploy.Pipeline(name, namespaceName, spec)
		if err != nil {
			logging.Logger.Error("Error deploying pipeline", "error", err)
			continue
		}

		logging.Logger.Info("Pipeline deployed successfully",
			"name", name, "resourceName", resourceName, "resourceNamespace", resourceNamespace)
	}

	return
}

// getMD5Hash returns the MD5 hash of the text
func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
