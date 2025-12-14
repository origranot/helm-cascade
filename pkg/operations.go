package helmutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"sync"
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

	chartsByLevel := make(map[int][]Chart)
	maxLevel := 0
	for _, chart := range charts {
		level := chart.level
		chartsByLevel[level] = append(chartsByLevel[level], chart)
		if level > maxLevel {
			maxLevel = level
		}
	}

	// Process levels from deepest to shallowest (maxLevel down to 0)
	// Charts at the same level can be processed in parallel
	for level := maxLevel; level >= 0; level-- {
		levelCharts, exists := chartsByLevel[level]
		if !exists {
			continue
		}

		// Process all charts at this level in parallel
		var wg sync.WaitGroup
		var mu sync.Mutex
		var firstErr error

		for _, chart := range levelCharts {
			wg.Add(1)
			go func(ch Chart) {
				defer wg.Done()

				var helmCmd *exec.Cmd

				switch operation {
				case OperationDependency:
					helmCmd = exec.Command("helm", "dependency", string(cmd), ch.path)
				case OperationLint:
					helmCmd = exec.Command("helm", "lint", ch.path, "--values", path.Join(rootDir, "values.yaml"))
				default:
					mu.Lock()
					if firstErr == nil {
						errorMsg := fmt.Sprintf("Error: unsupported operation: %s", operation)
						firstErr = fmt.Errorf("%s", ErrorColor.Sprint(errorMsg))
					}
					mu.Unlock()
					return
				}

				var stdout, stderr bytes.Buffer
				helmCmd.Stdout = &stdout
				helmCmd.Stderr = &stderr

				if err := helmCmd.Run(); err != nil {
					mu.Lock()
					if firstErr == nil {
						switch operation {
						case OperationDependency:
							errorMsg := fmt.Sprintf("Error: failed to %s dependencies for %s:\n%s", cmd, ch.path, stderr.String())
							firstErr = fmt.Errorf("%s", ErrorColor.Sprint(errorMsg))
						case OperationLint:
							errorMsg := fmt.Sprintf("Error: failed to lint chart %s:\n%s", ch.path, stderr.String())
							firstErr = fmt.Errorf("%s", ErrorColor.Sprint(errorMsg))
						}
					}
					mu.Unlock()
					return
				}

				// Synchronize output to avoid interleaving
				mu.Lock()
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
					fmt.Printf("%s %s\n", ChartColor.Sprint(ch.metadata.Name), successMsg)
				case OperationLint:
					fmt.Printf("%s linted successfully\n", ChartColor.Sprint(ch.metadata.Name))
				}
				mu.Unlock()
			}(chart)
		}

		// Wait for all charts at this level to complete
		wg.Wait()

		// If any chart at this level failed, return the error
		if firstErr != nil {
			return firstErr
		}
	}

	return nil
}
