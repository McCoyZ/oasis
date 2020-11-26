package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"zmc.io/oasis/cmd/controller-manager/app/options"
	apis "zmc.io/oasis/pkg/apis/servicemesh/v1alpha1"
	controllerconfig "zmc.io/oasis/pkg/apiserver/config"
	"zmc.io/oasis/pkg/informers"
	"zmc.io/oasis/pkg/simple/client/k8s"
	"zmc.io/oasis/pkg/utils/term"
)

func NewControllerManagerCommand() *cobra.Command {
	s := options.NewControllerManagerOptions()
	conf, err := controllerconfig.TryLoadFromDisk()
	if err == nil {
		// make sure LeaderElection is not nil
		s = &options.ControllerManagerOptions{
			KubernetesOptions: conf.KubernetesOptions,
			// ServiceMeshOptions: conf.ServiceMeshOptions,
			LeaderElection: s.LeaderElection,
			LeaderElect:    s.LeaderElect,
		}
	} else {
		klog.Fatal("Failed to load configuration from disk", err)
	}

	cmd := &cobra.Command{
		Use:  "controller-manager",
		Long: `Mesh controller manager is a daemon that`,
		Run: func(cmd *cobra.Command, args []string) {
			if errs := s.Validate(); len(errs) != 0 {
				klog.Error(utilerrors.NewAggregate(errs))
				os.Exit(1)
			}

			if err = run(s, signals.SetupSignalHandler()); err != nil {
				klog.Error(err)
				os.Exit(1)
			}
		},
		SilenceUsage: true,
	}

	fs := cmd.Flags()
	namedFlagSets := s.Flags()

	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
	return cmd
}

func run(s *options.ControllerManagerOptions, stopCh <-chan struct{}) error {
	kubernetesClient, err := k8s.NewKubernetesClient(s.KubernetesOptions)
	if err != nil {
		klog.Errorf("Failed to create kubernetes clientset %v", err)
		return err
	}

	informerFactory := informers.NewInformerFactories(
		kubernetesClient.Kubernetes(),
		kubernetesClient.Mesh(),
		kubernetesClient.Istio(),
	)

	mgrOptions := manager.Options{}

	if s.LeaderElect {
		mgrOptions = manager.Options{
			LeaderElection:          s.LeaderElect,
			LeaderElectionNamespace: "linkedcare-system",
			LeaderElectionID:        "ms-controller-manager-leader-election",
			LeaseDuration:           &s.LeaderElection.LeaseDuration,
			RetryPeriod:             &s.LeaderElection.RetryPeriod,
			RenewDeadline:           &s.LeaderElection.RenewDeadline,
		}
	}

	klog.V(0).Info("setting up manager")

	// Use 8443 instead of 443 cause we need root permission to bind port 443
	mgr, err := manager.New(kubernetesClient.Config(), mgrOptions)
	if err != nil {
		klog.Fatalf("unable to set up overall controller manager: %v", err)
	}

	if err = apis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatalf("unable add APIs to scheme: %v", err)
	}

	// servicemeshEnabled := s.ServiceMeshOptions != nil && len(s.ServiceMeshOptions.IstioPilotHost) != 0
	servicemeshEnabled := true

	// Add controllers
	if err = addControllers(mgr,
		kubernetesClient,
		informerFactory,
		servicemeshEnabled,
		stopCh); err != nil {
		klog.Fatalf("unable to register controllers to the manager: %v", err)
	}

	// Start cache data after all informer is registered
	klog.V(0).Info("Starting cache resource from apiserver...")
	informerFactory.Start(stopCh)

	klog.V(0).Info("Starting the controllers.")
	if err = mgr.Start(stopCh); err != nil {
		klog.Fatalf("unable to run the manager: %v", err)
	}

	return nil
}
