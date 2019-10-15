package commands_test

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands/companygroup"
	"reflect"

	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
)

var _ = Describe("company group commands", func() {
	var (
		field reflect.StructField
		fakeCompanyGroupClient *commandsfakes.FakeCompanyGroupClient
	)

	BeforeEach(func() {
		fakeCompanyGroupClient = &commandsfakes.FakeCompanyGroupClient{}

		commands.NewCompanyGroupClient = func(companygroup.PivnetClient) commands.CompanyGroupClient {
			return fakeCompanyGroupClient
		}
	})

	Describe("CompanyGroupsCommand", func() {
		var (
			cmd commands.CompanyGroupsCommand
		)

		BeforeEach(func() {
			cmd = commands.CompanyGroupsCommand{}
		})

		It("invokes the CompanyGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeCompanyGroupClient.ListCallCount()).To(Equal(1))
		})

		Context("when the CompanyGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeCompanyGroupClient.ListReturns(expectedErr)
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
	})

	Describe("CompanyGroupCommand", func() {
		var (
			cmd commands.CompanyGroupCommand
		)

		BeforeEach(func() {
			cmd = commands.CompanyGroupCommand{}
		})

		It("invokes the CompanyGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeCompanyGroupClient.GetCallCount()).To(Equal(1))
		})

		Context("when the CompanyGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeCompanyGroupClient.GetReturns(expectedErr)
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

		Describe("CompanyGroupId flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.CompanyGroupCommand{}, "CompanyGroupID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("company-group-id"))
			})
		})
	})

	Describe("AddCompanyGroupMemberCommand", func() {
		var (
			cmd commands.AddCompanyGroupMemberCommand
		)

		BeforeEach(func() {
			cmd = commands.AddCompanyGroupMemberCommand{}
		})

		It("invokes the CompanyGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeCompanyGroupClient.AddMemberCallCount()).To(Equal(1))
		})

		Context("when the CompanyGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeCompanyGroupClient.AddMemberReturns(expectedErr)
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

		Describe("CompanyGroupId flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddCompanyGroupMemberCommand{}, "CompanyGroupID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("company-group-id"))
			})
		})

		Describe("MemberEmail flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddCompanyGroupMemberCommand{}, "MemberEmail")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("member-email"))
			})
		})

		Describe("IsAdmin flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddCompanyGroupMemberCommand{}, "IsAdmin")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("admin"))
			})
		})
	})

	Describe("RemoveCompanyGroupMemberCommand", func() {
		var (
			cmd commands.RemoveCompanyGroupMemberCommand
		)

		BeforeEach(func() {
			cmd = commands.RemoveCompanyGroupMemberCommand{}
		})

		It("invokes the CompanyGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeCompanyGroupClient.RemoveMemberCallCount()).To(Equal(1))
		})

		Context("when the CompanyGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeCompanyGroupClient.RemoveMemberReturns(expectedErr)
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

		Describe("CompanyGroupId flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveCompanyGroupMemberCommand{}, "CompanyGroupID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("company-group-id"))
			})
		})

		Describe("MemberEmail flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveCompanyGroupMemberCommand{}, "MemberEmail")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("member-email"))
			})
		})
	})
})