package imagereference_test

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
	"github.com/pivotal-cf/pivnet-cli/commands/imagereference"
	"github.com/pivotal-cf/pivnet-cli/commands/imagereference/imagereferencefakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("imagereference commands", func() {
	var (
		l                logger.Logger
		fakePivnetClient *imagereferencefakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer
		logBuffer bytes.Buffer

		imageReferences []pivnet.ImageReference

		client *imagereference.ImageReferenceClient
	)

	BeforeEach(func() {
		infoLogger := log.New(GinkgoWriter, "", 0)
		debugLogger := log.New(GinkgoWriter, "", 0)
		l = logshim.NewLogShim(infoLogger, debugLogger, true)

		fakePivnetClient = &imagereferencefakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}
		logBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		imageReferences = []pivnet.ImageReference{
			{
				ID:                 1234,
				Name:               "my name",
				ImagePath:          "my/path:123",
				Description:        "my description",
				DocsURL:            "my.docs.url",
				Digest:             "sha256:mydigest",
				SystemRequirements: []string{"requirement1", "requirement2"},
			},
		}

		client = imagereference.NewImageReferenceClient(
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

			fakePivnetClient.ImageReferencesReturns(imageReferences, nil)
		})

		It("lists all ImageReferences", func() {
			err := client.List(productSlug, releaseVersion, digest)
			Expect(err).NotTo(HaveOccurred())

			var returnedImageReferences []pivnet.ImageReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedImageReferences)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedImageReferences).To(Equal(imageReferences))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("imageReferences error")
				fakePivnetClient.ImageReferencesReturns(nil, expectedErr)
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
				fakePivnetClient.ImageReferencesForReleaseReturns(imageReferences, nil)
			})

			It("lists all ImageReferences", func() {
				err := client.List(productSlug, releaseVersion, digest)
				Expect(err).NotTo(HaveOccurred())

				var returnedImageReferences []pivnet.ImageReference
				err = json.Unmarshal(outBuffer.Bytes(), &returnedImageReferences)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedImageReferences).To(Equal(imageReferences))
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
					expectedErr = errors.New("imageReferences error")
					fakePivnetClient.ImageReferencesForReleaseReturns(nil, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.List(productSlug, releaseVersion, digest)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})
		})

		Context("when an image digest is provided", func() {
			BeforeEach(func() {
				digest = "sha256:digest"
				imageReferences = []pivnet.ImageReference{
					{
						ID:                 1234,
						Name:               "my name",
						ImagePath:          "my/path:123",
						Description:        "my description",
						DocsURL:            "my.docs.url",
						Digest:             "sha256:mydigest",
						SystemRequirements: []string{"requirement1", "requirement2"},
						ReleaseVersions:    []string{"1.0.0", "1.2.3"},
					},
					{
						ID:                 9876,
						Name:               "other image name",
						ImagePath:          "other image/path:123",
						Description:        "other image description",
						DocsURL:            "other image.docs.url",
						Digest:             "sha256:mydigest",
						SystemRequirements: []string{"requirement3", "requirement4"},
						ReleaseVersions:    []string{},
					},
				}

				fakePivnetClient.ImageReferencesForDigestReturns(imageReferences, nil)
			})

			It("invoke proper list command on client", func() {
				Expect(fakePivnetClient.ImageReferencesForDigestCallCount()).To(Equal(0))

				err := client.List(productSlug, releaseVersion, digest)

				Expect(err).ToNot(HaveOccurred())
				Expect(fakePivnetClient.ImageReferencesForDigestCallCount()).To(Equal(1))

				invokedProductSlug, invokedDigest := fakePivnetClient.ImageReferencesForDigestArgsForCall(0)

				Expect(invokedProductSlug).To(Equal(productSlug))
				Expect(invokedDigest).To(Equal(digest))

				var returnedImageReferences []pivnet.ImageReference
				err = json.Unmarshal(outBuffer.Bytes(), &returnedImageReferences)
				Expect(err).ToNot(HaveOccurred())

				expectedImageReferences := imageReferences
				expectedImageReferences[1].ReleaseVersions = nil
				Expect(returnedImageReferences).To(Equal(expectedImageReferences))
			})

			Context("when there is an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("imageReferences error")
					fakePivnetClient.ImageReferencesForDigestReturns(nil, expectedErr)
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
			productSlug      string
			releaseVersion   string
			imageReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = ""
			imageReferenceID = imageReferences[0].ID

			fakePivnetClient.ImageReferenceReturns(imageReferences[0], nil)
		})

		It("gets ImageReference", func() {
			err := client.Get(productSlug, releaseVersion, imageReferenceID)
			Expect(err).NotTo(HaveOccurred())

			var returnedImageReference pivnet.ImageReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedImageReference)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedImageReference).To(Equal(imageReferences[0]))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("imagereference error")
				fakePivnetClient.ImageReferenceReturns(pivnet.ImageReference{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(productSlug, releaseVersion, imageReferenceID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when release version is not empty", func() {
			BeforeEach(func() {
				releaseVersion = "some-release-version"
				fakePivnetClient.ImageReferenceForReleaseReturns(imageReferences[0], nil)
			})

			It("gets ImageReference", func() {
				err := client.Get(productSlug, releaseVersion, imageReferenceID)
				Expect(err).NotTo(HaveOccurred())

				var returnedImageReference pivnet.ImageReference
				err = json.Unmarshal(outBuffer.Bytes(), &returnedImageReference)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedImageReference).To(Equal(imageReferences[0]))
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
					err := client.Get(productSlug, releaseVersion, imageReferenceID)
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
					expectedErr = errors.New("imageReferences error")
					fakePivnetClient.ImageReferenceForReleaseReturns(pivnet.ImageReference{}, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.Get(productSlug, releaseVersion, imageReferenceID)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})
		})
	})

	Describe("Update", func() {
		var (
			imageReferenceID int
			productSlug      string

			existingName               string
			existingDescription        string
			existingDocsURL            string
			existingSystemRequirements []string

			name               string
			description        string
			docsURL            string
			systemRequirements []string

			existingImageReference pivnet.ImageReference
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			imageReferenceID = imageReferences[0].ID

			existingName = "some-name"
			existingDescription = "some-description"
			existingDocsURL = "some-url.net"
			existingSystemRequirements = []string{"1", "2"}

			name = "some-new-name"
			description = "some-new-description"
			docsURL = "some-new-url.net"
			systemRequirements = []string{"3", "4"}

			existingImageReference = pivnet.ImageReference{
				ID:                 imageReferenceID,
				Name:               existingName,
				Description:        existingDescription,
				DocsURL:            existingDocsURL,
				SystemRequirements: existingSystemRequirements,
			}

			fakePivnetClient.ImageReferenceReturns(existingImageReference, nil)
			fakePivnetClient.UpdateImageReferenceReturns(imageReferences[0], nil)
		})

		It("updates ImageReference", func() {
			err := client.Update(
				productSlug,
				imageReferenceID,
				&name,
				&description,
				&docsURL,
				&systemRequirements,
			)
			Expect(err).NotTo(HaveOccurred())

			var returnedImageReference pivnet.ImageReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedImageReference)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedImageReference).To(Equal(imageReferences[0]))

			invokedProductSlug, invokedImageReference := fakePivnetClient.UpdateImageReferenceArgsForCall(0)
			Expect(invokedProductSlug).To(Equal(productSlug))
			Expect(invokedImageReference.ID).To(Equal(imageReferenceID))
			Expect(invokedImageReference.Name).To(Equal(name))
			Expect(invokedImageReference.Description).To(Equal(description))
			Expect(invokedImageReference.DocsURL).To(Equal(docsURL))
			Expect(invokedImageReference.SystemRequirements).To(Equal(systemRequirements))
		})

		Context("when optional fields are nil", func() {
			It("updates Image Reference with previous values", func() {
				err := client.Update(
					productSlug,
					imageReferenceID,
					nil,
					nil,
					nil,
					nil,
				)
				Expect(err).NotTo(HaveOccurred())

				var returnedImageReference pivnet.ImageReference
				err = json.Unmarshal(outBuffer.Bytes(), &returnedImageReference)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedImageReference).To(Equal(imageReferences[0]))

				invokedProductSlug, invokedProductFile := fakePivnetClient.UpdateImageReferenceArgsForCall(0)
				Expect(invokedProductSlug).To(Equal(productSlug))
				Expect(invokedProductFile.ID).To(Equal(imageReferenceID))
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
				expectedErr = errors.New("imageReference error")
				fakePivnetClient.UpdateImageReferenceReturns(pivnet.ImageReference{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Update(
					productSlug,
					imageReferenceID,
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
			config pivnet.CreateImageReferenceConfig
		)

		BeforeEach(func() {
			config = pivnet.CreateImageReferenceConfig{
				Name: "some-name",
			}

			fakePivnetClient.CreateImageReferenceReturns(imageReferences[0], nil)
		})

		It("creates ImageReference", func() {
			err := client.Create(config)
			Expect(err).NotTo(HaveOccurred())

			var returnedImageReference pivnet.ImageReference
			err = json.Unmarshal(outBuffer.Bytes(), &returnedImageReference)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedImageReference).To(Equal(imageReferences[0]))
			Expect(fakePivnetClient.CreateImageReferenceArgsForCall(0)).To(Equal(config))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("imagereference error")
				fakePivnetClient.CreateImageReferenceReturns(pivnet.ImageReference{}, expectedErr)
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
			imageReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			imageReferenceID = imageReferences[0].ID

			fakePivnetClient.DeleteImageReferenceReturns(imageReferences[0], nil)
		})

		It("deletes ImageReference", func() {
			err := client.Delete(productSlug, imageReferenceID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("imagereference error")
				fakePivnetClient.DeleteImageReferenceReturns(pivnet.ImageReference{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Delete(productSlug, imageReferenceID)
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
			imageReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			imageReferenceID = imageReferences[0].ID

			fakePivnetClient.AddImageReferenceToReleaseReturns(nil)
		})

		It("adds ImageReference to release", func() {
			err := client.AddToRelease(
				productSlug,
				imageReferenceID,
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
					imageReferenceID,
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
				expectedErr = errors.New("image reference error")
				fakePivnetClient.AddImageReferenceToReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddToRelease(
					productSlug,
					imageReferenceID,
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
			imageReferenceID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			imageReferenceID = imageReferences[0].ID

			fakePivnetClient.RemoveImageReferenceFromReleaseReturns(nil)
		})

		It("removes ImageReference from release", func() {
			err := client.RemoveFromRelease(
				productSlug,
				imageReferenceID,
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
					imageReferenceID,
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
				expectedErr = errors.New("image reference error")
				fakePivnetClient.RemoveImageReferenceFromReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveFromRelease(
					productSlug,
					imageReferenceID,
					releaseVersion,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
