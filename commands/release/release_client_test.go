package release_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-cli/commands/release"
	"github.com/pivotal-cf/pivnet-cli/commands/release/releasefakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("release commands", func() {
	var (
		fakePivnetClient *releasefakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		releases []pivnet.Release

		client *release.ReleaseClient
	)

	BeforeEach(func() {
		fakePivnetClient = &releasefakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		releases = []pivnet.Release{
			{
				ID: 1234,
			},
			{
				ID: 2345,
			},
		}

		client = release.NewReleaseClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
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

			fakePivnetClient.ReleasesForProductSlugReturns(releases, nil)
		})

		It("lists all Releases", func() {
			err := client.List(productSlug)
			Expect(err).NotTo(HaveOccurred())

			var returnedReleases []pivnet.Release
			err = json.Unmarshal(outBuffer.Bytes(), &returnedReleases)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedReleases).To(Equal(releases))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleasesForProductSlugReturns(nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.List(productSlug)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Get", func() {
		var (
			productSlug    string
			releaseVersion string
			releaseID      int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "some-release-version"
			releaseID = releases[0].ID

			fakePivnetClient.ReleaseForVersionReturns(releases[0], nil)
		})

		It("gets Release", func() {
			err := client.Get(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())

			var returnedRelease pivnet.Release
			err = json.Unmarshal(outBuffer.Bytes(), &returnedRelease)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedRelease).To(Equal(releases[0]))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("release error")
				fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Create", func() {
		var (
			productSlug    string
			releaseVersion string
			releaseType    string
			eulaSlug       string

			validEULAs        []pivnet.EULA
			validReleaseTypes []pivnet.ReleaseType
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "some-release-version"
			releaseType = "some-release-type"
			eulaSlug = "some-eula-slug"

			validEULAs = []pivnet.EULA{
				{
					Slug: eulaSlug,
				},
				{
					Slug: "some-other-eula-slug",
				},
			}

			validReleaseTypes = []pivnet.ReleaseType{
				pivnet.ReleaseType(releaseType),
				"some-other-release-type",
			}

			fakePivnetClient.ReleaseForVersionReturns(releases[0], nil)
			fakePivnetClient.CreateReleaseReturns(releases[0], nil)
			fakePivnetClient.EULAsReturns(validEULAs, nil)
			fakePivnetClient.ReleaseTypesReturns(validReleaseTypes, nil)
		})

		It("creates Release", func() {
			err := client.Create(
				productSlug,
				releaseVersion,
				releaseType,
				eulaSlug,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("release error")
				fakePivnetClient.CreateReleaseReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Create(
					productSlug,
					releaseVersion,
					releaseType,
					eulaSlug,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error creating release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("release error")
				fakePivnetClient.CreateReleaseReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Create(
					productSlug,
					releaseVersion,
					releaseType,
					eulaSlug,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error getting eulas", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("eula error")
				fakePivnetClient.EULAsReturns(nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Create(
					productSlug,
					releaseVersion,
					releaseType,
					eulaSlug,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when the EULA is not valid", func() {
			BeforeEach(func() {
				validEULAs = []pivnet.EULA{
					{
						Slug: "a-different-eula-slug",
					},
					{
						Slug: "some-other-eula-slug",
					},
				}

				fakePivnetClient.EULAsReturns(validEULAs, nil)
			})

			It("invokes the error handler", func() {
				err := client.Create(
					productSlug,
					releaseVersion,
					releaseType,
					eulaSlug,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				actualErr := fakeErrorHandler.HandleErrorArgsForCall(0)
				Expect(actualErr.Error()).To(ContainSubstring(eulaSlug))
				Expect(actualErr.Error()).To(ContainSubstring(validEULAs[0].Slug))
				Expect(actualErr.Error()).To(ContainSubstring(validEULAs[1].Slug))
			})
		})

		Context("when there is an error getting release types", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("release type error")
				fakePivnetClient.ReleaseTypesReturns(nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Create(
					productSlug,
					releaseVersion,
					releaseType,
					eulaSlug,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when the release type is not valid", func() {
			BeforeEach(func() {
				validReleaseTypes = []pivnet.ReleaseType{
					"a-different-release-type",
					"some-other-release-type",
				}

				fakePivnetClient.ReleaseTypesReturns(validReleaseTypes, nil)
			})

			It("invokes the error handler", func() {
				err := client.Create(
					productSlug,
					releaseVersion,
					releaseType,
					eulaSlug,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				actualErr := fakeErrorHandler.HandleErrorArgsForCall(0)
				Expect(actualErr.Error()).To(ContainSubstring(releaseType))
				Expect(actualErr.Error()).To(ContainSubstring(string(validReleaseTypes[0])))
				Expect(actualErr.Error()).To(ContainSubstring(string(validReleaseTypes[1])))
			})
		})

	})

	Describe("Delete", func() {
		var (
			productSlug    string
			releaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = releases[0].Version

			fakePivnetClient.ReleaseForVersionReturns(releases[0], nil)
			fakePivnetClient.DeleteReleaseReturns(nil)
		})

		It("deletes Release", func() {
			err := client.Delete(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("release error")
				fakePivnetClient.DeleteReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Delete(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error getting release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("release error")
				fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Delete(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
