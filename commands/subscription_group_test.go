package commands_test

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands/subscriptiongroup"
	"reflect"

	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
)

var _ = Describe("subscription group commands", func() {
	var (
		field reflect.StructField
		fakeSubscriptionGroupClient *commandsfakes.FakeSubscriptionGroupClient
	)

	BeforeEach(func() {
		fakeSubscriptionGroupClient = &commandsfakes.FakeSubscriptionGroupClient{}

		commands.NewSubscriptionGroupClient = func(subscriptiongroup.PivnetClient) commands.SubscriptionGroupClient {
			return fakeSubscriptionGroupClient
		}
	})

	Describe("SubscriptionGroupsCommand", func() {
		var (
			cmd commands.SubscriptionGroupsCommand
		)

		BeforeEach(func() {
			cmd = commands.SubscriptionGroupsCommand{}
		})

		It("invokes the SubscriptionGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeSubscriptionGroupClient.ListCallCount()).To(Equal(1))
		})

		Context("when the SubscriptionGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeSubscriptionGroupClient.ListReturns(expectedErr)
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

	Describe("SubscriptionGroupCommand", func() {
		var (
			cmd commands.SubscriptionGroupCommand
		)

		BeforeEach(func() {
			cmd = commands.SubscriptionGroupCommand{}
		})

		It("invokes the SubscriptionGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeSubscriptionGroupClient.GetCallCount()).To(Equal(1))
		})

		Context("when the SubscriptionGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeSubscriptionGroupClient.GetReturns(expectedErr)
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

		Describe("SubscriptionGroupId flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.SubscriptionGroupCommand{}, "SubscriptionGroupID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("subscription-group-id"))
			})
		})
	})

	Describe("AddSubscriptionGroupMemberCommand", func() {
		var (
			cmd commands.AddSubscriptionGroupMemberCommand
		)

		BeforeEach(func() {
			cmd = commands.AddSubscriptionGroupMemberCommand{}
		})

		It("invokes the SubscriptionGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeSubscriptionGroupClient.AddMemberCallCount()).To(Equal(1))
		})

		Context("when the SubscriptionGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeSubscriptionGroupClient.AddMemberReturns(expectedErr)
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

		Describe("SubscriptionGroupId flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddSubscriptionGroupMemberCommand{}, "SubscriptionGroupID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("subscription-group-id"))
			})
		})

		Describe("MemberEmail flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddSubscriptionGroupMemberCommand{}, "MemberEmail")
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
				field = fieldFor(commands.AddSubscriptionGroupMemberCommand{}, "IsAdmin")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("admin"))
			})
		})
	})

	Describe("RemoveSubscriptionGroupMemberCommand", func() {
		var (
			cmd commands.RemoveSubscriptionGroupMemberCommand
		)

		BeforeEach(func() {
			cmd = commands.RemoveSubscriptionGroupMemberCommand{}
		})

		It("invokes the SubscriptionGroup client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeSubscriptionGroupClient.RemoveMemberCallCount()).To(Equal(1))
		})

		Context("when the SubscriptionGroup client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeSubscriptionGroupClient.RemoveMemberReturns(expectedErr)
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

		Describe("SubscriptionGroupId flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveSubscriptionGroupMemberCommand{}, "SubscriptionGroupID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("subscription-group-id"))
			})
		})

		Describe("MemberEmail flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveSubscriptionGroupMemberCommand{}, "MemberEmail")
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