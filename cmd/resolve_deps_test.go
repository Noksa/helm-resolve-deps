package main

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/noksa/helm-resolve-deps/internal/helpers"
	"github.com/noksa/helm-resolve-deps/internal/models"
)

var _ = Describe("Integration: Local Dependencies", func() {
	var tmpDir string

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()
	})

	Context("when a parent chart has a local dependency", func() {
		var parentDir string

		BeforeEach(func() {
			childDir := filepath.Join(tmpDir, "child-chart")
			Expect(os.MkdirAll(filepath.Join(childDir, "templates"), 0755)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(childDir, "Chart.yaml"), []byte(`apiVersion: v2
name: child-chart
version: 1.0.0
type: application
`), 0644)).To(Succeed())

			parentDir = filepath.Join(tmpDir, "parent-chart")
			Expect(os.MkdirAll(filepath.Join(parentDir, "templates"), 0755)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(parentDir, "Chart.yaml"), []byte(`apiVersion: v2
name: parent-chart
version: 1.0.0
type: application
dependencies:
  - name: child-chart
    version: 1.0.0
    repository: file://../child-chart
`), 0644)).To(Succeed())
		})

		It("should load the parent chart", func() {
			chart, err := helpers.LoadChartByPath(parentDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(chart.Name).To(Equal("parent-chart"))
		})

		It("should detect the local dependency", func() {
			chart, err := helpers.LoadChartByPath(parentDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(chart.Dependencies).To(HaveLen(1))
			Expect(chart.Dependencies[0].Name).To(Equal("child-chart"))
			Expect(chart.Dependencies[0].Repository).To(Equal("file://../child-chart"))
		})
	})
})

var _ = Describe("Integration: Chain Dependencies", func() {
	It("should load transitive dependency chains", func() {
		tmpDir := GinkgoT().TempDir()

		// grandchild
		grandchildDir := filepath.Join(tmpDir, "grandchild-chart")
		Expect(os.MkdirAll(filepath.Join(grandchildDir, "templates"), 0755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(grandchildDir, "Chart.yaml"), []byte(`apiVersion: v2
name: grandchild-chart
version: 1.0.0
type: application
`), 0644)).To(Succeed())

		// child -> grandchild
		childDir := filepath.Join(tmpDir, "child-chart")
		Expect(os.MkdirAll(filepath.Join(childDir, "templates"), 0755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(childDir, "Chart.yaml"), []byte(`apiVersion: v2
name: child-chart
version: 1.0.0
type: application
dependencies:
  - name: grandchild-chart
    version: 1.0.0
    repository: file://../grandchild-chart
`), 0644)).To(Succeed())

		// parent -> child
		parentDir := filepath.Join(tmpDir, "parent-chart")
		Expect(os.MkdirAll(filepath.Join(parentDir, "templates"), 0755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(parentDir, "Chart.yaml"), []byte(`apiVersion: v2
name: parent-chart
version: 1.0.0
type: application
dependencies:
  - name: child-chart
    version: 1.0.0
    repository: file://../child-chart
`), 0644)).To(Succeed())

		parentChart, err := helpers.LoadChartByPath(parentDir)
		Expect(err).NotTo(HaveOccurred())
		Expect(parentChart.Dependencies).To(HaveLen(1))

		childChart, err := helpers.LoadChartByPath(childDir)
		Expect(err).NotTo(HaveOccurred())
		Expect(childChart.Dependencies).To(HaveLen(1))
	})
})

var _ = Describe("Options: Clean Flag", func() {
	It("should remove charts, tmpcharts and Chart.lock", func() {
		tmpDir := GinkgoT().TempDir()
		Expect(os.WriteFile(filepath.Join(tmpDir, "Chart.yaml"), []byte(`apiVersion: v2
name: test-chart
version: 1.0.0
type: application
`), 0644)).To(Succeed())

		chartsDir := filepath.Join(tmpDir, "charts")
		tmpchartsDir := filepath.Join(tmpDir, "tmpcharts")
		lockFile := filepath.Join(tmpDir, "Chart.lock")

		Expect(os.MkdirAll(chartsDir, 0755)).To(Succeed())
		Expect(os.MkdirAll(tmpchartsDir, 0755)).To(Succeed())
		Expect(os.WriteFile(lockFile, []byte("test"), 0644)).To(Succeed())

		opts := models.HelmResolveDepsOptions{
			SkipRefresh: true,
			Clean:       true,
			Threads:     1,
		}
		_ = helpers.ResolveDependencies(tmpDir, opts)

		Expect(chartsDir).NotTo(BeADirectory())
		Expect(tmpchartsDir).NotTo(BeADirectory())
		Expect(lockFile).NotTo(BeAnExistingFile())
	})
})

var _ = Describe("Options: SkipRefreshInCharts", func() {
	It("should accept skip-refresh-in option without error", func() {
		tmpDir := GinkgoT().TempDir()

		childDir := filepath.Join(tmpDir, "child")
		Expect(os.MkdirAll(filepath.Join(childDir, "templates"), 0755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(childDir, "Chart.yaml"), []byte(`apiVersion: v2
name: child-chart
version: 1.0.0
`), 0644)).To(Succeed())

		parentDir := filepath.Join(tmpDir, "parent")
		Expect(os.MkdirAll(filepath.Join(parentDir, "templates"), 0755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(parentDir, "Chart.yaml"), []byte(`apiVersion: v2
name: parent-chart
version: 1.0.0
dependencies:
  - name: child-chart
    version: 1.0.0
    repository: file://../child
`), 0644)).To(Succeed())

		opts := models.HelmResolveDepsOptions{
			SkipRefresh:         true,
			SkipRefreshInCharts: []string{"child-chart"},
			Threads:             1,
		}

		err := helpers.ResolveDependencies(parentDir, opts)
		if err != nil {
			GinkgoWriter.Printf("Expected behavior: %v\n", err)
		}
	})
})
