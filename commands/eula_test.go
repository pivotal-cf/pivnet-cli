package commands_test

import (
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
)

var _ = Describe("eula commands", func() {
	var (
		field reflect.StructField

		fakeEULAClient *commandsfakes.FakeEULAClient
	)

	BeforeEach(func() {
		fakeEULAClient = &commandsfakes.FakeEULAClient{}

		commands.NewEULAClient = func() commands.EULAClient {
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
