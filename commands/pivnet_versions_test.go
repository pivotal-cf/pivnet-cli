package commands_test

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/commands/pivnetversions"
)

var _ = Describe("pivnet versions command", func() {
	var (
		fakePivnetVersionsClient *commandsfakes.FakePivnetVersionsClient
	)

	BeforeEach(func() {
		fakePivnetVersionsClient = &commandsfakes.FakePivnetVersionsClient{}

		commands.NewPivnetVersionsClient = func(pivnetversions.PivnetClient) commands.PivnetVersionsClient {
			return fakePivnetVersionsClient
		}
	})

	Describe("PivnetVersionsCommand", func() {
		var (
			cmd *commands.PivnetVersionsCommand
		)

		BeforeEach(func() {
			cmd = &commands.PivnetVersionsCommand{}
		})

		It("invokes the PivnetVersions client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakePivnetVersionsClient.ListCallCount()).To(Equal(1))
		})

		Context("when the PivnetVersions client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakePivnetVersionsClient.ListReturns(expectedErr)
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
})
