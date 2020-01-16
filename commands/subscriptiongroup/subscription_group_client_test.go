package subscriptiongroup_test

import (
	"bytes"
	"errors"
	"github.com/pivotal-cf/pivnet-cli/commands/subscriptiongroup"
	"github.com/pivotal-cf/pivnet-cli/commands/subscriptiongroup/subscriptiongroupfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet/v4"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("subscription group client", func() {
	var (
		fakePivnetClient *subscriptiongroupfakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		client *subscriptiongroup.SubscriptionGroupClient
	)

	BeforeEach(func() {
		fakePivnetClient = &subscriptiongroupfakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		client = subscriptiongroup.NewSubscriptionGroupClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("List", func() {
		It("List the SubscriptionGroups", func() {
			Expect(fakePivnetClient.SubscriptionGroupsCallCount()).To(Equal(0))

			_ = client.List()

			Expect(fakePivnetClient.SubscriptionGroupsCallCount()).To(Equal(1))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("subscriptiongroups error")
				fakePivnetClient.SubscriptionGroupsReturns(nil, expectedErr)
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
		It("Get the SubscriptionGroup", func() {
			Expect(fakePivnetClient.SubscriptionGroupCallCount()).To(Equal(0))

			_ = client.Get(4567)

			Expect(fakePivnetClient.SubscriptionGroupCallCount()).To(Equal(1))

			subscriptionGroupId := fakePivnetClient.SubscriptionGroupArgsForCall(0)
			Expect(subscriptionGroupId).To(Equal(4567))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("subscriptiongroup error")
				fakePivnetClient.SubscriptionGroupReturns(pivnet.SubscriptionGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(4567)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("AddSubscriptionGroupMember", func() {
		It("Add member to the subscription group and return the SubscriptionGroup", func() {
			Expect(fakePivnetClient.AddSubscriptionGroupMemberCallCount()).To(Equal(0))

			_ = client.AddMember(4567, "dude@dude.dude", "false")

			Expect(fakePivnetClient.AddSubscriptionGroupMemberCallCount()).To(Equal(1))

			subscriptionGroupId, memberEmail, isAdmin := fakePivnetClient.AddSubscriptionGroupMemberArgsForCall(0)
			Expect(subscriptionGroupId).To(Equal(4567))
			Expect(memberEmail).To(Equal("dude@dude.dude"))
			Expect(isAdmin).To(Equal("false"))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("subscriptiongroup error")
				fakePivnetClient.AddSubscriptionGroupMemberReturns(pivnet.SubscriptionGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddMember(123,"blah", "")
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("RemoveSubscriptionGroupMember", func() {
		It("Remove member to the subscription group and return the SubscriptionGroup", func() {
			Expect(fakePivnetClient.RemoveSubscriptionGroupMemberCallCount()).To(Equal(0))

			_ = client.RemoveMember(4567, "dude@dude.dude")

			Expect(fakePivnetClient.RemoveSubscriptionGroupMemberCallCount()).To(Equal(1))

			subscriptionGroupId, memberEmail := fakePivnetClient.RemoveSubscriptionGroupMemberArgsForCall(0)
			Expect(subscriptionGroupId).To(Equal(4567))
			Expect(memberEmail).To(Equal("dude@dude.dude"))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("subscriptiongroup error")
				fakePivnetClient.RemoveSubscriptionGroupMemberReturns(pivnet.SubscriptionGroup{}, expectedErr)
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
