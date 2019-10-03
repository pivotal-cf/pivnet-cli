package companygroup_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/pivotal-cf/pivnet-cli/commands/companygroup"
	"github.com/pivotal-cf/pivnet-cli/commands/companygroup/companygroupfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("company group client", func() {
	var (
		fakePivnetClient *companygroupfakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		companygroups []pivnet.CompanyGroup

		client *companygroup.CompanyGroupClient
	)

	BeforeEach(func() {
		fakePivnetClient = &companygroupfakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		companygroups = []pivnet.CompanyGroup{
			{
				ID:          1234,
				Name:        "company-group-1234",
			},
			{
				ID:          2345,
				Name:        "company-group-2345",
			},
		}

		client = companygroup.NewCompanyGroupClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("List", func() {
		BeforeEach(func() {
			fakePivnetClient.CompanyGroupsReturns(companygroups, nil)
		})

		It("lists all CompanyGroups", func() {
			err := client.List()
			Expect(err).NotTo(HaveOccurred())

			var returnedCompanyGroups []pivnet.CompanyGroup
			err = json.Unmarshal(outBuffer.Bytes(), &returnedCompanyGroups)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedCompanyGroups).To(Equal(companygroups))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("companygroups error")
				fakePivnetClient.CompanyGroupsReturns(nil, expectedErr)
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
