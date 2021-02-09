package commands_test

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet/v7"
	"github.com/pivotal-cf/pivnet-cli/v2/commands"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/artifactreference"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/commandsfakes"
)

var _ = Describe("artifact reference commands", func() {
	var (
		field reflect.StructField

		fakeArtifactReferenceClient *commandsfakes.FakeArtifactReferenceClient
	)

	BeforeEach(func() {
		fakeArtifactReferenceClient = &commandsfakes.FakeArtifactReferenceClient{}

		commands.NewArtifactReferenceClient = func(artifactreference.PivnetClient) commands.ArtifactReferenceClient {
			return fakeArtifactReferenceClient
		}
	})

	Describe("ArtifactReferencesCommand", func() {
		var (
			cmd commands.ArtifactReferencesCommand
		)

		BeforeEach(func() {
			cmd = commands.ArtifactReferencesCommand{}
		})

		It("invokes the ArtifactReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeArtifactReferenceClient.ListCallCount()).To(Equal(1))
		})

		Context("when the ArtifactReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeArtifactReferenceClient.ListReturns(expectedErr)
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
				field = fieldFor(commands.ArtifactReferencesCommand{}, "ProductSlug")
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
				field = fieldFor(commands.ArtifactReferencesCommand{}, "ReleaseVersion")
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

		Describe("ArtifactDigest flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ArtifactReferencesCommand{}, "ArtifactDigest")
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

	Describe("ArtifactReferenceCommand", func() {
		var (
			cmd commands.ArtifactReferenceCommand
		)

		BeforeEach(func() {
			cmd = commands.ArtifactReferenceCommand{}
		})

		It("invokes the ArtifactReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeArtifactReferenceClient.GetCallCount()).To(Equal(1))
		})

		Context("when the ArtifactReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeArtifactReferenceClient.GetReturns(expectedErr)
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
				field = fieldFor(commands.ArtifactReferenceCommand{}, "ProductSlug")
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
				field = fieldFor(commands.ArtifactReferenceCommand{}, "ReleaseVersion")
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

		Describe("ArtifactReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ArtifactReferenceCommand{}, "ArtifactReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("artifact-reference-id"))
			})
		})
	})

	Describe("CreateArtifactReferenceCommand", func() {
		var (
			productSlug  string
			name         string
			artifactPath string
			digest       string

			description        string
			docsURL            string
			systemRequirements []string

			cmd commands.CreateArtifactReferenceCommand
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			name = "some artifact reference"
			artifactPath = "some artifact path"
			digest = "some digest"

			description = "some description"
			docsURL = "some-docs-url"
			systemRequirements = []string{"system1", "system2"}

			cmd = commands.CreateArtifactReferenceCommand{
				ProductSlug:        productSlug,
				Name:               name,
				ArtifactPath:       artifactPath,
				Digest:             digest,
				Description:        description,
				DocsURL:            docsURL,
				SystemRequirements: systemRequirements,
			}
		})

		It("invokes the ArtifactReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			config := pivnet.CreateArtifactReferenceConfig{
				ProductSlug:        productSlug,
				Name:               name,
				ArtifactPath:       artifactPath,
				Digest:             digest,
				Description:        description,
				DocsURL:            docsURL,
				SystemRequirements: systemRequirements,
			}

			Expect(fakeArtifactReferenceClient.CreateCallCount()).To(Equal(1))
			Expect(fakeArtifactReferenceClient.CreateArgsForCall(0)).To(Equal(config))
		})

		Context("when the ArtifactReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeArtifactReferenceClient.CreateReturns(expectedErr)
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

		Describe("ArtifactPath flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ArtifactPath")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("artifact-path"))
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

	Describe("DeleteArtifactReferenceCommand", func() {
		var (
			cmd commands.DeleteArtifactReferenceCommand
		)

		BeforeEach(func() {
			cmd = commands.DeleteArtifactReferenceCommand{}
		})

		It("invokes the ArtifactReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeArtifactReferenceClient.DeleteCallCount()).To(Equal(1))
		})

		Context("when the ArtifactReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeArtifactReferenceClient.DeleteReturns(expectedErr)
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
				field = fieldFor(commands.DeleteArtifactReferenceCommand{}, "ProductSlug")
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

		Describe("ArtifactReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DeleteArtifactReferenceCommand{}, "ArtifactReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("artifact-reference-id"))
			})
		})
	})

	Describe("UpdateArtifactreferenceCommand", func() {
		var (
			productSlug         string
			artifactReferenceID int

			description        string
			name               string
			docsURL            string
			systemRequirements []string

			cmd commands.UpdateArtifactReferenceCommand
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			artifactReferenceID = 1234

			description = "some description"
			name = "some artifact reference"
			docsURL = "some-docs-url"
			systemRequirements = []string{"system1", "system2"}

			cmd = commands.UpdateArtifactReferenceCommand{
				ProductSlug:         productSlug,
				ArtifactReferenceID: artifactReferenceID,
				Name:                &name,
				Description:         &description,
				DocsURL:             &docsURL,
				SystemRequirements:  &systemRequirements,
			}
		})

		It("invokes the ArtifactReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeArtifactReferenceClient.UpdateCallCount()).To(Equal(1))

			invokedProductSlug,
			invokedArtifactID,
			invokedName,
			invokedDescription,
			invokedDocsURL,
			invokedSystemRequirements := fakeArtifactReferenceClient.UpdateArgsForCall(0)

			Expect(invokedArtifactID).To(Equal(artifactReferenceID))
			Expect(invokedProductSlug).To(Equal(productSlug))
			Expect(*invokedName).To(Equal(name))
			Expect(*invokedDescription).To(Equal(description))
			Expect(*invokedDocsURL).To(Equal(docsURL))
			Expect(*invokedSystemRequirements).To(Equal(systemRequirements))
		})

		Context("when the ArtifactReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeArtifactReferenceClient.UpdateReturns(expectedErr)
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

		Describe("ArtifactReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ArtifactReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("artifact-reference-id"))
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

	Describe("AddArtifactReferenceToReleaseCommand", func() {
		var (
			cmd commands.AddArtifactReferenceToReleaseCommand
		)

		BeforeEach(func() {
			cmd = commands.AddArtifactReferenceToReleaseCommand{}
		})

		It("invokes the ArtifactReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeArtifactReferenceClient.AddToReleaseCallCount()).To(Equal(1))
		})

		Context("when the ArtifactReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeArtifactReferenceClient.AddToReleaseReturns(expectedErr)
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

		Describe("ArtifactReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ArtifactReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("artifact-reference-id"))
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

	Describe("RemoveArtifactReferenceFromReleaseCommand", func() {
		var (
			cmd commands.RemoveArtifactReferenceFromReleaseCommand
		)

		BeforeEach(func() {
			cmd = commands.RemoveArtifactReferenceFromReleaseCommand{}
		})

		It("invokes the ArtifactReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeArtifactReferenceClient.RemoveFromReleaseCallCount()).To(Equal(1))
		})

		Context("when the ArtifactReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeArtifactReferenceClient.RemoveFromReleaseReturns(expectedErr)
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

		Describe("ArtifactReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ArtifactReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("artifact-reference-id"))
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
