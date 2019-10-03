package commands_test

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands/companygroup"

	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
)

var _ = Describe("user group commands", func() {
	var (
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
})