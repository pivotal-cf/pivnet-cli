package commands_test

import (
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
)

var _ = Describe("product commands", func() {
	var (
		field reflect.StructField

		fakeProductClient *commandsfakes.FakeProductClient
	)

	BeforeEach(func() {
		fakeProductClient = &commandsfakes.FakeProductClient{}

		commands.NewProductClient = func() commands.ProductClient {
			return fakeProductClient
		}
	})

	Describe("ProductsCommand", func() {
		var (
			cmd commands.ProductsCommand
		)

		BeforeEach(func() {
			cmd = commands.ProductsCommand{}
		})

		It("invokes the Product client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeProductClient.ListCallCount()).To(Equal(1))
		})

		Context("when the Product client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeProductClient.ListReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})
	})

	Describe("ProductCommand", func() {
		var (
			cmd commands.ProductCommand

			productSlug string
		)

		BeforeEach(func() {
			productSlug = "some product slug"

			cmd = commands.ProductCommand{
				ProductSlug: productSlug,
			}
		})

		It("invokes the Product client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeProductClient.GetCallCount()).To(Equal(1))
		})

		Context("when the Product client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeProductClient.GetReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ProductCommand{}, "ProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("p"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-slug"))
			})
		})
	})
})
