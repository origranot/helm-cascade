package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cascade [command] CHART",
	Short: "A tool to manage Helm dependencies recursively",
	Long:  `This is a tool to manage Helm dependencies recursively.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	updateCmd.Aliases = []string{"up"}
}
