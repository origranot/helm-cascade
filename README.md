# Helm Cascade Plugin

A Helm plugin for recursive dependency management across all subcharts.

## Motivation

I needed a way to manage dependencies across multiple subcharts in a large Helm chart. Helm's built-in dependency management only handles dependencies for a single chart, requiring manual updates for each subchart. This becomes tedious and error-prone with complex chart hierarchies. Cascade automates dependency management by recursively processing the entire chart tree in one command.

## Installation

```bash
helm plugin install https://github.com/origranot/helm-cascade
```

## Usage

### List dependencies

```bash
helm cascade list <chart-dir>
```

Example output:

```bash
├── parent-chart
│   subchart-a    1.0.0    repo/subchart-a    unpacked
│   subchart-b    2.0.0    repo/subchart-b    unpacked
│   ├── subchart-a
│   │   dependency-x    1.2.0    repo/dep-x    missing # This would never been shown by helm dependency list
```

### Build dependencies

```bash
helm cascade build <chart-dir>
```

### Update dependencies

```bash
helm cascade update <chart-dir> # or helm cascade up
```
