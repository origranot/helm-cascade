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
	if err != nil || stderr.String() != "" {
		return "", stderr.String(), err
	}

	lines := strings.Split(stdout.String(), "\n")
	formattedDeps := make([]string, 0, len(lines)-1)

	if len(lines) > 1 {
		for _, line := range lines[1:] {
			formattedDeps = append(formattedDeps, formatDependency(line))
		}

		return strings.Join(formattedDeps, "\n"), "", nil
	}
	return "", stderr.String(), nil
}

func formatDependency(dep string) string {
	if strings.HasPrefix(dep, "WARNING:") {
		return color.YellowString(dep)
	}

	parts := strings.Split(dep, "\t")
	if len(parts) < 4 {
		return dep
	}

	name := strings.TrimSpace(parts[0])
	version := strings.TrimSpace(parts[1])
	repo := strings.TrimSpace(parts[2])
	status := strings.TrimSpace(parts[3])

	switch strings.ToLower(status) {
	case "ok", "unpacked":
		status = color.GreenString(status)
	case "missing":
		status = color.RedString(status)
	default:
		status = color.YellowString(status)
	}

	namePad := 25
	versionPad := 8
	repoPad := 52

	return fmt.Sprintf("%-*s %-*s %-*s %s",
		namePad, name,
		versionPad, version,
		repoPad, repo,
		status)
}
