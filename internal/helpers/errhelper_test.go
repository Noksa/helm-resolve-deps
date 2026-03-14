package helpers

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Must", func() {
	Context("when error is nil", func() {
		It("should not panic", func() {
			Expect(func() { Must(nil) }).NotTo(Panic())
		})
	})

	// Must calls os.Exit(1) on error — cannot be directly tested in-process
})
