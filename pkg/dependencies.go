package helmutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Command string

const (
	Build  Command = "build"
	Update Command = "update"
)

func ManageDependencies(rootDir string, cmd Command) error {
	var charts []ChartInfo
	err := collectCharts(rootDir, &charts)
	if err != nil {
		return err
	}

	for i := len(charts) - 1; i >= 0; i-- {
		chart := charts[i]
		helmCmd := exec.Command("helm", "dependency", string(cmd))
		helmCmd.Dir = chart.Path

		var stdout, stderr bytes.Buffer
		helmCmd.Stdout = &stdout
		helmCmd.Stderr = &stderr

		if err := helmCmd.Run(); err != nil {
			return fmt.Errorf("failed to %s dependencies for %s:\n%s", cmd, chart.Path, color.RedString(stderr.String()))
		}
		if stderr.String() != "" {
			fmt.Printf("%s\n", color.YellowString(stderr.String()))
		}

		var successMsg string
		if cmd == Build {
			successMsg = "Successfully built"
		} else {
			successMsg = "Successfully updated"
		}

		fmt.Printf("%s %s\n", color.BlueString(chart.Name), successMsg)
	}
	return nil
}

func collectCharts(rootDir string, charts *[]ChartInfo) error {

	_, err := os.Stat(path.Join(rootDir, "Chart.yaml"))
	if err != nil {
		return fmt.Errorf("chart not found: %s", rootDir)
	}

	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return err
		}

		chartPath := filepath.Join(path, "Chart.yaml")
		if _, err := os.Stat(chartPath); err == nil {
			level := strings.Count(strings.TrimPrefix(path, rootDir), string(filepath.Separator))
			*charts = append(*charts, ChartInfo{
				Name:  filepath.Base(path),
				Path:  path,
				Level: level,
			})
		}
		return nil
	})
}
