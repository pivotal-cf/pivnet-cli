package logout_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands/logout"
	"github.com/pivotal-cf/pivnet-cli/commands/logout/logoutfakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("logout commands", func() {
	var (
		fakeRCHandler    *logoutfakes.FakeRCHandler
		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		client *logout.LogoutClient
	)

	BeforeEach(func() {
		fakeRCHandler = &logoutfakes.FakeRCHandler{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		client = logout.NewLogoutClient(
			fakeRCHandler,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("Logout", func() {
		var (
			profileName string

			removeProfileErr error
		)

		BeforeEach(func() {
			profileName = "some-logout-slug"

			removeProfileErr = nil
		})

		JustBeforeEach(func() {
			fakeRCHandler.RemoveProfileWithNameReturns(removeProfileErr)
		})

		It("removes profile", func() {
			err := client.Logout(profileName)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeRCHandler.RemoveProfileWithNameCallCount()).To(Equal(1))
			invokedProfileName := fakeRCHandler.RemoveProfileWithNameArgsForCall(0)

			Expect(invokedProfileName).To(Equal(profileName))
		})

		Context("when removing profile returns an error", func() {
			BeforeEach(func() {
				removeProfileErr = fmt.Errorf("remove profile error")
			})

			It("invokes the error handler", func() {
				err := client.Logout(profileName)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(removeProfileErr))
			})
		})
	})
})
