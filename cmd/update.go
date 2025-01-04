package cmd

var updateCmd = createDependencyCommand("update", "Update dependencies across all sub charts")

func init() {
	rootCmd.AddCommand(updateCmd)
}
