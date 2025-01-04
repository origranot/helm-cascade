package cmd

var buildCmd = createDependencyCommand("build", "Build dependencies across all sub charts")

func init() {
	rootCmd.AddCommand(buildCmd)
}
