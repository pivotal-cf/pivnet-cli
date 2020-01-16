package companygroup_test

import (
	"bytes"
	"errors"
	"github.com/pivotal-cf/pivnet-cli/commands/companygroup"
	"github.com/pivotal-cf/pivnet-cli/commands/companygroup/companygroupfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet/v3"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("company group client", func() {
	var (
		fakePivnetClient *companygroupfakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		client *companygroup.CompanyGroupClient
	)

	BeforeEach(func() {
		fakePivnetClient = &companygroupfakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		client = companygroup.NewCompanyGroupClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("List", func() {
		It("List the CompanyGroups", func() {
			Expect(fakePivnetClient.CompanyGroupsCallCount()).To(Equal(0))

			_ = client.List()

			Expect(fakePivnetClient.CompanyGroupsCallCount()).To(Equal(1))
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

	Describe("Get", func() {
		It("Get the CompanyGroup", func() {
			Expect(fakePivnetClient.CompanyGroupCallCount()).To(Equal(0))

			_ = client.Get(4567)

			Expect(fakePivnetClient.CompanyGroupCallCount()).To(Equal(1))

			companyGroupId := fakePivnetClient.CompanyGroupArgsForCall(0)
			Expect(companyGroupId).To(Equal(4567))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("companygroup error")
				fakePivnetClient.CompanyGroupReturns(pivnet.CompanyGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(4567)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("AddCompanyGroupMember", func() {
		It("Add member to the company group and return the CompanyGroup", func() {
			Expect(fakePivnetClient.AddCompanyGroupMemberCallCount()).To(Equal(0))

			_ = client.AddMember(4567, "dude@dude.dude", "false")

			Expect(fakePivnetClient.AddCompanyGroupMemberCallCount()).To(Equal(1))

			companyGroupId, memberEmail, isAdmin := fakePivnetClient.AddCompanyGroupMemberArgsForCall(0)
			Expect(companyGroupId).To(Equal(4567))
			Expect(memberEmail).To(Equal("dude@dude.dude"))
			Expect(isAdmin).To(Equal("false"))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("companygroup error")
				fakePivnetClient.AddCompanyGroupMemberReturns(pivnet.CompanyGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddMember(123,"blah", "")
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("RemoveCompanyGroupMember", func() {
		It("Remove member to the company group and return the CompanyGroup", func() {
			Expect(fakePivnetClient.RemoveCompanyGroupMemberCallCount()).To(Equal(0))

			_ = client.RemoveMember(4567, "dude@dude.dude")

			Expect(fakePivnetClient.RemoveCompanyGroupMemberCallCount()).To(Equal(1))

			companyGroupId, memberEmail := fakePivnetClient.RemoveCompanyGroupMemberArgsForCall(0)
			Expect(companyGroupId).To(Equal(4567))
			Expect(memberEmail).To(Equal("dude@dude.dude"))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("companygroup error")
				fakePivnetClient.RemoveCompanyGroupMemberReturns(pivnet.CompanyGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveMember(123,"blah")
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
