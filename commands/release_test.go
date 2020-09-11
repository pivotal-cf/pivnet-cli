package commands_test

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/v2/commands"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/release"
)

var _ = Describe("release commands", func() {
	var (
		field reflect.StructField

		fakeReleaseClient *commandsfakes.FakeReleaseClient
	)

	BeforeEach(func() {
		fakeReleaseClient = &commandsfakes.FakeReleaseClient{}

		commands.NewReleaseClient = func(release.PivnetClient) commands.ReleaseClient {
			return fakeReleaseClient
		}
	})

	Describe("ReleasesCommand", func() {
		var (
			cmd commands.ReleasesCommand
		)

		BeforeEach(func() {
			cmd = commands.ReleasesCommand{}
		})

		It("invokes the Release client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseClient.ListCallCount()).To(Equal(1))
		})

		Context("when a limit is provided", func() {
			var (
				slug string
				limit string
			)

			BeforeEach(func() {
				slug = "product-slug"
				limit = "5"
				cmd = commands.ReleasesCommand{slug, limit}
			})

			It("invokes the Release client for the method ListWithLimit", func() {
				Expect(fakeReleaseClient.ListWithLimitCallCount()).To(Equal(0))
				err := cmd.Execute(nil)

				Expect(err).ToNot(HaveOccurred())
				Expect(fakeReleaseClient.ListWithLimitCallCount()).To(Equal(1))
				slug, params := fakeReleaseClient.ListWithLimitArgsForCall(0)

				Expect(slug).To(Equal(slug))
				Expect(params).To(Equal(limit))
			})

			Context("when method ListWithLimit return an error", func() {
				var expectedErr error

				BeforeEach(func() {
					expectedErr = errors.New("expected error")
					fakeReleaseClient.ListWithLimitReturns(expectedErr)
				})

				It("forward the error", func() {
					err := cmd.Execute(nil)

					Expect(err).To(Equal(expectedErr))
				})
			})
		})

		Context("when the Release client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeReleaseClient.ListReturns(expectedErr)
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
				field = fieldFor(commands.ReleasesCommand{}, "ProductSlug")
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

		Describe("limit flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ReleasesCommand{}, "Limit")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("l"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("limit"))
			})
		})
	})

	Describe("ReleaseCommand", func() {
		var (
			cmd commands.ReleaseCommand
		)

		BeforeEach(func() {
			cmd = commands.ReleaseCommand{}
		})

		It("invokes the Release client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseClient.GetCallCount()).To(Equal(1))
		})

		Context("when the Release client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeReleaseClient.GetReturns(expectedErr)
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
				field = fieldFor(commands.ReleaseCommand{}, "ProductSlug")
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
				field = fieldFor(commands.ReleaseCommand{}, "ReleaseVersion")
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

	Describe("CreateReleaseCommand", func() {
		var (
			cmd commands.CreateReleaseCommand
		)

		BeforeEach(func() {
			cmd = cmd
		})

		It("invokes the Release client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseClient.CreateCallCount()).To(Equal(1))
		})

		Context("when the Release client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeReleaseClient.CreateReturns(expectedErr)
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
				field = fieldFor(cmd, "ProductSlug")
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
				field = fieldFor(cmd, "ReleaseVersion")
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

		Describe("ReleaseType flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ReleaseType")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("t"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-type"))
			})
		})

		Describe("EULASlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "EULASlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("e"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("eula-slug"))
			})
		})
	})

	Describe("UpdateReleaseCommand", func() {
		var (
			cmd commands.UpdateReleaseCommand
		)

		BeforeEach(func() {
			cmd = commands.UpdateReleaseCommand{}
		})

		It("invokes the Release client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseClient.UpdateCallCount()).To(Equal(1))
		})

		Context("when the Release client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeReleaseClient.UpdateReturns(expectedErr)
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
				field = fieldFor(cmd, "ProductSlug")
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
				field = fieldFor(cmd, "ReleaseVersion")
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

		Describe("Availability flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Availability")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("availability"))
			})
		})

		Describe("ReleaseType flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ReleaseType")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-type"))
			})
		})
	})

	Describe("DeleteReleaseCommand", func() {
		var (
			cmd commands.DeleteReleaseCommand
		)

		BeforeEach(func() {
			cmd = commands.DeleteReleaseCommand{}
		})

		It("invokes the Release client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseClient.DeleteCallCount()).To(Equal(1))
		})

		Context("when the Release client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeReleaseClient.DeleteReturns(expectedErr)
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
				field = fieldFor(commands.DeleteReleaseCommand{}, "ProductSlug")
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
				field = fieldFor(commands.DeleteReleaseCommand{}, "ReleaseVersion")
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
