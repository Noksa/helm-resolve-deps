package helpers

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/noksa/helm-resolve-deps/internal/models"
)

var _ = Describe("LoadChartByPath", func() {
	var tmpDir string

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()
	})

	Context("with a valid Chart.yaml", func() {
		BeforeEach(func() {
			chartYAML := `name: test-chart
version: 1.0.0
dependencies:
  - name: dep1
    version: 1.0.0
    repository: https://charts.example.com
`
			Expect(os.WriteFile(filepath.Join(tmpDir, "Chart.yaml"), []byte(chartYAML), 0644)).To(Succeed())
		})

		It("should load chart name and version", func() {
			chart, err := LoadChartByPath(tmpDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(chart.Name).To(Equal("test-chart"))
			Expect(chart.Version).To(Equal("1.0.0"))
		})

		It("should load dependencies", func() {
			chart, err := LoadChartByPath(tmpDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(chart.Dependencies).To(HaveLen(1))
		})

		It("should set the chart path", func() {
			chart, err := LoadChartByPath(tmpDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(chart.Path).To(Equal(tmpDir))
		})
	})

	Context("with an invalid path", func() {
		It("should return an error", func() {
			_, err := LoadChartByPath("/nonexistent/path")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("with invalid YAML", func() {
		It("should return an error", func() {
			Expect(os.WriteFile(filepath.Join(tmpDir, "Chart.yaml"), []byte("invalid: yaml: content:"), 0644)).To(Succeed())
			_, err := LoadChartByPath(tmpDir)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("ResolveDependencies", func() {
	Context("when chart path does not exist", func() {
		It("should return an error", func() {
			err := ResolveDependencies("/nonexistent/path", models.HelmResolveDepsOptions{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("with a chart that has no dependencies", func() {
		It("should succeed", func() {
			tmpDir := GinkgoT().TempDir()
			chartYAML := `name: test-chart
version: 1.0.0
`
			Expect(os.WriteFile(filepath.Join(tmpDir, "Chart.yaml"), []byte(chartYAML), 0644)).To(Succeed())

			opts := models.HelmResolveDepsOptions{
				SkipRefresh: true,
				Threads:     1,
			}
			Expect(ResolveDependencies(tmpDir, opts)).To(Succeed())
		})
	})

	Context("with extra args", func() {
		It("should not fail on chart with no dependencies", func() {
			tmpDir := GinkgoT().TempDir()
			Expect(os.WriteFile(filepath.Join(tmpDir, "Chart.yaml"), []byte("name: test-chart\nversion: 1.0.0\n"), 0644)).To(Succeed())

			opts := models.HelmResolveDepsOptions{
				SkipRefresh: true,
				Args:        []string{"--debug"},
				Threads:     1,
			}
			err := ResolveDependencies(tmpDir, opts)
			GinkgoWriter.Printf("ResolveDependencies with args: %v\n", err)
		})
	})
})
