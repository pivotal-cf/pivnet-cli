package curl_test

import (
	"bytes"
	"errors"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands/curl"
	"github.com/pivotal-cf/pivnet-cli/commands/curl/curlfakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("curl commands", func() {
	var (
		fakePivnetClient *curlfakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		client *curl.CurlClient
	)

	BeforeEach(func() {
		fakePivnetClient = &curlfakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		client = curl.NewCurlClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("MakeRequest", func() {
		var (
			method string
			data   string
			args   []string
		)

		BeforeEach(func() {
			method = "some-method"
			data = "some-data"
			args = []string{"some-endpoint", "not-used"}
		})

		It("invokes the client", func() {
			err := client.MakeRequest(method, data, args)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePivnetClient.MakeRequestCallCount()).To(Equal(1))

			invokedMethod,
				invokedEndpoint,
				invokedExpectedResponseCode,
				invokedBody,
				invokedResponse := fakePivnetClient.MakeRequestArgsForCall(0)

			Expect(invokedMethod).To(Equal(method))
			Expect(invokedEndpoint).To(Equal(args[0]))
			Expect(invokedExpectedResponseCode).To(Equal(0))
			Expect(invokedBody).To(Equal(strings.NewReader(data)))
			Expect(invokedResponse).To(BeNil())
		})

		Context("when data is empty", func() {
			BeforeEach(func() {
				data = ""
			})

			It("invokes the client with nil body", func() {
				err := client.MakeRequest(method, data, args)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.MakeRequestCallCount()).To(Equal(1))

				_,
					_,
					_,
					invokedBody,
					_ := fakePivnetClient.MakeRequestArgsForCall(0)

				Expect(invokedBody).To(BeNil())
			})
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("curl error")
				fakePivnetClient.MakeRequestReturns(nil, nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.MakeRequest(method, data, args)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error unmarshalling json", func() {
			var (
				returnedBytes []byte
			)

			BeforeEach(func() {
				returnedBytes = []byte(`[garbage-json!`)
				fakePivnetClient.MakeRequestReturns(nil, returnedBytes, nil)
			})

			It("invokes the error handler", func() {
				err := client.MakeRequest(method, data, args)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(HaveOccurred())
			})
		})
	})
})
