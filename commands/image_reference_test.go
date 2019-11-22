package commands_test

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet/v3"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/commands/imagereference"
)

var _ = Describe("image reference commands", func() {
	var (
		field reflect.StructField

		fakeImageReferenceClient *commandsfakes.FakeImageReferenceClient
	)

	BeforeEach(func() {
		fakeImageReferenceClient = &commandsfakes.FakeImageReferenceClient{}

		commands.NewImageReferenceClient = func(imagereference.PivnetClient) commands.ImageReferenceClient {
			return fakeImageReferenceClient
		}
	})

	Describe("ImageReferencesCommand", func() {
		var (
			cmd commands.ImageReferencesCommand
		)

		BeforeEach(func() {
			cmd = commands.ImageReferencesCommand{}
		})

		It("invokes the ImageReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeImageReferenceClient.ListCallCount()).To(Equal(1))
		})

		Context("when the ImageReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeImageReferenceClient.ListReturns(expectedErr)
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
				field = fieldFor(commands.ImageReferencesCommand{}, "ProductSlug")
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
				field = fieldFor(commands.ImageReferencesCommand{}, "ReleaseVersion")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("r"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-version"))
			})
		})

		Describe("ImageDigest flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ImageReferencesCommand{}, "ImageDigest")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("d"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("digest"))
			})
		})
	})

	Describe("ImageReferenceCommand", func() {
		var (
			cmd commands.ImageReferenceCommand
		)

		BeforeEach(func() {
			cmd = commands.ImageReferenceCommand{}
		})

		It("invokes the ImageReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeImageReferenceClient.GetCallCount()).To(Equal(1))
		})

		Context("when the ImageReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeImageReferenceClient.GetReturns(expectedErr)
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
				field = fieldFor(commands.ImageReferenceCommand{}, "ProductSlug")
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
				field = fieldFor(commands.ImageReferenceCommand{}, "ReleaseVersion")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("r"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-version"))
			})
		})

		Describe("ImageReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ImageReferenceCommand{}, "ImageReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("image-reference-id"))
			})
		})
	})

	Describe("CreateImageReferenceCommand", func() {
		var (
			productSlug string
			name        string
			imagePath   string
			digest      string

			description        string
			docsURL            string
			systemRequirements []string

			cmd commands.CreateImageReferenceCommand
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			name = "some image reference"
			imagePath = "some image path"
			digest = "some digest"

			description = "some description"
			docsURL = "some-docs-url"
			systemRequirements = []string{"system1", "system2"}

			cmd = commands.CreateImageReferenceCommand{
				ProductSlug:        productSlug,
				Name:               name,
				ImagePath:          imagePath,
				Digest:             digest,
				Description:        description,
				DocsURL:            docsURL,
				SystemRequirements: systemRequirements,
			}
		})

		It("invokes the ImageReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			config := pivnet.CreateImageReferenceConfig{
				ProductSlug:        productSlug,
				Name:               name,
				ImagePath:          imagePath,
				Digest:             digest,
				Description:        description,
				DocsURL:            docsURL,
				SystemRequirements: systemRequirements,
			}

			Expect(fakeImageReferenceClient.CreateCallCount()).To(Equal(1))
			Expect(fakeImageReferenceClient.CreateArgsForCall(0)).To(Equal(config))
		})

		Context("when the ImageReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeImageReferenceClient.CreateReturns(expectedErr)
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

		Describe("Name flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Name")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("name"))
			})
		})

		Describe("ImagePath flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ImagePath")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("image-path"))
			})
		})

		Describe("Digest flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Digest")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("digest"))
			})
		})

		Describe("Description flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Description")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("description"))
			})
		})

		Describe("DocsURL flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "DocsURL")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("docs-url"))
			})
		})

		Describe("SystemRequirements flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "SystemRequirements")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("system-requirement"))
			})
		})
	})

	Describe("DeleteImageReferenceCommand", func() {
		var (
			cmd commands.DeleteImageReferenceCommand
		)

		BeforeEach(func() {
			cmd = commands.DeleteImageReferenceCommand{}
		})

		It("invokes the ImageReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeImageReferenceClient.DeleteCallCount()).To(Equal(1))
		})

		Context("when the ImageReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeImageReferenceClient.DeleteReturns(expectedErr)
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
				field = fieldFor(commands.DeleteImageReferenceCommand{}, "ProductSlug")
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

		Describe("ImageReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DeleteImageReferenceCommand{}, "ImageReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("image-reference-id"))
			})
		})
	})

	Describe("UpdateImagereferenceCommand", func() {
		var (
			productSlug      string
			imageReferenceID int

			description        string
			name               string
			docsURL            string
			systemRequirements []string

			cmd commands.UpdateImageReferenceCommand
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			imageReferenceID = 1234

			description = "some description"
			name = "some image reference"
			docsURL = "some-docs-url"
			systemRequirements = []string{"system1", "system2"}

			cmd = commands.UpdateImageReferenceCommand{
				ProductSlug:        productSlug,
				ImageReferenceID:   imageReferenceID,
				Name:               &name,
				Description:        &description,
				DocsURL:            &docsURL,
				SystemRequirements: &systemRequirements,
			}
		})

		It("invokes the ProductFile client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeImageReferenceClient.UpdateCallCount()).To(Equal(1))

			invokedProductSlug,
				invokedImageID,
				invokedName,
				invokedDescription,
				invokedDocsURL,
				invokedSystemRequirements := fakeImageReferenceClient.UpdateArgsForCall(0)

			Expect(invokedImageID).To(Equal(imageReferenceID))
			Expect(invokedProductSlug).To(Equal(productSlug))
			Expect(*invokedName).To(Equal(name))
			Expect(*invokedDescription).To(Equal(description))
			Expect(*invokedDocsURL).To(Equal(docsURL))
			Expect(*invokedSystemRequirements).To(Equal(systemRequirements))
		})

		Context("when the ProductFile client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeImageReferenceClient.UpdateReturns(expectedErr)
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

		Describe("ImageReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ImageReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("image-reference-id"))
			})
		})

		Describe("Name flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Name")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("name"))
			})
		})

		Describe("Description flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Description")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("description"))
			})
		})

		Describe("DocsURL flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "DocsURL")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("docs-url"))
			})
		})

		Describe("SystemRequirements flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "SystemRequirements")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("system-requirement"))
			})
		})
	})

	Describe("AddImageReferenceToReleaseCommand", func() {
		var (
			cmd commands.AddImageReferenceToReleaseCommand
		)

		BeforeEach(func() {
			cmd = commands.AddImageReferenceToReleaseCommand{}
		})

		It("invokes the ImageReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeImageReferenceClient.AddToReleaseCallCount()).To(Equal(1))
		})

		Context("when the ImageReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeImageReferenceClient.AddToReleaseReturns(expectedErr)
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

		Describe("ImageReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ImageReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("image-reference-id"))
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
	})

	Describe("RemoveImageReferenceFromReleaseCommand", func() {
		var (
			cmd commands.RemoveImageReferenceFromReleaseCommand
		)

		BeforeEach(func() {
			cmd = commands.RemoveImageReferenceFromReleaseCommand{}
		})

		It("invokes the ImageReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeImageReferenceClient.RemoveFromReleaseCallCount()).To(Equal(1))
		})

		Context("when the ImageReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeImageReferenceClient.RemoveFromReleaseReturns(expectedErr)
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

		Describe("ImageReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ImageReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("image-reference-id"))
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
	})
})
