# 🔗 helm-resolve-deps

> ⚡ A Helm plugin that properly resolves local chain dependencies in charts

---

## 🤔 Why?

Helm's built-in `helm dependency update` doesn't handle transitive local dependencies well. If you have charts with `file://` dependencies that themselves have dependencies, Helm won't resolve the entire chain correctly.

**helm-resolve-deps** solves this by recursively resolving all local dependencies in the correct order.

Inspired by [helm/helm#2247](https://github.com/helm/helm/issues/2247)

> 📝 **Note**: I've created [an issue](https://github.com/helm/helm/issues/31496) and [a pull request](https://github.com/helm/helm/pull/31524) to fix this in Helm itself. Until it's merged, this plugin provides the solution.

### ✨ Key Benefits

| Feature                        | Description                                              |
|--------------------------------|----------------------------------------------------------|
| 🔄 **Chain Resolution**        | Recursively resolves local (`file://`) dependencies      |
| ⚡ **Parallel Processing**     | Multi-threaded dependency resolution                     |
| 📦 **Unpack Mode**             | Extract `.tgz` archives for debugging                    |
| 🧹 **Clean Mode**              | Remove old dependencies before updating                  |
| 🎯 **Selective Refresh**       | Skip repository updates for specific charts              |

---

## 📋 Requirements

- 🎯 **Helm 3** installed on host machine

---

## 🚀 Installation

```bash
(helm plugin uninstall resolve-deps || true) && helm plugin install https://github.com/Noksa/helm-resolve-deps.git
```

> 💡 **Helm 4 users**: Add `--verify=false` flag to the install command

---

## 📖 Usage

### 🔍 Getting Help

```bash
helm resolve-deps --help
```

### 🎯 Basic Syntax

```bash
helm resolve-deps [PATH] [FLAGS]
```

### 🔧 Available Flags

| Flag                  | Short | Description                                              |
|-----------------------|-------|----------------------------------------------------------|
| `--untar`             | `-u`  | Unpack dependent charts as directories instead of `.tgz` |
| `--clean`             | `-c`  | Remove `charts/`, `tmpcharts/`, and `Chart.lock`         |
| `--skip-refresh`      |       | Skip fetching updates from Helm repositories             |
| `--skip-refresh-in`   |       | Skip refresh for specific charts (comma-separated)       |
| `--threads`           |       | Number of parallel workers (default: CPU count - 1)      |
| `--help`              | `-h`  | Show help                                                |

---

## ⚙️ How It Works

Given this structure:
```
parent-chart/
├── Chart.yaml (depends on child-chart via file://)
└── charts/
    └── child-chart/
        └── Chart.yaml (depends on grandchild-chart via file://)
```

Running `helm resolve-deps parent-chart` will:
1. 🔍 **Discover** all local dependencies recursively
2. 📦 **Resolve** grandchild-chart dependencies first
3. 🔗 **Resolve** child-chart dependencies (including grandchild)
4. ✅ **Resolve** parent-chart dependencies (including entire chain)

---

## 💡 Examples

### 🔍 Basic Operations

**Resolve dependencies in current directory**
```bash
helm resolve-deps .
```

**Resolve with repository refresh skipped**
```bash
helm resolve-deps . --skip-refresh
```

### 🧹 Clean Operations

**Clean before resolving**
```bash
helm resolve-deps . --clean
```

**Clean and unpack dependencies**
```bash
helm resolve-deps . --clean --untar
```

### ⚡ Performance Optimization

**Use multiple threads**
```bash
# Use 4 parallel workers
helm resolve-deps . --threads 4
```

**Skip refresh for specific charts**
```bash
# Skip refresh for chart1 and chart2
helm resolve-deps . --skip-refresh-in chart1,chart2
```

### 🔧 Advanced Usage

**Unpack dependencies for debugging**
```bash
# Extract all .tgz files to directories
helm resolve-deps ~/charts/my-chart --untar
```

> 💡 **Tip**: Use `--untar` to inspect and modify dependent charts directly in the `charts/` directory

**Pass additional flags to helm dependency update**
```bash
# Pass flags after --
helm resolve-deps . -- --kubeconfig myconfig
```

**Full example with all options**
```bash
helm resolve-deps ~/charts/my-chart \
  --clean \
  --untar \
  --skip-refresh-in my-chart1,my-chart2 \
  --threads 4
```

---

## 🧪 Testing

```bash
make test
```

---

## 🔍 Linting

```bash
make lint
```
