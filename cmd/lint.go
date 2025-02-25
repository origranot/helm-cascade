package cmd

var lintCmd = createDependencyCommand("lint", "Lint chart including his subcharts with correct values scope (this must be run from the root of the chart)")

func init() {
	rootCmd.AddCommand(lintCmd)
}
