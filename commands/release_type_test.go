package commands_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
)

var _ = Describe("release types commands", func() {
	var (
		fakeReleaseTypeClient *commandsfakes.FakeReleaseTypeClient
	)

	BeforeEach(func() {
		fakeReleaseTypeClient = &commandsfakes.FakeReleaseTypeClient{}

		commands.NewReleaseTypeClient = func() commands.ReleaseTypeClient {
			return fakeReleaseTypeClient
		}
	})

	Describe("ReleaseTypesCommand", func() {
		var (
			cmd *commands.ReleaseTypesCommand
		)

		BeforeEach(func() {
			cmd = &commands.ReleaseTypesCommand{}
		})

		It("invokes the ReleaseType client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseTypeClient.ListCallCount()).To(Equal(1))
		})

		Context("when the ReleaseType client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeReleaseTypeClient.ListReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})
	})
})
