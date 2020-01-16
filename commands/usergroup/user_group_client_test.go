package usergroup_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet/v4"
	"github.com/pivotal-cf/pivnet-cli/commands/usergroup"
	"github.com/pivotal-cf/pivnet-cli/commands/usergroup/usergroupfakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("usergroup commands", func() {
	var (
		fakePivnetClient *usergroupfakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		usergroups []pivnet.UserGroup

		client *usergroup.UserGroupClient
	)

	BeforeEach(func() {
		fakePivnetClient = &usergroupfakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		usergroups = []pivnet.UserGroup{
			{
				ID:          1234,
				Name:        "user-group-1234",
				Description: "first",
			},
			{
				ID:          2345,
				Name:        "user-group-2345",
				Description: "second",
			},
		}

		client = usergroup.NewUserGroupClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("List", func() {
		var (
			productSlug    string
			releaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = ""

			fakePivnetClient.UserGroupsReturns(usergroups, nil)
		})

		It("lists all UserGroups", func() {
			err := client.List(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())

			var returnedUserGroups []pivnet.UserGroup
			err = json.Unmarshal(outBuffer.Bytes(), &returnedUserGroups)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedUserGroups).To(Equal(usergroups))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("usergroups error")
				fakePivnetClient.UserGroupsReturns(nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.List(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when release version is not empty", func() {
			BeforeEach(func() {
				releaseVersion = "some-release-version"
				fakePivnetClient.UserGroupsForReleaseReturns(usergroups, nil)
			})

			It("lists all UserGroups", func() {
				err := client.List(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				var returnedUserGroups []pivnet.UserGroup
				err = json.Unmarshal(outBuffer.Bytes(), &returnedUserGroups)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedUserGroups).To(Equal(usergroups))
			})

			Context("when there is an error getting release", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("releases error")
					fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.List(productSlug, releaseVersion)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})

			Context("when there is an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("usergroups error")
					fakePivnetClient.UserGroupsForReleaseReturns(nil, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.List(productSlug, releaseVersion)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})
		})
	})

	Describe("Create", func() {
		var (
			name        string
			description string
			members     []string

			createdUserGroup pivnet.UserGroup
		)

		BeforeEach(func() {
			name = "new-name"
			description = "new-description"
			members = []string{"member-1", "member-2"}

			createdUserGroup = pivnet.UserGroup{
				Name:        name,
				Description: description,
				Members:     members,
			}

			fakePivnetClient.CreateUserGroupReturns(createdUserGroup, nil)
		})

		It("gets UserGroup", func() {
			err := client.Create(name, description, members)
			Expect(err).NotTo(HaveOccurred())

			var returnedUserGroup pivnet.UserGroup
			err = json.Unmarshal(outBuffer.Bytes(), &returnedUserGroup)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedUserGroup).To(Equal(createdUserGroup))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("usergroup error")
				fakePivnetClient.CreateUserGroupReturns(pivnet.UserGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Create(name, description, members)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Get", func() {
		var (
			userGroupID int
		)

		BeforeEach(func() {
			userGroupID = usergroups[0].ID

			fakePivnetClient.UserGroupReturns(usergroups[0], nil)
		})

		It("gets UserGroup", func() {
			err := client.Get(userGroupID)
			Expect(err).NotTo(HaveOccurred())

			var returnedUserGroup pivnet.UserGroup
			err = json.Unmarshal(outBuffer.Bytes(), &returnedUserGroup)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedUserGroup).To(Equal(usergroups[0]))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("usergroup error")
				fakePivnetClient.UserGroupReturns(pivnet.UserGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(userGroupID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("AddToRelease", func() {
		var (
			productSlug    string
			releaseVersion string
			userGroupID    int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			userGroupID = usergroups[0].ID

			fakePivnetClient.AddUserGroupReturns(nil)
		})

		It("adds UserGroup", func() {
			err := client.AddToRelease(productSlug, releaseVersion, userGroupID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error getting release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddToRelease(productSlug, releaseVersion, userGroupID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("usergroup error")
				fakePivnetClient.AddUserGroupReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddToRelease(productSlug, releaseVersion, userGroupID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("RemoveFromRelease", func() {
		var (
			productSlug    string
			releaseVersion string
			userGroupID    int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			userGroupID = usergroups[0].ID

			fakePivnetClient.RemoveUserGroupReturns(nil)
		})

		It("deletes UserGroup", func() {
			err := client.RemoveFromRelease(productSlug, releaseVersion, userGroupID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error getting release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveFromRelease(productSlug, releaseVersion, userGroupID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("usergroup error")
				fakePivnetClient.RemoveUserGroupReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveFromRelease(productSlug, releaseVersion, userGroupID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Delete", func() {
		var (
			userGroupID int
		)

		BeforeEach(func() {
			userGroupID = usergroups[0].ID

			fakePivnetClient.DeleteUserGroupReturns(nil)
		})

		It("deletes UserGroup", func() {
			err := client.Delete(userGroupID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("usergroup error")
				fakePivnetClient.DeleteUserGroupReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Delete(userGroupID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Update", func() {
		var (
			userGroupID int
			name        *string
			description *string
		)

		BeforeEach(func() {
			userGroupID = usergroups[0].ID
			name = nil
			description = nil

			fakePivnetClient.UserGroupReturns(usergroups[0], nil)
			fakePivnetClient.UpdateUserGroupReturns(usergroups[0], nil)
		})

		It("updates UserGroup with existing values", func() {
			err := client.Update(userGroupID, name, description)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePivnetClient.UpdateUserGroupCallCount()).To(Equal(1))
			invokedUserGroup := fakePivnetClient.UpdateUserGroupArgsForCall(0)

			Expect(invokedUserGroup.ID).To(Equal(userGroupID))
			Expect(invokedUserGroup.Name).To(Equal(usergroups[0].Name))
			Expect(invokedUserGroup.Description).To(Equal(usergroups[0].Description))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("usergroup error")
				fakePivnetClient.UpdateUserGroupReturns(pivnet.UserGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Update(userGroupID, name, description)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when the name is non-nil", func() {
			var (
				nameVal string
			)

			BeforeEach(func() {
				nameVal = "some-name"
				name = &nameVal
			})

			It("updates UserGroup with the provided name", func() {
				err := client.Update(userGroupID, name, description)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateUserGroupCallCount()).To(Equal(1))
				invokedUserGroup := fakePivnetClient.UpdateUserGroupArgsForCall(0)

				Expect(invokedUserGroup.ID).To(Equal(userGroupID))
				Expect(invokedUserGroup.Name).To(Equal(nameVal))
				Expect(invokedUserGroup.Description).To(Equal(usergroups[0].Description))
			})
		})

		Context("when the description is non-nil", func() {
			var (
				descriptionVal string
			)

			BeforeEach(func() {
				descriptionVal = "some-description"
				description = &descriptionVal
			})

			It("updates UserGroup with the provided name", func() {
				err := client.Update(userGroupID, name, description)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateUserGroupCallCount()).To(Equal(1))
				invokedUserGroup := fakePivnetClient.UpdateUserGroupArgsForCall(0)

				Expect(invokedUserGroup.ID).To(Equal(userGroupID))
				Expect(invokedUserGroup.Name).To(Equal(usergroups[0].Name))
				Expect(invokedUserGroup.Description).To(Equal(descriptionVal))
			})
		})
	})

	Describe("AddUserGroupMember", func() {
		var (
			userGroupID int
			memberEmail string
			admin       bool
		)

		BeforeEach(func() {
			userGroupID = usergroups[0].ID
			memberEmail = "some-product-slug"
			admin = true

			fakePivnetClient.AddMemberToGroupReturns(usergroups[0], nil)
		})

		It("adds user to UserGroup", func() {
			err := client.AddUserGroupMember(userGroupID, memberEmail, admin)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("usergroup error")
				fakePivnetClient.AddMemberToGroupReturns(pivnet.UserGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddUserGroupMember(userGroupID, memberEmail, admin)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("RemoveUserGroupMember", func() {
		var (
			userGroupID int
			memberEmail string
		)

		BeforeEach(func() {
			userGroupID = usergroups[0].ID
			memberEmail = "some-product-slug"

			fakePivnetClient.RemoveMemberFromGroupReturns(usergroups[0], nil)
		})

		It("removes user from UserGroup", func() {
			err := client.RemoveUserGroupMember(userGroupID, memberEmail)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("usergroup error")
				fakePivnetClient.RemoveMemberFromGroupReturns(pivnet.UserGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveUserGroupMember(userGroupID, memberEmail)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
