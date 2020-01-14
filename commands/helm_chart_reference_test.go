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
	"github.com/pivotal-cf/pivnet-cli/commands/helmchartreference"
)

var _ = Describe("helm chart reference commands", func() {
	var (
		field reflect.StructField

		fakeHelmChartReferenceClient *commandsfakes.FakeHelmChartReferenceClient
	)

	BeforeEach(func() {
		fakeHelmChartReferenceClient = &commandsfakes.FakeHelmChartReferenceClient{}

		commands.NewHelmChartReferenceClient = func(helmchartreference.PivnetClient) commands.HelmChartReferenceClient {
			return fakeHelmChartReferenceClient
		}
	})

	Describe("HelmChartReferencesCommand", func() {
		var (
			cmd commands.HelmChartReferencesCommand
		)

		BeforeEach(func() {
			cmd = commands.HelmChartReferencesCommand{}
		})

		It("invokes the HelmChartReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeHelmChartReferenceClient.ListCallCount()).To(Equal(1))
		})

		Context("when the HelmChartReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeHelmChartReferenceClient.ListReturns(expectedErr)
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
				field = fieldFor(commands.HelmChartReferencesCommand{}, "ProductSlug")
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
				field = fieldFor(commands.HelmChartReferencesCommand{}, "ReleaseVersion")
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
	})

	Describe("HelmChartReferenceCommand", func() {
		var (
			cmd commands.HelmChartReferenceCommand
		)

		BeforeEach(func() {
			cmd = commands.HelmChartReferenceCommand{}
		})

		It("invokes the HelmChartReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeHelmChartReferenceClient.GetCallCount()).To(Equal(1))
		})

		Context("when the HelmChartReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeHelmChartReferenceClient.GetReturns(expectedErr)
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
				field = fieldFor(commands.HelmChartReferenceCommand{}, "ProductSlug")
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
				field = fieldFor(commands.HelmChartReferenceCommand{}, "ReleaseVersion")
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

		Describe("HelmChartReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.HelmChartReferenceCommand{}, "HelmChartReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("helm-chart-reference-id"))
			})
		})
	})

	Describe("CreateHelmChartReferenceCommand", func() {
		var (
			productSlug string
			name        string
			version     string

			description        string
			docsURL            string
			systemRequirements []string

			cmd commands.CreateHelmChartReferenceCommand
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			name = "some helm chart reference"
			version = "some version"

			description = "some description"
			docsURL = "some-docs-url"
			systemRequirements = []string{"system1", "system2"}

			cmd = commands.CreateHelmChartReferenceCommand{
				ProductSlug:        productSlug,
				Name:               name,
				Version:            version,
				Description:        description,
				DocsURL:            docsURL,
				SystemRequirements: systemRequirements,
			}
		})

		It("invokes the HelmChartReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			config := pivnet.CreateHelmChartReferenceConfig{
				ProductSlug:        productSlug,
				Name:               name,
				Version:            version,
				Description:        description,
				DocsURL:            docsURL,
				SystemRequirements: systemRequirements,
			}

			Expect(fakeHelmChartReferenceClient.CreateCallCount()).To(Equal(1))
			Expect(fakeHelmChartReferenceClient.CreateArgsForCall(0)).To(Equal(config))
		})

		Context("when the HelmChartReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeHelmChartReferenceClient.CreateReturns(expectedErr)
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

		Describe("Version flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Version")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("version"))
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

	Describe("DeleteHelmChartReferenceCommand", func() {
		var (
			cmd commands.DeleteHelmChartReferenceCommand
		)

		BeforeEach(func() {
			cmd = commands.DeleteHelmChartReferenceCommand{}
		})

		It("invokes the HelmChartReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeHelmChartReferenceClient.DeleteCallCount()).To(Equal(1))
		})

		Context("when the HelmChartReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeHelmChartReferenceClient.DeleteReturns(expectedErr)
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
				field = fieldFor(commands.DeleteHelmChartReferenceCommand{}, "ProductSlug")
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

		Describe("HelmChartReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DeleteHelmChartReferenceCommand{}, "HelmChartReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("helm-chart-reference-id"))
			})
		})
	})

	Describe("UpdateHelmChartReferenceCommand", func() {
		var (
			productSlug      string
			helmChartReferenceID int

			description        string
			docsURL            string
			systemRequirements []string

			cmd commands.UpdateHelmChartReferenceCommand
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			helmChartReferenceID = 1234

			description = "some description"
			docsURL = "some-docs-url"
			systemRequirements = []string{"system1", "system2"}

			cmd = commands.UpdateHelmChartReferenceCommand{
				ProductSlug:        productSlug,
				HelmChartReferenceID:   helmChartReferenceID,
				Description:        &description,
				DocsURL:            &docsURL,
				SystemRequirements: &systemRequirements,
			}
		})

		It("invokes the HelmChartReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeHelmChartReferenceClient.UpdateCallCount()).To(Equal(1))

			invokedProductSlug,
				invokedHelmChartReferenceID,
				invokedDescription,
				invokedDocsURL,
				invokedSystemRequirements := fakeHelmChartReferenceClient.UpdateArgsForCall(0)

			Expect(invokedHelmChartReferenceID).To(Equal(helmChartReferenceID))
			Expect(invokedProductSlug).To(Equal(productSlug))
			Expect(*invokedDescription).To(Equal(description))
			Expect(*invokedDocsURL).To(Equal(docsURL))
			Expect(*invokedSystemRequirements).To(Equal(systemRequirements))
		})

		Context("when the HelmChartReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeHelmChartReferenceClient.UpdateReturns(expectedErr)
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

		Describe("HelmChartReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "HelmChartReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("helm-chart-reference-id"))
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

	Describe("AddHelmChartReferenceToReleaseCommand", func() {
		var (
			cmd commands.AddHelmChartReferenceToReleaseCommand
		)

		BeforeEach(func() {
			cmd = commands.AddHelmChartReferenceToReleaseCommand{}
		})

		It("invokes the HelmChartReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeHelmChartReferenceClient.AddToReleaseCallCount()).To(Equal(1))
		})

		Context("when the HelmChartReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeHelmChartReferenceClient.AddToReleaseReturns(expectedErr)
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

		Describe("HelmChartReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "HelmChartReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("helm-chart-reference-id"))
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

	Describe("RemoveHelmChartReferenceFromReleaseCommand", func() {
		var (
			cmd commands.RemoveHelmChartReferenceFromReleaseCommand
		)

		BeforeEach(func() {
			cmd = commands.RemoveHelmChartReferenceFromReleaseCommand{}
		})

		It("invokes the HelmChartReference client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeHelmChartReferenceClient.RemoveFromReleaseCallCount()).To(Equal(1))
		})

		Context("when the HelmChartReference client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeHelmChartReferenceClient.RemoveFromReleaseReturns(expectedErr)
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

		Describe("HelmChartReferenceID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "HelmChartReferenceID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("helm-chart-reference-id"))
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
