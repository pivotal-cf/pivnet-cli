package commands_test

import (
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
)

var _ = Describe("release dependency commands", func() {
	var (
		field reflect.StructField

		fakeReleaseDependencyClient *commandsfakes.FakeReleaseDependencyClient
	)

	BeforeEach(func() {
		fakeReleaseDependencyClient = &commandsfakes.FakeReleaseDependencyClient{}

		commands.NewReleaseDependencyClient = func() commands.ReleaseDependencyClient {
			return fakeReleaseDependencyClient
		}
	})

	Describe("ReleasesDependenciesCommand", func() {
		var (
			cmd *commands.ReleaseDependenciesCommand
		)

		BeforeEach(func() {
			cmd = &commands.ReleaseDependenciesCommand{}
		})

		It("invokes the ReleaseDependency client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseDependencyClient.ListCallCount()).To(Equal(1))
		})

		Context("when the ReleaseDependency client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeReleaseDependencyClient.ListReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ReleaseDependenciesCommand{}, "ProductSlug")
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
				field = fieldFor(commands.ReleaseDependenciesCommand{}, "ReleaseVersion")
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

	Describe("AddReleasesDependenciesCommand", func() {
		var (
			cmd commands.AddReleaseDependencyCommand
		)

		BeforeEach(func() {
			cmd = commands.AddReleaseDependencyCommand{}
		})

		It("invokes the ReleaseDependency client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseDependencyClient.AddCallCount()).To(Equal(1))
		})

		Context("when the ReleaseDependency client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeReleaseDependencyClient.AddReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddReleaseDependencyCommand{}, "ProductSlug")
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
				field = fieldFor(commands.AddReleaseDependencyCommand{}, "ReleaseVersion")
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

		Describe("DependentProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddReleaseDependencyCommand{}, "DependentProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("s"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("dependent-product-slug"))
			})
		})

		Describe("DependentReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddReleaseDependencyCommand{}, "DependentReleaseVersion")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("u"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("dependent-release-version"))
			})
		})
	})

	Describe("RemoveReleasesDependenciesCommand", func() {
		var (
			cmd commands.RemoveReleaseDependencyCommand
		)

		BeforeEach(func() {
			cmd = commands.RemoveReleaseDependencyCommand{}
		})

		It("invokes the ReleaseDependency client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeReleaseDependencyClient.RemoveCallCount()).To(Equal(1))
		})

		Context("when the ReleaseDependency client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeReleaseDependencyClient.RemoveReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
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

		Describe("DependentProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "DependentProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("s"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("dependent-product-slug"))
			})
		})

		Describe("DependentReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "DependentReleaseVersion")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("u"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("dependent-release-version"))
			})
		})
	})
})
