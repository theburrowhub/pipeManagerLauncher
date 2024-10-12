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
	Use:   "launcher",
	Short: "Pipe Manager Launcher CLI",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

// init initializes the application
func init() {
	// Persistent flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", defaultConfigFile, "Path to the config file")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version")

	// Bind the version flag
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("version")
		if v {
			fmt.Println(version.GetVersion())
			os.Exit(0)
		}
	}

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(artifactsCmd)
	rootCmd.AddCommand(cacheCmd)
}
