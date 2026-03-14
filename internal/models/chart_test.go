package models

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("MiniHelmChart", func() {
	Context("when unmarshalling a full Chart.yaml", func() {
		var chart MiniHelmChart

		BeforeEach(func() {
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
			Expect(yaml.Unmarshal([]byte(yamlData), &chart)).To(Succeed())
		})

		It("should parse chart name and version", func() {
			Expect(chart.Name).To(Equal("my-chart"))
			Expect(chart.Version).To(Equal("1.2.3"))
		})

		It("should parse all dependencies", func() {
			Expect(chart.Dependencies).To(HaveLen(2))
		})

		It("should parse the first dependency correctly", func() {
			dep := chart.Dependencies[0]
			Expect(dep.Name).To(Equal("nginx"))
			Expect(dep.Version).To(Equal("1.0.0"))
		})

		It("should parse dependency metadata", func() {
			dep := chart.Dependencies[1]
			Expect(dep.Condition).To(Equal("postgres.enabled"))
			Expect(dep.Enabled).To(BeTrue())
			Expect(dep.Alias).To(Equal("db"))
			Expect(dep.Tags).To(ConsistOf("database"))
		})
	})
})

var _ = Describe("Dependency", func() {
	Context("with a local file:// repository", func() {
		It("should preserve the repository path", func() {
			yamlData := `
name: parent-chart
version: 1.0.0
dependencies:
  - name: local-chart
    version: 1.0.0
    repository: file://../local-chart
`
			var chart MiniHelmChart
			Expect(yaml.Unmarshal([]byte(yamlData), &chart)).To(Succeed())
			Expect(chart.Dependencies).To(HaveLen(1))
			Expect(chart.Dependencies[0].Repository).To(Equal("file://../local-chart"))
		})
	})
})
