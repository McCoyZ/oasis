package app

import (
	"fmt"

	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
	"zmc.io/oasis/cmd/apiserver/app/options"
	apiserverconfig "zmc.io/oasis/pkg/apiserver/config"
	"zmc.io/oasis/pkg/utils/signals"
	"zmc.io/oasis/pkg/utils/term"
)

func NewAPIServerCommand() *cobra.Command {
	s := options.NewServerRunOptions()

	// Load configuration from file
	conf, err := apiserverconfig.TryLoadFromDisk()
	if err == nil {
		s = &options.ServerRunOptions{
			GenericServerRunOptions: s.GenericServerRunOptions,
			Config:                  conf,
		}
	}

	cmd := &cobra.Command{
		Use: "ks-apiserver",
		Long: `The Kubernets API server validates and configures data for the API objects. 
The API Server services REST operations and provides the frontend to the
cluster's shared state through which all other components interact.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if errs := s.Validate(); len(errs) != 0 {
				return utilerrors.NewAggregate(errs)
			}
			return Run(s, signals.SetupSignalHandler())
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
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
	return cmd
}

func Run(s *options.ServerRunOptions, stopCh <-chan struct{}) error {

	apiserver, err := s.NewAPIServer(stopCh)
	if err != nil {
		return err
	}

	err = apiserver.PrepareRun(stopCh)
	if err != nil {
		return nil
	}

	return apiserver.Run(stopCh)
}
