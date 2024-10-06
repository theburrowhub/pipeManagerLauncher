// Package main contains the main entrypoint for the application.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/sergiotejon/pipeManager/internal/pkg/version"
)

const (
	defaultConfigFile = "/etc/pipe-manager/config.yaml" // defaultConfigFile is the default configuration file
	templateFolder    = "/etc/pipe-manager/templates"   // templateFolder is the folder where the templates are stored
	repoDir           = "/tmp/repo"                     // repoDir is the directory where the repository is cloned
)

var (
	configFile  string // configFile is the path to the configuration file
	showVersion bool   // showVersion is a flag to show the version
)

const (
	ErrCodeOK             = 0
	ErrCodeLoadConfig     = 1
	ErrCodeCloneRepo      = 2
	ErrCodeMixFiles       = 3
	ErrCodeNormalize      = 4
	ErrCodeClone          = 5
	ErrCodeBucketDownload = 6
	ErrCodeBucketUpload   = 7
)

// main is the entrypoint for the application
// It sets up the root command and executes the application
func main() {
	rootCmd := &cobra.Command{
		Use:   "pipe-manager-launcher",
		Short: "Pipe Manager Launcher CLI",
		Run: func(cmd *cobra.Command, args []string) {
			// Show version
			if showVersion {
				fmt.Println(version.GetVersion())
				os.Exit(0)
			}

			// Run the application
			initApp()
			app()
		},
	}

	rootCmd.Flags().StringVarP(&configFile, "config", "c", defaultConfigFile, "Path to the config file")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print the version")

	// Add commands to rootCmd
	rootCmd.AddCommand(cloneCmd, bucketCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
