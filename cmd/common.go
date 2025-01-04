package cmd

import (
	"path/filepath"

	helmutil "github.com/origranot/helm-cascade/pkg"
	"github.com/spf13/cobra"
)

func createDependencyCommand(name, desc string) *cobra.Command {
	return &cobra.Command{
		Use:                   name + " CHART",
		Short:                 desc,
		DisableFlagsInUseLine: true,
		Example:               "helm cascade " + name + " .",
		Args:                  cobra.ExactArgs(1),
		ValidArgs:             []string{"chart"},
		Run: func(cmd *cobra.Command, args []string) {
			absPath, err := filepath.Abs(args[0])
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			err = helmutil.ManageDependencies(absPath, helmutil.Command(name))
			if err != nil {
				cmd.PrintErrln(err)
			}
		},
	}
}
