package cmd

import (
	helmutil "github.com/origranot/helm-cascade/pkg"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:                   "list CHART",
	Short:                 "List all dependencies across all sub charts",
	Example:               "helm cascade list .",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	ValidArgs:             []string{"chart"},
	RunE: func(_ *cobra.Command, args []string) error {
		return helmutil.ListSubchartDependencies(args[0])
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
