package pivnetversions_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet/v3"
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

	Describe("PivnetVersions.List", func() {
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
				expectedErr = errors.New("pivnetversions error")
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

	Describe("PivnetVersions.Warn", func() {
		It("returns a warning if the current Pivnet CLI version is out of date", func() {
			result := client.Warn("1.2.2")
			Expect(result).NotTo(BeEmpty())
		})

		It("returns empty if the current Pivnet CLI version is up to date", func() {
			result := client.Warn("1.2.3")
			Expect(result).To(BeEmpty())
		})

		It("returns empty if the current Pivnet CLI version is newer", func() {
			result := client.Warn("1.2.4")
			Expect(result).To(BeEmpty())
		})

		It("returns empty if there is an error getting the latest Pivnet CLI version", func() {
			expectedErr := errors.New("pivnetversions error")
			fakePivnetClient.PivnetVersionsReturns(pivnet.PivnetVersions{}, expectedErr)

			result := client.Warn("1.2.3")
			Expect(result).To(BeEmpty())
		})
	})
})
