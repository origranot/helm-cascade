package cmd

import (
	"path/filepath"

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
	Run: func(cmd *cobra.Command, args []string) {
		absPath, err := filepath.Abs(args[0])

		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		err = helmutil.ListSubchartDependencies(absPath)

		if err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
