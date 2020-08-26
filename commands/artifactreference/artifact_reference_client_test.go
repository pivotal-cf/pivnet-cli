package artifactreference_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet/v6"
	"github.com/pivotal-cf/go-pivnet/v6/logger"
	"github.com/pivotal-cf/go-pivnet/v6/logshim"
	"github.com/pivotal-cf/pivnet-cli/commands/artifactreference"
	"github.com/pivotal-cf/pivnet-cli/commands/artifactreference/artifactreferencefakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("artifactreference commands", func() {
	var (
		l                logger.Logger
		fakePivnetClient *artifactreferencefakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer
		logBuffer bytes.Buffer

		artifactReferences []pivnet.ArtifactReference

		client *artifactreference.ArtifactReferenceClient
	)

	BeforeEach(func() {
		infoLogger := log.New(GinkgoWriter, "", 0)
		debugLogger := log.New(GinkgoWriter, "", 0)
		l = logshim.NewLogShim(infoLogger, debugLogger, true)

		fakePivnetClient = &artifactreferencefakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}
		logBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		artifactReferences = []pivnet.ArtifactReference{
			{
				ID:                 1234,
				Name:               "my name",
				ArtifactPath:       "my/path:123",
				Description:        "my description",
				DocsURL:            "my.docs.url",
				Digest:             "sha256:mydigest",
				SystemRequirements: []string{"requirement1", "requirement2"},
			},
		}

		client = artifactreference.NewArtifactReferenceClient(
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
			digest         string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = ""
			digest = ""

			fakePivnetClient.ArtifactReferencesReturns(artifactReferences, nil)
		})

		It("lists all ArtifactReferences", func() {
			err := client.List(productSlug, releaseVersion, digest)
			Expect(err).NotTo(HaveOccurred())

			var returnedArtifactReferences []pivnet.ArtifactReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedArtifactReferences)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedArtifactReferences).To(Equal(artifactReferences))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("artifactReferences error")
				fakePivnetClient.ArtifactReferencesReturns(nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.List(productSlug, releaseVersion, digest)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when release version is not empty", func() {
			BeforeEach(func() {
				releaseVersion = "some-release-version"
				fakePivnetClient.ArtifactReferencesForReleaseReturns(artifactReferences, nil)
			})

			It("lists all ArtifactReferences", func() {
				err := client.List(productSlug, releaseVersion, digest)
				Expect(err).NotTo(HaveOccurred())

				var returnedArtifactReferences []pivnet.ArtifactReference
				err = json.Unmarshal(outBuffer.Bytes(), &returnedArtifactReferences)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedArtifactReferences).To(Equal(artifactReferences))
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
					err := client.List(productSlug, releaseVersion, digest)
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
					expectedErr = errors.New("artifactReferences error")
					fakePivnetClient.ArtifactReferencesForReleaseReturns(nil, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.List(productSlug, releaseVersion, digest)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})
		})

		Context("when an artifact digest is provided", func() {
			BeforeEach(func() {
				digest = "sha256:digest"
				artifactReferences = []pivnet.ArtifactReference{
					{
						ID:                 1234,
						Name:               "my name",
						ArtifactPath:       "my/path:123",
						Description:        "my description",
						DocsURL:            "my.docs.url",
						Digest:             "sha256:mydigest",
						SystemRequirements: []string{"requirement1", "requirement2"},
						ReleaseVersions:    []string{"1.0.0", "1.2.3"},
					},
					{
						ID:                 9876,
						Name:               "other artifact name",
						ArtifactPath:       "other artifact/path:123",
						Description:        "other artifact description",
						DocsURL:            "other artifact.docs.url",
						Digest:             "sha256:mydigest",
						SystemRequirements: []string{"requirement3", "requirement4"},
						ReleaseVersions:    []string{},
					},
				}

				fakePivnetClient.ArtifactReferencesForDigestReturns(artifactReferences, nil)
			})

			It("invoke proper list command on client", func() {
				Expect(fakePivnetClient.ArtifactReferencesForDigestCallCount()).To(Equal(0))

				err := client.List(productSlug, releaseVersion, digest)

				Expect(err).ToNot(HaveOccurred())
				Expect(fakePivnetClient.ArtifactReferencesForDigestCallCount()).To(Equal(1))

				invokedProductSlug, invokedDigest := fakePivnetClient.ArtifactReferencesForDigestArgsForCall(0)

				Expect(invokedProductSlug).To(Equal(productSlug))
				Expect(invokedDigest).To(Equal(digest))

				var returnedArtifactReferences []pivnet.ArtifactReference
				err = json.Unmarshal(outBuffer.Bytes(), &returnedArtifactReferences)
				Expect(err).ToNot(HaveOccurred())

				expectedArtifactReferences := artifactReferences
				expectedArtifactReferences[1].ReleaseVersions = nil
				Expect(returnedArtifactReferences).To(Equal(expectedArtifactReferences))
			})

			Context("when there is an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("artifactReferences error")
					fakePivnetClient.ArtifactReferencesForDigestReturns(nil, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.List(productSlug, releaseVersion, digest)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})
		})
	})

	Describe("Get", func() {
		var (
			productSlug         string
			releaseVersion      string
			artifactReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = ""
			artifactReferenceID = artifactReferences[0].ID

			fakePivnetClient.ArtifactReferenceReturns(artifactReferences[0], nil)
		})

		It("gets ArtifactReference", func() {
			err := client.Get(productSlug, releaseVersion, artifactReferenceID)
			Expect(err).NotTo(HaveOccurred())

			var returnedArtifactReference pivnet.ArtifactReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedArtifactReference)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedArtifactReference).To(Equal(artifactReferences[0]))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("artifactreference error")
				fakePivnetClient.ArtifactReferenceReturns(pivnet.ArtifactReference{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(productSlug, releaseVersion, artifactReferenceID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when release version is not empty", func() {
			BeforeEach(func() {
				releaseVersion = "some-release-version"
				fakePivnetClient.ArtifactReferenceForReleaseReturns(artifactReferences[0], nil)
			})

			It("gets ArtifactReference", func() {
				err := client.Get(productSlug, releaseVersion, artifactReferenceID)
				Expect(err).NotTo(HaveOccurred())

				var returnedArtifactReference pivnet.ArtifactReference
				err = json.Unmarshal(outBuffer.Bytes(), &returnedArtifactReference)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedArtifactReference).To(Equal(artifactReferences[0]))
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
					err := client.Get(productSlug, releaseVersion, artifactReferenceID)
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
					expectedErr = errors.New("artifactReferences error")
					fakePivnetClient.ArtifactReferenceForReleaseReturns(pivnet.ArtifactReference{}, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.Get(productSlug, releaseVersion, artifactReferenceID)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})
		})
	})

	Describe("Update", func() {
		var (
			artifactReferenceID int
			productSlug         string

			existingName               string
			existingDescription        string
			existingDocsURL            string
			existingSystemRequirements []string

			name               string
			description        string
			docsURL            string
			systemRequirements []string

			existingArtifactReference pivnet.ArtifactReference
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			artifactReferenceID = artifactReferences[0].ID

			existingName = "some-name"
			existingDescription = "some-description"
			existingDocsURL = "some-url.net"
			existingSystemRequirements = []string{"1", "2"}

			name = "some-new-name"
			description = "some-new-description"
			docsURL = "some-new-url.net"
			systemRequirements = []string{"3", "4"}

			existingArtifactReference = pivnet.ArtifactReference{
				ID:                 artifactReferenceID,
				Name:               existingName,
				Description:        existingDescription,
				DocsURL:            existingDocsURL,
				SystemRequirements: existingSystemRequirements,
			}

			fakePivnetClient.ArtifactReferenceReturns(existingArtifactReference, nil)
			fakePivnetClient.UpdateArtifactReferenceReturns(artifactReferences[0], nil)
		})

		It("updates ArtifactReference", func() {
			err := client.Update(
				productSlug,
				artifactReferenceID,
				&name,
				&description,
				&docsURL,
				&systemRequirements,
			)
			Expect(err).NotTo(HaveOccurred())

			var returnedArtifactReference pivnet.ArtifactReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedArtifactReference)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedArtifactReference).To(Equal(artifactReferences[0]))

			invokedProductSlug, invokedArtifactReference := fakePivnetClient.UpdateArtifactReferenceArgsForCall(0)
			Expect(invokedProductSlug).To(Equal(productSlug))
			Expect(invokedArtifactReference.ID).To(Equal(artifactReferenceID))
			Expect(invokedArtifactReference.Name).To(Equal(name))
			Expect(invokedArtifactReference.Description).To(Equal(description))
			Expect(invokedArtifactReference.DocsURL).To(Equal(docsURL))
			Expect(invokedArtifactReference.SystemRequirements).To(Equal(systemRequirements))
		})

		Context("when optional fields are nil", func() {
			It("updates Artifact Reference with previous values", func() {
				err := client.Update(
					productSlug,
					artifactReferenceID,
					nil,
					nil,
					nil,
					nil,
				)
				Expect(err).NotTo(HaveOccurred())

				var returnedArtifactReference pivnet.ArtifactReference
				err = json.Unmarshal(outBuffer.Bytes(), &returnedArtifactReference)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedArtifactReference).To(Equal(artifactReferences[0]))

				invokedProductSlug, invokedProductFile := fakePivnetClient.UpdateArtifactReferenceArgsForCall(0)
				Expect(invokedProductSlug).To(Equal(productSlug))
				Expect(invokedProductFile.ID).To(Equal(artifactReferenceID))
				Expect(invokedProductFile.Name).To(Equal(existingName))
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
				expectedErr = errors.New("artifactReference error")
				fakePivnetClient.UpdateArtifactReferenceReturns(pivnet.ArtifactReference{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Update(
					productSlug,
					artifactReferenceID,
					&name,
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
			config pivnet.CreateArtifactReferenceConfig
		)

		BeforeEach(func() {
			config = pivnet.CreateArtifactReferenceConfig{
				Name: "some-name",
			}

			fakePivnetClient.CreateArtifactReferenceReturns(artifactReferences[0], nil)
		})

		It("creates ArtifactReference", func() {
			err := client.Create(config)
			Expect(err).NotTo(HaveOccurred())

			var returnedArtifactReference pivnet.ArtifactReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedArtifactReference)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedArtifactReference).To(Equal(artifactReferences[0]))
			Expect(fakePivnetClient.CreateArtifactReferenceArgsForCall(0)).To(Equal(config))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("artifactreference error")
				fakePivnetClient.CreateArtifactReferenceReturns(pivnet.ArtifactReference{}, expectedErr)
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
			productSlug         string
			artifactReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			artifactReferenceID = artifactReferences[0].ID

			fakePivnetClient.DeleteArtifactReferenceReturns(artifactReferences[0], nil)
		})

		It("deletes ArtifactReference", func() {
			err := client.Delete(productSlug, artifactReferenceID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("artifactreference error")
				fakePivnetClient.DeleteArtifactReferenceReturns(pivnet.ArtifactReference{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Delete(productSlug, artifactReferenceID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("AddToRelease", func() {
		var (
			productSlug         string
			releaseVersion      string
			artifactReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			artifactReferenceID = artifactReferences[0].ID

			fakePivnetClient.AddArtifactReferenceToReleaseReturns(nil)
		})

		It("adds ArtifactReference to release", func() {
			err := client.AddToRelease(
				productSlug,
				artifactReferenceID,
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
					artifactReferenceID,
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
				expectedErr = errors.New("artifact reference error")
				fakePivnetClient.AddArtifactReferenceToReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddToRelease(
					productSlug,
					artifactReferenceID,
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
			productSlug         string
			releaseVersion      string
			artifactReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			artifactReferenceID = artifactReferences[0].ID

			fakePivnetClient.RemoveArtifactReferenceFromReleaseReturns(nil)
		})

		It("removes ArtifactReference from release", func() {
			err := client.RemoveFromRelease(
				productSlug,
				artifactReferenceID,
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
					artifactReferenceID,
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
				expectedErr = errors.New("artifact reference error")
				fakePivnetClient.RemoveArtifactReferenceFromReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveFromRelease(
					productSlug,
					artifactReferenceID,
					releaseVersion,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
