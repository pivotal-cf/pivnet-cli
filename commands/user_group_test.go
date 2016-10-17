package commands_test

import (
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
)

var _ = Describe("user group commands", func() {
	var (
		field reflect.StructField

		fakeUserGroupClient *commandsfakes.FakeUserGroupClient
	)

	BeforeEach(func() {
		fakeUserGroupClient = &commandsfakes.FakeUserGroupClient{}

		commands.NewUserGroupClient = func() commands.UserGroupClient {
			return fakeUserGroupClient
		}
	})

	Describe("UserGroupsCommand", func() {
		var (
			cmd commands.UserGroupsCommand
		)

		BeforeEach(func() {
			cmd = commands.UserGroupsCommand{}
		})

		It("invokes the UserGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeUserGroupClient.ListCallCount()).To(Equal(1))
		})

		Context("when the UserGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeUserGroupClient.ListReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.UserGroupsCommand{}, "ProductSlug")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("p"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-slug"))
			})
		})

		Describe("ReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.UserGroupsCommand{}, "ReleaseVersion")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("r"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-version"))
			})
		})
	})

	Describe("UserGroupCommand", func() {
		var (
			cmd commands.UserGroupCommand
		)

		BeforeEach(func() {
			cmd = commands.UserGroupCommand{}
		})

		It("invokes the UserGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeUserGroupClient.GetCallCount()).To(Equal(1))
		})

		Context("when the UserGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeUserGroupClient.GetReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})
	})

	Describe("CreateUserGroupCommand", func() {
		var (
			cmd commands.CreateUserGroupCommand
		)

		BeforeEach(func() {
			cmd = commands.CreateUserGroupCommand{}
		})

		It("invokes the UserGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeUserGroupClient.CreateCallCount()).To(Equal(1))
		})

		Context("when the UserGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeUserGroupClient.CreateReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Describe("Name flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.CreateUserGroupCommand{}, "Name")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("name"))
			})
		})

		Describe("Description flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.CreateUserGroupCommand{}, "Description")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("description"))
			})
		})

		Describe("Members flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.CreateUserGroupCommand{}, "Members")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("member"))
			})
		})
	})

	Describe("UpdateUserGroupCommand", func() {
		var (
			cmd commands.UpdateUserGroupCommand
		)

		BeforeEach(func() {
			cmd = commands.UpdateUserGroupCommand{}
		})

		It("invokes the UserGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeUserGroupClient.UpdateCallCount()).To(Equal(1))
		})

		Context("when the UserGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeUserGroupClient.UpdateReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Describe("Name flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.UpdateUserGroupCommand{}, "Name")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("name"))
			})
		})

		Describe("Description flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.UpdateUserGroupCommand{}, "Description")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("description"))
			})
		})
	})

	Describe("AddUserGroupCommand", func() {
		var (
			cmd commands.AddUserGroupCommand
		)

		BeforeEach(func() {
			cmd = commands.AddUserGroupCommand{}
		})

		It("invokes the UserGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeUserGroupClient.AddToReleaseCallCount()).To(Equal(1))
		})

		Context("when the UserGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeUserGroupClient.AddToReleaseReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddUserGroupCommand{}, "ProductSlug")
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

		Describe("UserGroupID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddUserGroupCommand{}, "UserGroupID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("user-group-id"))
			})
		})

		Describe("ReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddUserGroupCommand{}, "ReleaseVersion")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("r"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-version"))
			})
		})
	})

	Describe("RemoveUserGroupCommand", func() {
		var (
			cmd commands.RemoveUserGroupCommand
		)

		BeforeEach(func() {
			cmd = commands.RemoveUserGroupCommand{}
		})

		It("invokes the UserGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeUserGroupClient.RemoveFromReleaseCallCount()).To(Equal(1))
		})

		Context("when the UserGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeUserGroupClient.RemoveFromReleaseReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveUserGroupCommand{}, "ProductSlug")
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

		Describe("UserGroupID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveUserGroupCommand{}, "UserGroupID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("user-group-id"))
			})
		})

		Describe("ReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveUserGroupCommand{}, "ReleaseVersion")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("r"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-version"))
			})
		})
	})

	Describe("DeleteUserGroupCommand", func() {
		var (
			cmd commands.DeleteUserGroupCommand
		)

		BeforeEach(func() {
			cmd = commands.DeleteUserGroupCommand{}
		})

		It("invokes the UserGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeUserGroupClient.DeleteCallCount()).To(Equal(1))
		})

		Context("when the UserGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeUserGroupClient.DeleteReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Describe("UserGroupID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DeleteUserGroupCommand{}, "UserGroupID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("user-group-id"))
			})
		})
	})

	Describe("AddUserGroupMemberCommand", func() {
		var (
			cmd commands.AddUserGroupMemberCommand
		)

		BeforeEach(func() {
			cmd = commands.AddUserGroupMemberCommand{}
		})

		It("invokes the UserGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeUserGroupClient.AddUserGroupMemberCallCount()).To(Equal(1))
		})

		Context("when the UserGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeUserGroupClient.AddUserGroupMemberReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Describe("UserGroupID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddUserGroupMemberCommand{}, "UserGroupID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("user-group-id"))
			})
		})

		Describe("MemberEmailAddress flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddUserGroupMemberCommand{}, "MemberEmailAddress")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("member-email"))
			})
		})

		Describe("Admin flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddUserGroupMemberCommand{}, "Admin")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("admin"))
			})
		})
	})

	Describe("RemoveUserGroupMemberCommand", func() {
		var (
			cmd commands.RemoveUserGroupMemberCommand
		)

		BeforeEach(func() {
			cmd = commands.RemoveUserGroupMemberCommand{}
		})

		It("invokes the UserGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeUserGroupClient.RemoveUserGroupMemberCallCount()).To(Equal(1))
		})

		Context("when the UserGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeUserGroupClient.RemoveUserGroupMemberReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Describe("UserGroupID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveUserGroupMemberCommand{}, "UserGroupID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("user-group-id"))
			})
		})

		Describe("MemberEmailAddress flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveUserGroupMemberCommand{}, "MemberEmailAddress")
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
