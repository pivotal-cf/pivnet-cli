package commands_test

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/commands/curl"
)

var _ = Describe("curl commands", func() {
	var (
		field reflect.StructField

		fakeCurlClient *commandsfakes.FakeCurlClient
	)

	BeforeEach(func() {
		fakeCurlClient = &commandsfakes.FakeCurlClient{}

		commands.NewCurlClient = func(curl.PivnetClient) commands.CurlClient {
			return fakeCurlClient
		}
	})

	Describe("CurlCommand", func() {
		var (
			cmd *commands.CurlCommand
		)

		BeforeEach(func() {
			cmd = &commands.CurlCommand{}
		})

		It("invokes the curl client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeCurlClient.MakeRequestCallCount()).To(Equal(1))
		})

		Context("when the curl client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeCurlClient.MakeRequestReturns(expectedErr)
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

		Describe("Method flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.CurlCommand{}, "Method")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("request"))
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("X"))
			})
		})

		Describe("Data flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.CurlCommand{}, "Data")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("data"))
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("d"))
			})
		})
	})
})
