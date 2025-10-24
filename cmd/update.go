package cmd

var updateCmd = createDependencyCommand("update", "Update dependencies across all sub charts")

func init() {
	updateCmd.Aliases = append(updateCmd.Aliases, "up")
	rootCmd.AddCommand(updateCmd)
}
