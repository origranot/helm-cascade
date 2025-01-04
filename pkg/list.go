package helmutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type ChartInfo struct {
	Name         string
	Path         string
	Dependencies string
	Level        int
}

func ListSubchartDependencies(rootDir string) error {
	var charts []ChartInfo
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return err
		}

		chartPath := filepath.Join(path, "Chart.yaml")
		if _, err := os.Stat(chartPath); err == nil {
			level := strings.Count(strings.TrimPrefix(path, rootDir), string(filepath.Separator))
			stdout, stderr, err := getDependencies(path)
			if stderr != "" {
				fmt.Printf("%s%s\n", strings.Repeat("│   ", level), color.YellowString(stderr))
			}
			if err == nil {
				charts = append(charts, ChartInfo{
					Name:         filepath.Base(path),
					Dependencies: stdout,
					Level:        level,
				})
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(charts) == 0 {
		fmt.Println("No charts found")
		return nil
	}

	for _, chart := range charts {
		prefix := strings.Repeat("│   ", chart.Level)
		fmt.Printf("%s├── %s\n", prefix, color.BlueString(chart.Name))

		if chart.Dependencies != "" {
			for _, dep := range strings.Split(chart.Dependencies, "\n") {
				if strings.TrimSpace(dep) != "" {
					fmt.Printf("%s│   %s\n", prefix, dep)
				}
			}
		}
	}
	return nil
}

func getDependencies(path string) (string, string, error) {
	cmd := exec.Command("helm", "dependency", "list")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", stderr.String(), err
	}

	lines := strings.Split(stdout.String(), "\n")
	if len(lines) > 1 {
		return strings.Join(lines[1:], "\n"), stderr.String(), nil
	}
	return "", stderr.String(), nil
}
