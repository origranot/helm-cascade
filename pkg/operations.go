package helmutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
)

type DependencyCommand string

const (
	Build  DependencyCommand = "build"
	Update DependencyCommand = "update"
)

type ChartOperation string

const (
	OperationDependency ChartOperation = "dependency"
	OperationLint       ChartOperation = "lint"
)

func ProcessCharts(rootDir string, operation ChartOperation, cmd DependencyCommand) error {
	var charts []Chart
	err := walkAndFindCharts(rootDir, &charts)
	if err != nil {
		return err
	}

	for i := len(charts) - 1; i >= 0; i-- {
		chart := charts[i]
		var helmCmd *exec.Cmd

		switch operation {
		case OperationDependency:
			helmCmd = exec.Command("helm", "dependency", string(cmd), chart.path)
		case OperationLint:
			helmCmd = exec.Command("helm", "lint", chart.path, "--values", path.Join(rootDir, "values.yaml"))
		default:
			errorMsg := fmt.Sprintf("Error: unsupported operation: %s", operation)
			return fmt.Errorf("%s", ErrorColor.Sprint(errorMsg))
		}

		var stdout, stderr bytes.Buffer
		helmCmd.Stdout = &stdout
		helmCmd.Stderr = &stderr

		if err := helmCmd.Run(); err != nil {
			switch operation {
			case OperationDependency:
				errorMsg := fmt.Sprintf("Error: failed to %s dependencies for %s:\n%s", cmd, chart.path, stderr.String())
				return fmt.Errorf("%s", ErrorColor.Sprint(errorMsg))
			case OperationLint:
				errorMsg := fmt.Sprintf("Error: failed to lint chart %s:\n%s", chart.path, stderr.String())
				return fmt.Errorf("%s", ErrorColor.Sprint(errorMsg))
			}
		}
		if stderr.String() != "" {
			fmt.Printf("%s\n", WarningColor.Sprint(stderr.String()))
		}

		switch operation {
		case OperationDependency:
			var successMsg string
			if cmd == Build {
				successMsg = "dependencies built successfully"
			} else {
				successMsg = "dependencies updated successfully"
			}
			fmt.Printf("%s %s\n", ChartColor.Sprint(chart.metadata.Name), successMsg)
		case OperationLint:
			fmt.Printf("%s linted successfully\n", ChartColor.Sprint(chart.metadata.Name))
		}
	}
	return nil
}
