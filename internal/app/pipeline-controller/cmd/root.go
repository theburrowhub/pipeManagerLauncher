package cmd

import (
	"flag"
)

const (
	ErrCodeOK         = 0
	ErrCodeLoadConfig = 1
)

//// rootCmd is the root command for the CLI
//var rootCmd = &cobra.Command{
//	Use:   "pipieline-controller",
//	Short: "Kubernetes pipeline manager controller",
//	Run: func(cmd *cobra.Command, args []string) {
//		err := cmd.Help()
//		if err != nil {
//			return
//		}
//	},
//}

var (
	configFile           string // configFile is the path to the configuration file
	versionFlag          bool   // versionFlag is the flag to print the version
	probeAddr            string // probeAddr is the address the probe endpoint binds to
	enableLeaderElection bool   // enableLeaderElection enables leader election for controller manager
)

//
//// init initializes the application
//func init() {
//	// Persistent flags
//	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to the config file")
//	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version")
//	rootCmd.Flags().StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
//	rootCmd.PersistentFlags().BoolVar(&enableLeaderElection, "leader-elect", false,
//		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
//	rootCmd.Flags().StringVar(&logLevel, "zap-log-level", "info", "Log level")
//
//	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
//		v, err := cmd.Flags().GetBool("version")
//		if err != nil {
//			log.Fatalf("Error getting version flag: %v", err)
//		}
//		if v {
//			fmt.Println(version.GetVersion())
//			os.Exit(0)
//		}
//
//		// Set up the app
//		setup()
//		// Run the app
//		app()
//
//		os.Exit(0)
//	}
//}

// Execute executes the root command
func Execute() {
	flag.StringVar(&configFile, "config", "", "Path to the config file")
	flag.BoolVar(&versionFlag, "version", false, "Print the version")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	setup()
	app()

	//if err := rootCmd.Execute(); err != nil {
	//	log.Fatalf("Error executing command: %v", err)
	//}
}
