package release_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet/v6"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/release"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/release/releasefakes"
	"github.com/pivotal-cf/pivnet-cli/v2/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/v2/printer"
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
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"

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

	Describe("ListWithLimit", func() {
		var (
			productSlug    string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"

			fakePivnetClient.ReleasesForProductSlugReturns(releases, nil)
		})

		It("lists all Releases", func() {
			Expect(fakePivnetClient.ReleasesForProductSlugCallCount()).To(Equal(0))

			err := client.ListWithLimit(productSlug, "0")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakePivnetClient.ReleasesForProductSlugCallCount()).To(Equal(1))
			slug, params := fakePivnetClient.ReleasesForProductSlugArgsForCall(0)
			Expect(slug).To(Equal(productSlug))
			Expect(params).To(Equal([]pivnet.QueryParameter{{"limit", "0"}}))
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
				err := client.ListWithLimit(productSlug, "")
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
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "some-release-version"

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

	Describe("Update", func() {
		var (
			productSlug    string
			releaseVersion string
			availability   *string
			releaseType    *string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = releases[0].Version
			availability = nil
			releaseType = nil

			fakePivnetClient.ReleaseForVersionReturns(releases[0], nil)
			fakePivnetClient.UpdateReleaseReturns(releases[0], nil)
		})

		It("updates Release", func() {
			err := client.Update(
				productSlug,
				releaseVersion,
				availability,
				releaseType,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("release error")
				fakePivnetClient.UpdateReleaseReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
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
				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when availability is provided", func() {
			var (
				availabilityVar string
			)

			BeforeEach(func() {
				availabilityVar = "all"
				availability = &availabilityVar
			})

			It("sets availability on release", func() {
				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.Availability).To(Equal("All Users"))
			})

			It("correctly sets admins only", func() {
				availabilityVar = "admins"
				availability = &availabilityVar

				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.Availability).To(Equal("Admins Only"))
			})

			It("correctly sets selected user groups only", func() {
				availabilityVar = "selected-user-groups"
				availability = &availabilityVar

				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.Availability).To(Equal("Selected User Groups Only"))
			})

			Context("when an unrecognized availability is provided", func() {
				BeforeEach(func() {
					availabilityVar = "bad value"
					availability = &availabilityVar
				})

				It("invokes the error handler", func() {
					err := client.Update(
						productSlug,
						releaseVersion,
						availability,
						releaseType,
					)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0).Error()).To(MatchRegexp(".*bad value.*"))
				})
			})
		})

		Context("when release type is provided", func() {
			var (
				releaseTypeVar string
			)

			BeforeEach(func() {
				releaseTypeVar = "all-in-one"
				releaseType = &releaseTypeVar
			})

			It("sets availability on release", func() {
				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.ReleaseType).To(Equal(pivnet.ReleaseType("All-In-One")))
			})

			It("correctly sets major", func() {
				releaseTypeVar = "major"
				releaseType = &releaseTypeVar

				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.ReleaseType).To(Equal(pivnet.ReleaseType("Major Release")))
			})

			It("correctly sets minor", func() {
				releaseTypeVar = "minor"
				releaseType = &releaseTypeVar

				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.ReleaseType).To(Equal(pivnet.ReleaseType("Minor Release")))
			})

			It("correctly sets service", func() {
				releaseTypeVar = "service"
				releaseType = &releaseTypeVar

				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.ReleaseType).To(Equal(pivnet.ReleaseType("Service Release")))
			})

			It("correctly sets maintenance", func() {
				releaseTypeVar = "maintenance"
				releaseType = &releaseTypeVar

				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.ReleaseType).To(Equal(pivnet.ReleaseType("Maintenance Release")))
			})

			It("correctly sets security", func() {
				releaseTypeVar = "security"
				releaseType = &releaseTypeVar

				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.ReleaseType).To(Equal(pivnet.ReleaseType("Security Release")))
			})

			It("correctly sets alpha", func() {
				releaseTypeVar = "alpha"
				releaseType = &releaseTypeVar

				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.ReleaseType).To(Equal(pivnet.ReleaseType("Alpha Release")))
			})

			It("correctly sets beta", func() {
				releaseTypeVar = "beta"
				releaseType = &releaseTypeVar

				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.ReleaseType).To(Equal(pivnet.ReleaseType("Beta Release")))
			})

			It("correctly sets edge", func() {
				releaseTypeVar = "edge"
				releaseType = &releaseTypeVar

				err := client.Update(
					productSlug,
					releaseVersion,
					availability,
					releaseType,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateReleaseCallCount()).To(Equal(1))
				_, providedRelease := fakePivnetClient.UpdateReleaseArgsForCall(0)
				Expect(providedRelease.ReleaseType).To(Equal(pivnet.ReleaseType("Edge Release")))
			})

			Context("when an unrecognized release type is provided", func() {
				BeforeEach(func() {
					releaseTypeVar = "bad value"
					releaseType = &releaseTypeVar
				})

				It("invokes the error handler", func() {
					err := client.Update(
						productSlug,
						releaseVersion,
						availability,
						releaseType,
					)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0).Error()).To(MatchRegexp(".*bad value.*"))
				})
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
