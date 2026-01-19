package models

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMiniHelmChart_UnmarshalYAML(t *testing.T) {
	yamlData := `
name: my-chart
version: 1.2.3
dependencies:
  - name: nginx
    version: 1.0.0
    repository: https://charts.bitnami.com/bitnami
  - name: postgres
    version: 2.0.0
    repository: https://charts.bitnami.com/bitnami
    condition: postgres.enabled
    tags:
      - database
    enabled: true
    alias: db
`

	var chart MiniHelmChart
	err := yaml.Unmarshal([]byte(yamlData), &chart)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if chart.Name != "my-chart" {
		t.Errorf("expected name 'my-chart', got '%s'", chart.Name)
	}
	if chart.Version != "1.2.3" {
		t.Errorf("expected version '1.2.3', got '%s'", chart.Version)
	}
	if len(chart.Dependencies) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(chart.Dependencies))
	}

	dep1 := chart.Dependencies[0]
	if dep1.Name != "nginx" {
		t.Errorf("expected dependency name 'nginx', got '%s'", dep1.Name)
	}
	if dep1.Version != "1.0.0" {
		t.Errorf("expected dependency version '1.0.0', got '%s'", dep1.Version)
	}

	dep2 := chart.Dependencies[1]
	if dep2.Condition != "postgres.enabled" {
		t.Errorf("expected condition 'postgres.enabled', got '%s'", dep2.Condition)
	}
	if !dep2.Enabled {
		t.Error("expected enabled to be true")
	}
	if dep2.Alias != "db" {
		t.Errorf("expected alias 'db', got '%s'", dep2.Alias)
	}
	if len(dep2.Tags) != 1 || dep2.Tags[0] != "database" {
		t.Errorf("expected tags ['database'], got %v", dep2.Tags)
	}
}

func TestDependency_LocalRepository(t *testing.T) {
	yamlData := `
name: parent-chart
version: 1.0.0
dependencies:
  - name: local-chart
    version: 1.0.0
    repository: file://../local-chart
`

	var chart MiniHelmChart
	err := yaml.Unmarshal([]byte(yamlData), &chart)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(chart.Dependencies) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(chart.Dependencies))
	}

	dep := chart.Dependencies[0]
	if dep.Repository != "file://../local-chart" {
		t.Errorf("expected repository 'file://../local-chart', got '%s'", dep.Repository)
	}
}
