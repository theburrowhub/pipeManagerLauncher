package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/sergiotejon/pipeManager/internal/pkg/version"
)

// rootCmd is the root command for the CLI
var rootCmd = &cobra.Command{
	Use:   "pipieline-controller",
	Short: "Kubernetes pipeline controller",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

var configFile string // configFile is the path to the configuration file

// init initializes the application
func init() {
	// Persistent flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to the config file")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version")

	// Bind the version flag
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("version")
		if v {
			fmt.Println(version.GetVersion())
			os.Exit(0)
		}

		// TODO: Normalize
		fmt.Println("Normalize... Under construction.")
		// -- Read spec from k8s object (only one pipeline)
		// -- Refactor Normalize to work only with one pipeline
		// Normalize the pipelines
		//pipelines, err := normalize.Normalize(rawPipelines)
		//if err != nil {
		//	logging.Logger.Error("Error normalizing pipelines", "msg", err)
		//	os.Exit(ErrCodeNormalize)
		//}
	}
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
