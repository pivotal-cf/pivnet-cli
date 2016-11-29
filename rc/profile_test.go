package rc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/rc"
)

var _ = Describe("PivnetProfile", func() {
	var (
		profile rc.PivnetProfile
	)

	BeforeEach(func() {
		profile = rc.PivnetProfile{
			Name:     "some-name",
			APIToken: "some-api-token",
			Host:     "some-host",
		}
	})

	Describe("Validate", func() {
		It("returns nil when validation is successful", func() {
			err := profile.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when name is empty", func() {
			BeforeEach(func() {
				profile.Name = ""
			})

			It("returns an error", func() {
				err := profile.Validate()
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Name"))
			})
		})

		Context("when API token is empty", func() {
			BeforeEach(func() {
				profile.APIToken = ""
			})

			It("returns an error", func() {
				err := profile.Validate()
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("API token"))
			})
		})

		Context("when host is empty", func() {
			BeforeEach(func() {
				profile.Host = ""
			})

			It("returns an error", func() {
				err := profile.Validate()
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Host"))
			})
		})
	})
})
