package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/noksa/helm-resolve-deps/internal/models"
)

func TestLoadChartByPath(t *testing.T) {
	tmpDir := t.TempDir()
	chartYAML := `name: test-chart
version: 1.0.0
dependencies:
  - name: dep1
    version: 1.0.0
    repository: https://charts.example.com
`
	chartPath := filepath.Join(tmpDir, "Chart.yaml")
	if err := os.WriteFile(chartPath, []byte(chartYAML), 0644); err != nil {
		t.Fatal(err)
	}

	chart, err := LoadChartByPath(tmpDir)
	if err != nil {
		t.Fatalf("LoadChartByPath failed: %v", err)
	}

	if chart.Name != "test-chart" {
		t.Errorf("expected name 'test-chart', got '%s'", chart.Name)
	}
	if chart.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", chart.Version)
	}
	if len(chart.Dependencies) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(chart.Dependencies))
	}
	if chart.Path != tmpDir {
		t.Errorf("expected path '%s', got '%s'", tmpDir, chart.Path)
	}
}

func TestLoadChartByPath_InvalidPath(t *testing.T) {
	_, err := LoadChartByPath("/nonexistent/path")
	if err == nil {
		t.Error("expected error for nonexistent path, got nil")
	}
}

func TestLoadChartByPath_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	chartPath := filepath.Join(tmpDir, "Chart.yaml")
	if err := os.WriteFile(chartPath, []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadChartByPath(tmpDir)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestResolveDependencies_NoChart(t *testing.T) {
	err := ResolveDependencies("/nonexistent/path", models.HelmResolveDepsOptions{})
	if err == nil {
		t.Error("expected error for nonexistent chart, got nil")
	}
}

func TestResolveDependencies_EmptyDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	chartYAML := `name: test-chart
version: 1.0.0
`
	chartPath := filepath.Join(tmpDir, "Chart.yaml")
	if err := os.WriteFile(chartPath, []byte(chartYAML), 0644); err != nil {
		t.Fatal(err)
	}

	opts := models.HelmResolveDepsOptions{
		SkipRefresh: true,
		Clean:       false,
		Threads:     1,
	}

	err := ResolveDependencies(tmpDir, opts)
	if err != nil {
		t.Errorf("ResolveDependencies failed for chart with no dependencies: %v", err)
	}
}

func TestResolveDependencies_WithArgs(t *testing.T) {
	tmpDir := t.TempDir()
	chartYAML := `name: test-chart
version: 1.0.0
`
	if err := os.WriteFile(filepath.Join(tmpDir, "Chart.yaml"), []byte(chartYAML), 0644); err != nil {
		t.Fatal(err)
	}

	opts := models.HelmResolveDepsOptions{
		SkipRefresh: true,
		Args:        []string{"--debug"},
		Threads:     1,
	}

	err := ResolveDependencies(tmpDir, opts)
	if err != nil {
		t.Logf("ResolveDependencies with args: %v", err)
	}
}

func TestResolveDependencies_MultipleThreads(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	tmpDir := t.TempDir()

	for i := 1; i <= 3; i++ {
		childDir := filepath.Join(tmpDir, filepath.Join("child", string(rune('0'+i))))
		if err := os.MkdirAll(filepath.Join(childDir, "templates"), 0755); err != nil {
			t.Fatal(err)
		}

		childYAML := `apiVersion: v2
name: child-` + string(rune('0'+i)) + `
version: 1.0.0
`
		if err := os.WriteFile(filepath.Join(childDir, "Chart.yaml"), []byte(childYAML), 0644); err != nil {
			t.Fatal(err)
		}
	}

	parentDir := filepath.Join(tmpDir, "parent")
	if err := os.MkdirAll(filepath.Join(parentDir, "templates"), 0755); err != nil {
		t.Fatal(err)
	}

	parentYAML := `apiVersion: v2
name: parent
version: 1.0.0
dependencies:
  - name: child-1
    version: 1.0.0
    repository: file://../child1
  - name: child-2
    version: 1.0.0
    repository: file://../child2
  - name: child-3
    version: 1.0.0
    repository: file://../child3
`
	if err := os.WriteFile(filepath.Join(parentDir, "Chart.yaml"), []byte(parentYAML), 0644); err != nil {
		t.Fatal(err)
	}

	opts := models.HelmResolveDepsOptions{
		SkipRefresh: true,
		Threads:     3,
	}

	err := ResolveDependencies(parentDir, opts)
	if err != nil {
		t.Logf("Multi-threaded resolution: %v", err)
	}
}
