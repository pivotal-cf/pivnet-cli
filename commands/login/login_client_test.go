package login_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands/login"
	"github.com/pivotal-cf/pivnet-cli/commands/login/loginfakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("login commands", func() {
	var (
		fakePivnetClient *loginfakes.FakePivnetClient
		fakeRCHandler    *loginfakes.FakeRCHandler
		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		client *login.LoginClient
	)

	BeforeEach(func() {
		fakePivnetClient = &loginfakes.FakePivnetClient{}
		fakeRCHandler = &loginfakes.FakeRCHandler{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		client = login.NewLoginClient(
			fakePivnetClient,
			fakeRCHandler,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("Login", func() {
		var (
			profileName string
			apiToken    string

			authResult bool
			authErr    error

			saveProfileErr error
		)

		BeforeEach(func() {
			profileName = "some-login-slug"
			apiToken = "some-api-token"

			authResult = true
			authErr = nil

			saveProfileErr = nil
		})

		JustBeforeEach(func() {
			fakePivnetClient.AuthReturns(authResult, authErr)
			fakeRCHandler.SaveProfileReturns(saveProfileErr)
		})

		It("authenticates and saves profile", func() {
			err := client.Login(profileName, apiToken)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePivnetClient.AuthCallCount()).To(Equal(1))

			Expect(fakeRCHandler.SaveProfileCallCount()).To(Equal(1))
			invokedProfileName, invokedAPIToken := fakeRCHandler.SaveProfileArgsForCall(0)

			Expect(invokedProfileName).To(Equal(profileName))
			Expect(invokedAPIToken).To(Equal(apiToken))
		})

		Context("when there is an error authenticating", func() {
			BeforeEach(func() {
				authErr = fmt.Errorf("auth error")
			})

			It("invokes the error handler", func() {
				err := client.Login(profileName, apiToken)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(authErr))
			})
		})

		Context("when authenticating fails", func() {
			BeforeEach(func() {
				authResult = false
			})

			It("invokes the error handler", func() {
				err := client.Login(profileName, apiToken)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0).Error()).To(ContainSubstring("login"))
			})
		})

		Context("when saving profile returns an error", func() {
			BeforeEach(func() {
				saveProfileErr = fmt.Errorf("save profile error")
			})

			It("invokes the error handler", func() {
				err := client.Login(profileName, apiToken)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(saveProfileErr))
			})
		})
	})
})
