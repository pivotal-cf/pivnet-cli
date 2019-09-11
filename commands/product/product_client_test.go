package product_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/pivnet-cli/commands/product"
	"github.com/pivotal-cf/pivnet-cli/commands/product/productfakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("product commands", func() {
	var (
		fakePivnetClient *productfakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		products []pivnet.Product

		client *product.ProductClient
	)

	BeforeEach(func() {
		fakePivnetClient = &productfakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		products = []pivnet.Product{
			{
				ID: 1234,
			},
			{
				ID: 2345,
			},
		}

		client = product.NewProductClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("List", func() {
		BeforeEach(func() {
			fakePivnetClient.ProductsReturns(products, nil)
		})

		It("lists all Products", func() {
			err := client.List()
			Expect(err).NotTo(HaveOccurred())

			var returnedProducts []pivnet.Product
			err = json.Unmarshal(outBuffer.Bytes(), &returnedProducts)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedProducts).To(Equal(products))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("products error")
				fakePivnetClient.ProductsReturns(nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.List()
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Get", func() {
		var (
			productSlug string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"

			fakePivnetClient.FindProductForSlugReturns(products[0], nil)
		})

		It("gets Product", func() {
			err := client.Get(productSlug)
			Expect(err).NotTo(HaveOccurred())

			var returnedProduct pivnet.Product
			err = json.Unmarshal(outBuffer.Bytes(), &returnedProduct)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedProduct).To(Equal(products[0]))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("product error")
				fakePivnetClient.FindProductForSlugReturns(pivnet.Product{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(productSlug)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
