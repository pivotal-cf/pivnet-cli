package curl_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
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

		resp           *http.Response
		makeRequestErr error

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

		resp = &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("")),
		}
		makeRequestErr = nil
	})

	JustBeforeEach(func() {
		fakePivnetClient.MakeRequestReturns(resp, makeRequestErr)
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

			resp = &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader([]byte{})),
			}
		})

		JustBeforeEach(func() {
			fakePivnetClient.MakeRequestReturns(resp, makeRequestErr)
		})

		It("invokes the client", func() {
			err := client.MakeRequest(method, data, args)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePivnetClient.MakeRequestCallCount()).To(Equal(1))

			invokedMethod,
				invokedEndpoint,
				invokedExpectedResponseCode,
				invokedBody := fakePivnetClient.MakeRequestArgsForCall(0)

			Expect(invokedMethod).To(Equal(method))
			Expect(invokedEndpoint).To(Equal(args[0]))
			Expect(invokedExpectedResponseCode).To(Equal(0))
			Expect(invokedBody).To(Equal(strings.NewReader(data)))
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
					invokedBody := fakePivnetClient.MakeRequestArgsForCall(0)

				Expect(invokedBody).To(BeNil())
			})
		})

		Context("when there is an error", func() {
			BeforeEach(func() {
				makeRequestErr = errors.New("curl error")
			})

			It("invokes the error handler", func() {
				err := client.MakeRequest(method, data, args)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(makeRequestErr))
			})
		})

		Context("when there is an error unmarshalling json", func() {
			BeforeEach(func() {
				returnedBytes := []byte(`[garbage-json!`)
				resp = &http.Response{
					Body: ioutil.NopCloser(bytes.NewReader(returnedBytes)),
				}
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
