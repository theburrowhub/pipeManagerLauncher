package cmd

import (
	"flag"
	"fmt"
	"github.com/sergiotejon/pipeManager/internal/pkg/version"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/sergiotejon/pipeManager/internal/pkg/pipelinecrd"
)

const leaderElectionID = "pipeline-manager-controller"

var logger = ctrl.Log.WithName("pipeline-controller")

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(pipelinecrd.Scheme))
	utilruntime.Must(pipelinecrd.AddToScheme(pipelinecrd.Scheme))
}

// app runs the application
func app() {
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	if versionFlag {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 pipelinecrd.Scheme,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       leaderElectionID,
	})
	if err != nil {
		logger.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&PipelineReconcile{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		logger.Error(err, "unable to create controller", "controller", "Pipeline")
		os.Exit(1)
	}

	// Añadir verificaciones de salud
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		logger.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		logger.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	logger.Info("Starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		logger.Error(err, "problem running manager")
		os.Exit(1)
	}
}
