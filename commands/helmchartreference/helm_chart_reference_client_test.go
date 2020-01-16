package helmchartreference_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet/v4"
	"github.com/pivotal-cf/go-pivnet/v4/logger"
	"github.com/pivotal-cf/go-pivnet/v4/logshim"
	"github.com/pivotal-cf/pivnet-cli/commands/helmchartreference"
	"github.com/pivotal-cf/pivnet-cli/commands/helmchartreference/helmchartreferencefakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("helmchartreference commands", func() {
	var (
		l                logger.Logger
		fakePivnetClient *helmchartreferencefakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer
		logBuffer bytes.Buffer

		helmChartReferences []pivnet.HelmChartReference

		client *helmchartreference.HelmChartReferenceClient
	)

	BeforeEach(func() {
		infoLogger := log.New(GinkgoWriter, "", 0)
		debugLogger := log.New(GinkgoWriter, "", 0)
		l = logshim.NewLogShim(infoLogger, debugLogger, true)

		fakePivnetClient = &helmchartreferencefakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}
		logBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		helmChartReferences = []pivnet.HelmChartReference{
			{
				ID:                 1234,
				Name:               "myname",
				Version:            "1.2.3",
				Description:        "my description",
				DocsURL:            "my.docs.url",
				SystemRequirements: []string{"requirement1", "requirement2"},
			},
		}

		client = helmchartreference.NewHelmChartReferenceClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			&logBuffer,
			printer.NewPrinter(&outBuffer),
			l,
		)
	})

	Describe("List", func() {
		var (
			productSlug    string
			releaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = ""

			fakePivnetClient.HelmChartReferencesReturns(helmChartReferences, nil)
		})

		It("lists all HelmChartReferences", func() {
			err := client.List(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())

			var returnedHelmChartReferences []pivnet.HelmChartReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedHelmChartReferences)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedHelmChartReferences).To(Equal(helmChartReferences))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("helmChartReferences error")
				fakePivnetClient.HelmChartReferencesReturns(nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.List(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when release version is not empty", func() {
			BeforeEach(func() {
				releaseVersion = "some-release-version"
				fakePivnetClient.HelmChartReferencesForReleaseReturns(helmChartReferences, nil)
			})

			It("lists all HelmChartReferences", func() {
				err := client.List(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				var returnedHelmChartReferences []pivnet.HelmChartReference
				err = json.Unmarshal(outBuffer.Bytes(), &returnedHelmChartReferences)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedHelmChartReferences).To(Equal(helmChartReferences))
			})

			Context("when there is an error getting release", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("releases error")
					fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.List(productSlug, releaseVersion)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})

			Context("when there is an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("helmChartReferences error")
					fakePivnetClient.HelmChartReferencesForReleaseReturns(nil, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.List(productSlug, releaseVersion)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})
		})
	})

	Describe("Get", func() {
		var (
			productSlug      string
			releaseVersion   string
			helmChartReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = ""
			helmChartReferenceID = helmChartReferences[0].ID

			fakePivnetClient.HelmChartReferenceReturns(helmChartReferences[0], nil)
		})

		It("gets HelmChartReference", func() {
			err := client.Get(productSlug, releaseVersion, helmChartReferenceID)
			Expect(err).NotTo(HaveOccurred())

			var returnedHelmChartReference pivnet.HelmChartReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedHelmChartReference)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedHelmChartReference).To(Equal(helmChartReferences[0]))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("helmchartreference error")
				fakePivnetClient.HelmChartReferenceReturns(pivnet.HelmChartReference{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(productSlug, releaseVersion, helmChartReferenceID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when release version is not empty", func() {
			BeforeEach(func() {
				releaseVersion = "some-release-version"
				fakePivnetClient.HelmChartReferenceForReleaseReturns(helmChartReferences[0], nil)
			})

			It("gets HelmChartReference", func() {
				err := client.Get(productSlug, releaseVersion, helmChartReferenceID)
				Expect(err).NotTo(HaveOccurred())

				var returnedHelmChartReference pivnet.HelmChartReference
				err = json.Unmarshal(outBuffer.Bytes(), &returnedHelmChartReference)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedHelmChartReference).To(Equal(helmChartReferences[0]))
			})

			Context("when there is an error getting release", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("releases error")
					fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.Get(productSlug, releaseVersion, helmChartReferenceID)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})

			Context("when there is an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("helmChartReferences error")
					fakePivnetClient.HelmChartReferenceForReleaseReturns(pivnet.HelmChartReference{}, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.Get(productSlug, releaseVersion, helmChartReferenceID)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})
		})
	})

	Describe("Update", func() {
		var (
			helmChartReferenceID int
			productSlug      string

			existingName               string
			existingDescription        string
			existingDocsURL            string
			existingSystemRequirements []string

			description        string
			docsURL            string
			systemRequirements []string

			existingHelmChartReference pivnet.HelmChartReference
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			helmChartReferenceID = helmChartReferences[0].ID

			existingName = "some-name"
			existingDescription = "some-description"
			existingDocsURL = "some-url.net"
			existingSystemRequirements = []string{"1", "2"}

			description = "some-new-description"
			docsURL = "some-new-url.net"
			systemRequirements = []string{"3", "4"}

			existingHelmChartReference = pivnet.HelmChartReference{
				ID:                 helmChartReferenceID,
				Name:               existingName,
				Description:        existingDescription,
				DocsURL:            existingDocsURL,
				SystemRequirements: existingSystemRequirements,
			}

			fakePivnetClient.HelmChartReferenceReturns(existingHelmChartReference, nil)
			fakePivnetClient.UpdateHelmChartReferenceReturns(helmChartReferences[0], nil)
		})

		It("updates HelmChartReference", func() {
			err := client.Update(
				productSlug,
				helmChartReferenceID,
				&description,
				&docsURL,
				&systemRequirements,
			)
			Expect(err).NotTo(HaveOccurred())

			var returnedHelmChartReference pivnet.HelmChartReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedHelmChartReference)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedHelmChartReference).To(Equal(helmChartReferences[0]))

			invokedProductSlug, invokedHelmChartReference := fakePivnetClient.UpdateHelmChartReferenceArgsForCall(0)
			Expect(invokedProductSlug).To(Equal(productSlug))
			Expect(invokedHelmChartReference.ID).To(Equal(helmChartReferenceID))
			Expect(invokedHelmChartReference.Description).To(Equal(description))
			Expect(invokedHelmChartReference.DocsURL).To(Equal(docsURL))
			Expect(invokedHelmChartReference.SystemRequirements).To(Equal(systemRequirements))
		})

		Context("when optional fields are nil", func() {
			It("updates Helm Chart Reference with previous values", func() {
				err := client.Update(
					productSlug,
					helmChartReferenceID,
					nil,
					nil,
					nil,
				)
				Expect(err).NotTo(HaveOccurred())

				var returnedHelmChartReference pivnet.HelmChartReference
				err = json.Unmarshal(outBuffer.Bytes(), &returnedHelmChartReference)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedHelmChartReference).To(Equal(helmChartReferences[0]))

				invokedProductSlug, invokedProductFile := fakePivnetClient.UpdateHelmChartReferenceArgsForCall(0)
				Expect(invokedProductSlug).To(Equal(productSlug))
				Expect(invokedProductFile.ID).To(Equal(helmChartReferenceID))
				Expect(invokedProductFile.Description).To(Equal(existingDescription))
				Expect(invokedProductFile.DocsURL).To(Equal(existingDocsURL))
				Expect(invokedProductFile.SystemRequirements).To(Equal(existingSystemRequirements))
			})
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("helmChartReference error")
				fakePivnetClient.UpdateHelmChartReferenceReturns(pivnet.HelmChartReference{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Update(
					productSlug,
					helmChartReferenceID,
					&description,
					&docsURL,
					&systemRequirements,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Create", func() {
		var (
			config pivnet.CreateHelmChartReferenceConfig
		)

		BeforeEach(func() {
			config = pivnet.CreateHelmChartReferenceConfig{
				Name: "some-name",
			}

			fakePivnetClient.CreateHelmChartReferenceReturns(helmChartReferences[0], nil)
		})

		It("creates HelmChartReference", func() {
			err := client.Create(config)
			Expect(err).NotTo(HaveOccurred())

			var returnedHelmChartReference pivnet.HelmChartReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedHelmChartReference)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedHelmChartReference).To(Equal(helmChartReferences[0]))
			Expect(fakePivnetClient.CreateHelmChartReferenceArgsForCall(0)).To(Equal(config))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("helmchartreference error")
				fakePivnetClient.CreateHelmChartReferenceReturns(pivnet.HelmChartReference{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Create(config)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Delete", func() {
		var (
			productSlug      string
			helmChartReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			helmChartReferenceID = helmChartReferences[0].ID

			fakePivnetClient.DeleteHelmChartReferenceReturns(helmChartReferences[0], nil)
		})

		It("deletes HelmChartReference", func() {
			err := client.Delete(productSlug, helmChartReferenceID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("helmchartreference error")
				fakePivnetClient.DeleteHelmChartReferenceReturns(pivnet.HelmChartReference{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Delete(productSlug, helmChartReferenceID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("AddToRelease", func() {
		var (
			productSlug      string
			releaseVersion   string
			helmChartReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			helmChartReferenceID = helmChartReferences[0].ID

			fakePivnetClient.AddHelmChartReferenceToReleaseReturns(nil)
		})

		It("adds HelmChartReference to release", func() {
			err := client.AddToRelease(
				productSlug,
				helmChartReferenceID,
				releaseVersion,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error getting release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddToRelease(
					productSlug,
					helmChartReferenceID,
					releaseVersion,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("helm chart reference error")
				fakePivnetClient.AddHelmChartReferenceToReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddToRelease(
					productSlug,
					helmChartReferenceID,
					releaseVersion,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("RemoveFromRelease", func() {
		var (
			productSlug      string
			releaseVersion   string
			helmChartReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			helmChartReferenceID = helmChartReferences[0].ID

			fakePivnetClient.RemoveHelmChartReferenceFromReleaseReturns(nil)
		})

		It("removes HelmChartReference from release", func() {
			err := client.RemoveFromRelease(
				productSlug,
				helmChartReferenceID,
				releaseVersion,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error getting release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveFromRelease(
					productSlug,
					helmChartReferenceID,
					releaseVersion,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("helm chart reference error")
				fakePivnetClient.RemoveHelmChartReferenceFromReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveFromRelease(
					productSlug,
					helmChartReferenceID,
					releaseVersion,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
