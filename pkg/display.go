package helmutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"helm.sh/helm/v3/pkg/chart"
)

type DependencyStatus string

const (
	StatusOK       DependencyStatus = "ok"
	StatusMissing  DependencyStatus = "missing"
	StatusUnpacked DependencyStatus = "unpacked"
	StatusMismatch DependencyStatus = "version_mismatch"
)

type DependencyInfo struct {
	Name    string
	Version string
	Status  DependencyStatus
	Message string
}

func checkDependencyStatus(chartPath string, depName, depVersion string) DependencyInfo {
	cmd := exec.Command("helm", "dependency", "list", chartPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errorMsg := strings.TrimSpace(stderr.String())
		if errorMsg == "" {
			errorMsg = "failed to run helm dependency list"
		}
		return DependencyInfo{
			Name:    depName,
			Version: depVersion,
			Status:  StatusMissing,
			Message: errorMsg,
		}
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return DependencyInfo{Name: depName, Version: depVersion, Status: StatusMissing, Message: "no output from helm"}
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[0] == depName {
			status := fields[len(fields)-1]
			return parseDependencyStatus(depName, depVersion, status)
		}
	}

	return DependencyInfo{Name: depName, Version: depVersion, Status: StatusMissing, Message: "not found in helm output"}
}

func parseDependencyStatus(depName, depVersion, status string) DependencyInfo {
	switch status {
	case "ok":
		return DependencyInfo{Name: depName, Version: depVersion, Status: StatusOK, Message: "ok"}
	case "missing":
		return DependencyInfo{Name: depName, Version: depVersion, Status: StatusMissing, Message: "missing"}
	case "unpacked":
		return DependencyInfo{Name: depName, Version: depVersion, Status: StatusUnpacked, Message: "unpacked"}
	case "version":
		return DependencyInfo{Name: depName, Version: depVersion, Status: StatusMismatch, Message: "version mismatch"}
	default:
		return DependencyInfo{Name: depName, Version: depVersion, Status: StatusMissing, Message: status}
	}
}

func ListSubchartDependencies(rootDir string) error {
	var charts []Chart
	err := walkAndFindCharts(rootDir, &charts)
	if err != nil {
		return err
	}

	if len(charts) == 0 {
		fmt.Println(color.New(color.FgHiBlack).Sprint("No charts found"))
		return nil
	}

	for i, chart := range charts {
		prefix := strings.Repeat("│   ", chart.level)

		var chartTreeChar string
		if i == len(charts)-1 {
			chartTreeChar = "└──"
		} else {
			chartTreeChar = "├──"
		}

		chartName := ChartColor.Sprint(chart.metadata.Name)
		fmt.Printf("%s%s %s\n", prefix, chartTreeChar, chartName)

		if len(chart.metadata.Dependencies) > 0 {
			for _, dep := range chart.metadata.Dependencies {
				if strings.TrimSpace(dep.Name) != "" {
					depInfo := checkDependencyStatus(chart.path, dep.Name, dep.Version)
					displayDependency(prefix, dep, depInfo)
				}
			}
		} else {
			fmt.Printf("%s│   %s\n", prefix, VersionColor.Sprint("No dependencies"))
		}

	}

	return nil
}

func displayDependency(prefix string, dep *chart.Dependency, depInfo DependencyInfo) {
	statusIndicator, statusColor := getStatusDisplay(depInfo.Status)

	depName := DepColor.Sprint(dep.Name)
	if dep.Alias != "" {
		depName += fmt.Sprintf(" (%s)", AliasColor.Sprint(dep.Alias))
	}

	version := VersionColor.Sprint(dep.Version)
	status := statusColor.Sprint(depInfo.Message)

	fmt.Printf("%s│   %s %s %s %s\n",
		prefix,
		statusIndicator,
		depName,
		version,
		status,
	)
}

func getStatusDisplay(status DependencyStatus) (string, *color.Color) {
	switch status {
	case StatusOK:
		return SuccessColor.Sprint("✓"), SuccessColor
	case StatusMissing:
		return ErrorColor.Sprint("✗"), ErrorColor
	case StatusUnpacked:
		return WarningColor.Sprint("📦"), WarningColor
	case StatusMismatch:
		return ErrorColor.Sprint("⚠"), WarningColor
	default:
		return "?", VersionColor
	}
}
