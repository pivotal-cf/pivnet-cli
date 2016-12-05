package commands_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
)

var _ = Describe("logout commands", func() {
	var (
		fakeLogoutClient *commandsfakes.FakeLogoutClient
	)

	BeforeEach(func() {
		fakeLogoutClient = &commandsfakes.FakeLogoutClient{}

		commands.NewLogoutClient = func() commands.LogoutClient {
			return fakeLogoutClient
		}
	})

	Describe("LogoutCommand", func() {
		var (
			cmd commands.LogoutCommand
		)

		It("invokes the Logout client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeLogoutClient.LogoutCallCount()).To(Equal(1))
		})

		It("invokes the Init function with 'true'", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(initInvocationArg).To(BeTrue())
		})

		Context("when the Logout client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeLogoutClient.LogoutReturns(expectedErr)
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
	})
})
