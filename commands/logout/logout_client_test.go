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
			apiToken    string
			host        string

			saveProfileErr error
		)

		BeforeEach(func() {
			profileName = "some-logout-slug"
			apiToken = "some-api-token"
			host = "some-host"

			saveProfileErr = nil
		})

		JustBeforeEach(func() {
			fakeRCHandler.SaveProfileReturns(saveProfileErr)
		})

		It("saves profile", func() {
			err := client.Logout(profileName)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeRCHandler.SaveProfileCallCount()).To(Equal(1))
			invokedProfileName, invokedAPIToken, invokedHost := fakeRCHandler.SaveProfileArgsForCall(0)

			Expect(invokedProfileName).To(Equal(profileName))
			Expect(invokedAPIToken).To(BeEmpty())
			Expect(invokedHost).To(BeEmpty())
		})

		Context("when saving profile returns an error", func() {
			BeforeEach(func() {
				saveProfileErr = fmt.Errorf("save profile error")
			})

			It("invokes the error handler", func() {
				err := client.Logout(profileName)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(saveProfileErr))
			})
		})
	})
})
