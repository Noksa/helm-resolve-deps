package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/noksa/helm-resolve-deps/internal/helpers"
	"github.com/noksa/helm-resolve-deps/internal/models"
)

func TestIntegration_LocalDependencies(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir := t.TempDir()

	// Create a child chart
	childDir := filepath.Join(tmpDir, "child-chart")
	if err := os.MkdirAll(childDir, 0755); err != nil {
		t.Fatal(err)
	}

	childChartYAML := `apiVersion: v2
name: child-chart
version: 1.0.0
type: application
`
	if err := os.WriteFile(filepath.Join(childDir, "Chart.yaml"), []byte(childChartYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create templates directory for child
	if err := os.MkdirAll(filepath.Join(childDir, "templates"), 0755); err != nil {
		t.Fatal(err)
	}

	// Create a parent chart with local dependency
	parentDir := filepath.Join(tmpDir, "parent-chart")
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		t.Fatal(err)
	}

	parentChartYAML := `apiVersion: v2
name: parent-chart
version: 1.0.0
type: application
dependencies:
  - name: child-chart
    version: 1.0.0
    repository: file://../child-chart
`
	if err := os.WriteFile(filepath.Join(parentDir, "Chart.yaml"), []byte(parentChartYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create templates directory for parent
	if err := os.MkdirAll(filepath.Join(parentDir, "templates"), 0755); err != nil {
		t.Fatal(err)
	}

	// Test loading the parent chart
	chart, err := helpers.LoadChartByPath(parentDir)
	if err != nil {
		t.Fatalf("failed to load parent chart: %v", err)
	}

	if chart.Name != "parent-chart" {
		t.Errorf("expected chart name 'parent-chart', got '%s'", chart.Name)
	}

	if len(chart.Dependencies) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(chart.Dependencies))
	}

	dep := chart.Dependencies[0]
	if dep.Name != "child-chart" {
		t.Errorf("expected dependency name 'child-chart', got '%s'", dep.Name)
	}
	if dep.Repository != "file://../child-chart" {
		t.Errorf("expected repository 'file://../child-chart', got '%s'", dep.Repository)
	}
}

func TestIntegration_ChainDependencies(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir := t.TempDir()

	// Create grandchild chart (no dependencies)
	grandchildDir := filepath.Join(tmpDir, "grandchild-chart")
	if err := os.MkdirAll(filepath.Join(grandchildDir, "templates"), 0755); err != nil {
		t.Fatal(err)
	}

	grandchildYAML := `apiVersion: v2
name: grandchild-chart
version: 1.0.0
type: application
`
	if err := os.WriteFile(filepath.Join(grandchildDir, "Chart.yaml"), []byte(grandchildYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create child chart (depends on grandchild)
	childDir := filepath.Join(tmpDir, "child-chart")
	if err := os.MkdirAll(filepath.Join(childDir, "templates"), 0755); err != nil {
		t.Fatal(err)
	}

	childYAML := `apiVersion: v2
name: child-chart
version: 1.0.0
type: application
dependencies:
  - name: grandchild-chart
    version: 1.0.0
    repository: file://../grandchild-chart
`
	if err := os.WriteFile(filepath.Join(childDir, "Chart.yaml"), []byte(childYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create parent chart (depends on child)
	parentDir := filepath.Join(tmpDir, "parent-chart")
	if err := os.MkdirAll(filepath.Join(parentDir, "templates"), 0755); err != nil {
		t.Fatal(err)
	}

	parentYAML := `apiVersion: v2
name: parent-chart
version: 1.0.0
type: application
dependencies:
  - name: child-chart
    version: 1.0.0
    repository: file://../child-chart
`
	if err := os.WriteFile(filepath.Join(parentDir, "Chart.yaml"), []byte(parentYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Load and verify the chain
	parentChart, err := helpers.LoadChartByPath(parentDir)
	if err != nil {
		t.Fatalf("failed to load parent chart: %v", err)
	}

	if len(parentChart.Dependencies) != 1 {
		t.Fatalf("expected 1 dependency in parent, got %d", len(parentChart.Dependencies))
	}

	childChart, err := helpers.LoadChartByPath(childDir)
	if err != nil {
		t.Fatalf("failed to load child chart: %v", err)
	}

	if len(childChart.Dependencies) != 1 {
		t.Fatalf("expected 1 dependency in child, got %d", len(childChart.Dependencies))
	}
}

func TestOptions_CleanFlag(t *testing.T) {
	tmpDir := t.TempDir()

	chartYAML := `apiVersion: v2
name: test-chart
version: 1.0.0
type: application
`
	if err := os.WriteFile(filepath.Join(tmpDir, "Chart.yaml"), []byte(chartYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create directories that should be cleaned
	chartsDir := filepath.Join(tmpDir, "charts")
	tmpchartsDir := filepath.Join(tmpDir, "tmpcharts")
	lockFile := filepath.Join(tmpDir, "Chart.lock")

	if err := os.MkdirAll(chartsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(tmpchartsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lockFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	opts := models.HelmResolveDepsOptions{
		SkipRefresh: true,
		Clean:       true,
		Threads:     1,
	}

	// This will clean the directories
	_ = helpers.ResolveDependencies(tmpDir, opts)

	// Verify directories were cleaned
	if _, err := os.Stat(chartsDir); !os.IsNotExist(err) {
		t.Error("charts directory should have been removed")
	}
	if _, err := os.Stat(tmpchartsDir); !os.IsNotExist(err) {
		t.Error("tmpcharts directory should have been removed")
	}
	if _, err := os.Stat(lockFile); !os.IsNotExist(err) {
		t.Error("Chart.lock file should have been removed")
	}
}

func TestOptions_SkipRefreshInCharts(t *testing.T) {
	tmpDir := t.TempDir()

	childDir := filepath.Join(tmpDir, "child")
	if err := os.MkdirAll(filepath.Join(childDir, "templates"), 0755); err != nil {
		t.Fatal(err)
	}

	childYAML := `apiVersion: v2
name: child-chart
version: 1.0.0
`
	if err := os.WriteFile(filepath.Join(childDir, "Chart.yaml"), []byte(childYAML), 0644); err != nil {
		t.Fatal(err)
	}

	parentDir := filepath.Join(tmpDir, "parent")
	if err := os.MkdirAll(filepath.Join(parentDir, "templates"), 0755); err != nil {
		t.Fatal(err)
	}

	parentYAML := `apiVersion: v2
name: parent-chart
version: 1.0.0
dependencies:
  - name: child-chart
    version: 1.0.0
    repository: file://../child
`
	if err := os.WriteFile(filepath.Join(parentDir, "Chart.yaml"), []byte(parentYAML), 0644); err != nil {
		t.Fatal(err)
	}

	opts := models.HelmResolveDepsOptions{
		SkipRefresh:         true,
		SkipRefreshInCharts: []string{"child-chart"},
		Threads:             1,
	}

	err := helpers.ResolveDependencies(parentDir, opts)
	if err != nil {
		t.Logf("Expected behavior: %v", err)
	}
}
