package commands_test

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/v3/commands"
	"github.com/pivotal-cf/pivnet-cli/v3/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/v3/commands/eula"
)

var _ = Describe("eula commands", func() {
	var (
		field reflect.StructField

		fakeEULAClient *commandsfakes.FakeEULAClient
	)

	BeforeEach(func() {
		fakeEULAClient = &commandsfakes.FakeEULAClient{}

		commands.NewEULAClient = func(eula.PivnetClient) commands.EULAClient {
			return fakeEULAClient
		}
	})

	Describe("EULAsCommand", func() {
		var (
			cmd *commands.EULAsCommand
		)

		BeforeEach(func() {
			cmd = &commands.EULAsCommand{}
		})

		It("invokes the EULA client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeEULAClient.ListCallCount()).To(Equal(1))
		})

		Context("when the EULA client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeEULAClient.ListReturns(expectedErr)
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

	Describe("EULACommand", func() {
		var (
			cmd *commands.EULACommand

			eulaSlug string
		)

		BeforeEach(func() {
			eulaSlug = "some eula slug"

			cmd = &commands.EULACommand{
				EULASlug: eulaSlug,
			}
		})

		It("invokes the EULA client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeEULAClient.GetCallCount()).To(Equal(1))
			Expect(fakeEULAClient.GetArgsForCall(0)).To(Equal(eulaSlug))
		})

		Context("when the EULA client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeEULAClient.GetReturns(expectedErr)
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

		Describe("EULASlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.EULACommand{}, "EULASlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("eula-slug"))
			})
		})
	})

	Describe("AcceptEULACommand", func() {
		var (
			cmd *commands.AcceptEULACommand

			productSlug    string
			releaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			releaseVersion = "some release version"

			cmd = &commands.AcceptEULACommand{
				ProductSlug:    productSlug,
				ReleaseVersion: releaseVersion,
			}
		})

		It("invokes the EULA client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeEULAClient.AcceptEULACallCount()).To(Equal(1))
			ps, rv := fakeEULAClient.AcceptEULAArgsForCall(0)
			Expect(ps).To(Equal(productSlug))
			Expect(rv).To(Equal(releaseVersion))
		})

		Context("when the EULA client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeEULAClient.AcceptEULAReturns(expectedErr)
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

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AcceptEULACommand{}, "ProductSlug")
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

		Describe("ReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AcceptEULACommand{}, "ReleaseVersion")
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
})
