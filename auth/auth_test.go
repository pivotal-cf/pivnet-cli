package auth_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/auth"
	"github.com/pivotal-cf/pivnet-cli/auth/authfakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
)

var _ = Describe("Authenticator", func() {
	var (
		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		authenticator *auth.Authenticator
	)

	BeforeEach(func() {
		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}
		authenticator = auth.NewAuthenticator(fakeErrorHandler)
	})

	Describe("AuthenticateClient", func() {
		var (
			fakeAuthClient *authfakes.FakeAuthClient

			authOk    bool
			authError error
		)

		BeforeEach(func() {
			fakeAuthClient = &authfakes.FakeAuthClient{}

			authOk = true
			authError = nil
		})

		JustBeforeEach(func() {
			fakeAuthClient.AuthReturns(authOk, authError)
		})

		It("invokes the provided client", func() {
			err := authenticator.AuthenticateClient(fakeAuthClient)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeAuthClient.AuthCallCount()).To(Equal(1))
		})

		Context("when client auth returns an error", func() {
			BeforeEach(func() {
				authError = fmt.Errorf("auth error")
			})

			It("handles the error", func() {
				err := authenticator.AuthenticateClient(fakeAuthClient)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(authError))
			})
		})

		Context("when client auth returns not ok", func() {
			BeforeEach(func() {
				authOk = false
			})

			It("handles an error", func() {
				err := authenticator.AuthenticateClient(fakeAuthClient)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0).Error()).To(ContainSubstring("login again"))
			})
		})
	})
})
