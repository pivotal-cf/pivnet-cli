package releasetype_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet/v7"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/releasetype"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/releasetype/releasetypefakes"
	"github.com/pivotal-cf/pivnet-cli/v2/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/v2/printer"
)

var _ = Describe("releasetype commands", func() {
	var (
		fakePivnetClient *releasetypefakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		releasetypes []pivnet.ReleaseType

		client *releasetype.ReleaseTypeClient
	)

	BeforeEach(func() {
		fakePivnetClient = &releasetypefakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		releasetypes = []pivnet.ReleaseType{
			pivnet.ReleaseType("release-type-A"),
			pivnet.ReleaseType("release-type-B"),
		}

		fakePivnetClient.ReleaseTypesReturns(releasetypes, nil)

		client = releasetype.NewReleaseTypeClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("ReleaseTypes", func() {
		It("lists all ReleaseTypes", func() {
			err := client.List()
			Expect(err).NotTo(HaveOccurred())

			var returnedReleaseTypes []pivnet.ReleaseType
			err = json.Unmarshal(outBuffer.Bytes(), &returnedReleaseTypes)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedReleaseTypes).To(Equal(releasetypes))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releasetypes error")
				fakePivnetClient.ReleaseTypesReturns(nil, expectedErr)
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
