// Package cmd
// cmd package contains the Cobra commands for the launcher application.
package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/sergiotejon/pipeManager/internal/app/launcher/pipelineprocessor"
	"github.com/sergiotejon/pipeManager/internal/app/launcher/repository"
	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/envvars"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

const (
	ErrCodeOK             = 0
	ErrCodeLoadConfig     = 1
	ErrCodeCloneRepo      = 2
	ErrCodeMixFiles       = 3
	ErrCodeNormalize      = 4
	ErrCodeBucketDownload = 6
	ErrCodeBucketUpload   = 7
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
	logging.Logger.Info("Setup", "configFile", configFile,
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

	// Normalize the pipelines
	pipelines, err := pipelineprocessor.Normalize(rawPipelines)
	if err != nil {
		logging.Logger.Error("Error normalizing pipelines", "msg", err)
		os.Exit(ErrCodeNormalize)
	}
	d, _ := yaml.Marshal(pipelines)
	fmt.Println(string(d))

	// Temporal
	//if config.Common.Data.Log.Level == "debug" {
	//	data, err := yaml.Marshal(pipelines)
	//	if err != nil {
	//		os.Exit(1)
	//	}
	//	fmt.Println(string(data))
	//}
	// Temporal

	for key, _ := range pipelines {
		logging.Logger.Info("Launching pipeline", "pipeline", key)
		// TODO:
		// - Launch the pipeline (pipeline controller)
	}

	return
}

// getMD5Hash returns the MD5 hash of the text
func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
