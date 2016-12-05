package commands_test

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/commands/login"
)

var _ = Describe("login commands", func() {
	var (
		field reflect.StructField

		fakeLoginClient *commandsfakes.FakeLoginClient
	)

	BeforeEach(func() {
		fakeLoginClient = &commandsfakes.FakeLoginClient{}

		commands.NewLoginClient = func(login.PivnetClient) commands.LoginClient {
			return fakeLoginClient
		}
	})

	Describe("LoginCommand", func() {
		var (
			cmd commands.LoginCommand
		)

		BeforeEach(func() {
			cmd.APIToken = "some-api-token"
		})

		It("invokes the Login client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeLoginClient.LoginCallCount()).To(Equal(1))
		})

		It("invokes the Init function with 'false'", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(initInvocationArg).To(BeFalse())
		})

		Context("when the Login client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeLoginClient.LoginReturns(expectedErr)
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

		It("sanitizes new api token", func() {
			outBuffer := bytes.Buffer{}
			commands.OutputWriter = &outBuffer

			err := cmd.Execute(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = fmt.Fprintf(commands.OutputWriter, apiToken)
			Expect(err).NotTo(HaveOccurred())

			Expect(outBuffer.String()).Should(ContainSubstring("*** redacted api token ***"))
			Expect(outBuffer.String()).ShouldNot(ContainSubstring(apiToken))
		})

		Describe("APIToken flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "APIToken")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("api-token"))
			})
		})

		Describe("Host flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Host")
			})

			It("contains long flag", func() {
				Expect(longTag(field)).To(Equal("host"))
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("has a default value", func() {
				Expect(defaultVal(field)).To(Equal("https://network.pivotal.io"))
			})
		})

	})
})
