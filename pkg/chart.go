package helmutil

import (
	"os"
	"path/filepath"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
)

type Chart struct {
	metadata chart.Metadata
	path     string
	level    int
}

func walkAndFindCharts(rootDir string, charts *[]Chart) error {
	_, err := chartutil.IsChartDir(rootDir)
	if err != nil {
		return err
	}

	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return err
		}

		chartPath := filepath.Join(path, CHART_FILE_NAME)
		chartMetadata, err := chartutil.LoadChartfile(chartPath)
		if err != nil {
			return nil
		}

		chartLevel := strings.Count(strings.TrimPrefix(path, rootDir), string(filepath.Separator))

		*charts = append(*charts, Chart{
			metadata: *chartMetadata,
			path:     path,
			level:    chartLevel,
		})

		return nil
	})
}
