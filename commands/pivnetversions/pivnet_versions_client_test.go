package pivnetversions_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-cli/commands/pivnetversions"
	"github.com/pivotal-cf/pivnet-cli/commands/pivnetversions/pivnetversionsfakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("pivnetversions commands", func() {
	var (
		fakePivnetClient *pivnetversionsfakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		pivnetVersions pivnet.PivnetVersions

		client *pivnetversions.PivnetVersionsClient
	)

	BeforeEach(func() {
		fakePivnetClient = &pivnetversionsfakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		pivnetVersions = pivnet.PivnetVersions{"1.2.3", "3.2.1"}

		fakePivnetClient.PivnetVersionsReturns(pivnetVersions, nil)

		client = pivnetversions.NewPivnetVersionsClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("PivnetVersions", func() {
		It("lists all PivnetVersions", func() {
			err := client.List()
			Expect(err).NotTo(HaveOccurred())

			var returnedPivnetVersions pivnet.PivnetVersions
			err = json.Unmarshal(outBuffer.Bytes(), &returnedPivnetVersions)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedPivnetVersions).To(Equal(pivnetVersions))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("pivnetversionss error")
				fakePivnetClient.PivnetVersionsReturns(pivnet.PivnetVersions{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.List()
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
