package commands_test

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/v2/commands"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/dependencyspecifier"
)

var _ = Describe("dependency specifier commands", func() {
	var (
		field reflect.StructField

		fakeDependencySpecifierClient *commandsfakes.FakeDependencySpecifierClient
	)

	BeforeEach(func() {
		fakeDependencySpecifierClient = &commandsfakes.FakeDependencySpecifierClient{}

		commands.NewDependencySpecifierClient = func(dependencyspecifier.PivnetClient) commands.DependencySpecifierClient {
			return fakeDependencySpecifierClient
		}
	})

	Describe("DependencySpecifiersCommand", func() {
		var (
			cmd *commands.DependencySpecifiersCommand
		)

		BeforeEach(func() {
			cmd = &commands.DependencySpecifiersCommand{}
		})

		It("invokes the DependencySpecifier client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeDependencySpecifierClient.ListCallCount()).To(Equal(1))
		})

		Context("when the DependencySpecifier client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeDependencySpecifierClient.ListReturns(expectedErr)
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
				field = fieldFor(commands.DependencySpecifiersCommand{}, "ProductSlug")
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
				field = fieldFor(commands.DependencySpecifiersCommand{}, "ReleaseVersion")
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

	Describe("DependencySpecifierCommand", func() {
		var (
			cmd commands.DependencySpecifierCommand
		)

		BeforeEach(func() {
			cmd = commands.DependencySpecifierCommand{}
		})

		It("invokes the DependencySpecifier client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeDependencySpecifierClient.GetCallCount()).To(Equal(1))
		})

		Context("when the DependencySpecifier client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeDependencySpecifierClient.GetReturns(expectedErr)
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

		Describe("DependencySpecifierID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "DependencySpecifierID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("dependency-specifier-id"))
			})
		})
	})

	Describe("CreateDependencySpecifierCommand", func() {
		var (
			cmd commands.CreateDependencySpecifierCommand
		)

		BeforeEach(func() {
			cmd = commands.CreateDependencySpecifierCommand{}
		})

		It("invokes the DependencySpecifier client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeDependencySpecifierClient.CreateCallCount()).To(Equal(1))
		})

		Context("when the DependencySpecifier client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeDependencySpecifierClient.CreateReturns(expectedErr)
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
				field = fieldFor(commands.CreateDependencySpecifierCommand{}, "ProductSlug")
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
				field = fieldFor(commands.CreateDependencySpecifierCommand{}, "ReleaseVersion")
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
				field = fieldFor(commands.CreateDependencySpecifierCommand{}, "DependentProductSlug")
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

		Describe("Specifier flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.CreateDependencySpecifierCommand{}, "Specifier")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("u"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("specifier"))
			})
		})
	})

	Describe("DeleteDependencySpecifierCommand", func() {
		var (
			cmd commands.DeleteDependencySpecifierCommand
		)

		BeforeEach(func() {
			cmd = commands.DeleteDependencySpecifierCommand{}
		})

		It("invokes the DependencySpecifier client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeDependencySpecifierClient.DeleteCallCount()).To(Equal(1))
		})

		Context("when the DependencySpecifier client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeDependencySpecifierClient.DeleteReturns(expectedErr)
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

		Describe("DependencySpecifierID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "DependencySpecifierID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("dependency-specifier-id"))
			})
		})
	})
})
