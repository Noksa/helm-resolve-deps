package models

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HelmResolveDepsOptions", func() {
	Context("with zero-value defaults", func() {
		var opts HelmResolveDepsOptions

		BeforeEach(func() {
			opts = HelmResolveDepsOptions{}
		})

		It("should have Clean disabled", func() {
			Expect(opts.Clean).To(BeFalse())
		})

		It("should have Untar disabled", func() {
			Expect(opts.Untar).To(BeFalse())
		})

		It("should have SkipRefresh disabled", func() {
			Expect(opts.SkipRefresh).To(BeFalse())
		})

		It("should have zero Threads", func() {
			Expect(opts.Threads).To(BeZero())
		})

		It("should have empty SkipRefreshInCharts", func() {
			Expect(opts.SkipRefreshInCharts).To(BeEmpty())
		})

		It("should have empty Args", func() {
			Expect(opts.Args).To(BeEmpty())
		})
	})

	Context("with explicit values", func() {
		It("should retain all configured fields", func() {
			opts := HelmResolveDepsOptions{
				Clean:               true,
				Untar:               true,
				SkipRefresh:         true,
				SkipRefreshInCharts: []string{"chart1", "chart2"},
				Threads:             4,
				Args:                []string{"--debug", "--dry-run"},
			}

			Expect(opts.Clean).To(BeTrue())
			Expect(opts.Untar).To(BeTrue())
			Expect(opts.SkipRefresh).To(BeTrue())
			Expect(opts.Threads).To(Equal(4))
			Expect(opts.SkipRefreshInCharts).To(HaveLen(2))
			Expect(opts.Args).To(HaveLen(2))
		})
	})
})
