# Helm Cascade Plugin

A Helm plugin for recursive dependency management across all subcharts in complex chart hierarchies.

## 🚀 Features

- **Recursive Dependency Management**: Automatically processes dependencies across all subcharts in a single command
- **Visual Tree Display**: Tree view showing the complete dependency hierarchy with status indicators
- **Batch Operations**: Build, update, and lint all charts in the correct dependency order

## 🎯 Motivation

Managing dependencies in complex Helm chart hierarchies is challenging. Helm's built-in `dependency` commands only work on a single chart level, requiring you to manually navigate to each subchart directory and run commands individually. This becomes:

- **Time-consuming**: Manual navigation through dozens of subcharts
- **Error prone**: Easy to miss dependencies or run commands in the wrong order
- **Inconsistent**: Different team members might handle dependencies differently

Cascade solves this by automatically discovering all charts in your hierarchy and processing them in the correct dependency order.

## 📦 Installation

```bash
helm plugin install https://github.com/origranot/helm-cascade
```

## 🛠️ Usage

### List Dependencies

View the complete dependency tree with status indicators:

```bash
helm cascade list <chart-dir>
```

**Example Output:**

```bash
├── e-commerce-platform
│   ✓ web-frontend        1.2.3    ok
│   📦 api-gateway        2.1.0    unpacked
│   ✓ payment-service     0.5.2    ok
│   ├── web-frontend
│   │   ✓ nginx-ingress      3.1.0    ok
│   │   ✓ cert-manager       2.4.1    ok
│   │   ✗ monitoring-stack   1.0.0    missing
│   ├── api-gateway
│   │   ✓ istio-gateway      1.8.0    ok
│   │   ✓ rate-limiter       0.3.2    ok
│   │   ├── istio-gateway
│   │   │   ✓ istio-base         2.1.0    ok
│   │   │   ✓ istiod             1.5.0    ok
```

**Status Indicators:**

- ✓ **ok**: Dependency is properly installed and up-to-date
- ✗ **missing**: Dependency needs to be downloaded
- 📦 **unpacked**: Dependency is downloaded but not installed
- ⚠ **version_mismatch**: Installed version doesn't match requirements

### Build Dependencies

Download and install all dependencies across the entire chart hierarchy:

```bash
helm cascade build <chart-dir>
```

### Update Dependencies

Update all dependencies to their latest versions:

```bash
helm cascade update <chart-dir>

# or use the short alias:
helm cascade up <chart-dir>
```

### Lint Charts

Validate all charts in the hierarchy:

```bash
helm cascade lint <chart-dir>
```

## 🏗️ Real-World Example

Consider an e-commerce platform with this structure:

```
e-commerce-platform/
├── Chart.yaml
├── values.yaml
├── charts/
│   ├── web-frontend/
│   │   ├── Chart.yaml
│   │   ├── charts/
│   │   │   ├── nginx-ingress/
│   │   │   │   ├── Chart.yaml
│   │   │   │   └── charts/
│   │   │   │       └── cert-manager/
│   │   │   └── monitoring-stack/
│   │   └── charts/
│   │       └── redis-cache/
│   ├── api-gateway/
│   │   ├── Chart.yaml
│   │   └── charts/
│   │       ├── istio-gateway/
│   │       │   ├── Chart.yaml
│   │       │   └── charts/
│   │       │       ├── istio-base/
│   │       │       └── istiod/
│   │       └── rate-limiter/
│   └── payment-service/
```

**Without Cascade:**

```bash
# Manual process - error-prone and time-consuming
cd e-commerce-platform
helm dependency build
cd charts/web-frontend
helm dependency build
cd charts/nginx-ingress
helm dependency build
cd charts/cert-manager
helm dependency build
cd charts/monitoring-stack
helm dependency build
cd charts/redis-cache
helm dependency build
cd ../api-gateway
helm dependency build
cd charts/istio-gateway
helm dependency build
cd charts/istio-base
helm dependency build
cd charts/istiod
helm dependency build
cd ../rate-limiter
helm dependency build
cd ../../payment-service
helm dependency build
# ... repeat for every chart
```

**With Cascade:**

```bash
# One command handles everything
helm cascade build e-commerce-platform
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
