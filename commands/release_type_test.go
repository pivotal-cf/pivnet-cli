package commands_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/commands/releasetype"
)

var _ = Describe("release types commands", func() {
	var (
		fakeReleaseTypeClient *commandsfakes.FakeReleaseTypeClient
	)

	BeforeEach(func() {
		fakeReleaseTypeClient = &commandsfakes.FakeReleaseTypeClient{}

		commands.NewReleaseTypeClient = func(releasetype.PivnetClient) commands.ReleaseTypeClient {
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

		Context("when Init returns an error", func() {
			BeforeEach(func() {
				initErr = fmt.Errorf("init error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(initErr))
			})
		})

		Context("when Authentication returns an error", func() {
			BeforeEach(func() {
				authErr = fmt.Errorf("auth error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(authErr))
			})
		})

	})
})
