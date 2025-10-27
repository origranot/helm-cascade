package cmd

import (
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
		RunE: func(_ *cobra.Command, args []string) error {
			chartPath := args[0]
			if name == "lint" {
				return helmutil.ProcessCharts(chartPath, helmutil.OperationLint, "")
			}
			return helmutil.ProcessCharts(chartPath, helmutil.OperationDependency, helmutil.DependencyCommand(name))
		},
	}
}
