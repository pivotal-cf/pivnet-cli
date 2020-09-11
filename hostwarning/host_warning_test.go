package hostwarning_test

import (
	"github.com/pivotal-cf/pivnet-cli/v2/hostwarning"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Host Warning", func() {

	var (
		hw   *hostwarning.HostWarning
		host string
	)

	JustBeforeEach(func() {
		hw = hostwarning.NewHostWarning(host)
	})

	Describe("Warn", func() {

		Context("host is NOT production", func() {
			BeforeEach(func() {
				host = "http://localhost:3000"
			})
			It("returns a warning", func() {
				Expect(hw.Warn()).To(ContainSubstring("Warning: You are currently targeting http://localhost:3000"))
			})
		})

		Context("host is production", func() {
			BeforeEach(func() {
				host = "https://network.pivotal.io"
			})
			It("does not return a warning", func() {
				Expect(hw.Warn()).To(Equal(""))
			})
		})

		Context("host is empty", func() {
			BeforeEach(func() {
				host = ""
			})
			It("does not return a warning", func() {
				Expect(hw.Warn()).To(Equal(""))
			})
		})
	})
})
