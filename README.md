# kube-pod-update

**kube-pod-update** is a CLI tool built in Go that scans your Kubernetes cluster for pods running outdated container images. It compares the current images to the latest available versions in container registries and optionally updates deployments automatically.

---

## 🚀 Features

- **Detect Outdated Images:** Finds pods using container images that have newer versions available.
- **Automatic Updates:** Optionally updates Kubernetes deployments to use the latest image versions.
- **Supports Major Registries:** Compatible with Docker Hub and other OCI-compliant container registries.
- **Namespace Filtering:** Specify namespaces to target specific groups of pods.
- **Colored Logs:** Clearly readable, color-coded log outputs for ease of use.
- **Flexible Configuration:** Configure via environment variables, YAML files, or CLI flags.

---

## 📦 Installation

### Prerequisites

- **Go 1.24+**
- **kubectl** access to a Kubernetes cluster (`~/.kube/config`).

### Build

```bash
go build ./cmd/kube-pod-update
```

---

## 📂 Project Structure

```text
kube-pod-update/
├── cmd/
│   └── kube-pod-update/    # CLI entrypoint
├── internal/
│   ├── k8s/                # Kubernetes API logic
│   ├── registry/           # Container registry interactions
│   ├── compare/            # Version/digest comparison logic
│   ├── notifier/           # Logging and notifications
│   └── config/             # Configuration management
├── config/                 # Optional configuration files
├── .gitignore
├── go.mod
└── README.md
```

---

## ✅ Development

### Run Tests

```bash
go test ./...
```

### Build Locally

```bash
go build ./cmd/kube-pod-update
```