package releasedependency_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet/v7"
	"github.com/pivotal-cf/pivnet-cli/v3/commands/releasedependency"
	"github.com/pivotal-cf/pivnet-cli/v3/commands/releasedependency/releasedependencyfakes"
	"github.com/pivotal-cf/pivnet-cli/v3/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/v3/printer"
)

var _ = Describe("releasedependency commands", func() {
	var (
		fakePivnetClient *releasedependencyfakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		releasedependencies []pivnet.ReleaseDependency

		client *releasedependency.ReleaseDependencyClient
	)

	BeforeEach(func() {
		fakePivnetClient = &releasedependencyfakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		releasedependencies = []pivnet.ReleaseDependency{
			{
				Release: pivnet.DependentRelease{
					ID: 1234,
				},
			},
			{
				Release: pivnet.DependentRelease{
					ID: 2345,
				},
			},
		}

		fakePivnetClient.ReleaseDependenciesReturns(releasedependencies, nil)

		client = releasedependency.NewReleaseDependencyClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("ReleaseDependencies", func() {
		var (
			productSlug    string
			releaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			releaseVersion = "some release version"
		})

		It("lists all ReleaseDependencies", func() {
			err := client.List(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())

			var returnedReleaseDependencies []pivnet.ReleaseDependency
			err = json.Unmarshal(outBuffer.Bytes(), &returnedReleaseDependencies)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedReleaseDependencies).To(Equal(releasedependencies))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releasedependencies error")
				fakePivnetClient.ReleaseDependenciesReturns(nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.List(productSlug, releaseVersion)
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
	})

	Describe("AddReleaseDependency", func() {
		var (
			productSlug             string
			releaseVersion          string
			dependentProductSlug    string
			dependentReleaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			releaseVersion = "some release version"
			dependentProductSlug = "dependent product slug"
			dependentReleaseVersion = "dependent release version"
		})

		It("adds ReleaseDependency", func() {
			err := client.Add(
				productSlug,
				releaseVersion,
				dependentProductSlug,
				dependentReleaseVersion,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releaseDependencies error")
				fakePivnetClient.AddReleaseDependencyReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Add(
					productSlug,
					releaseVersion,
					dependentProductSlug,
					dependentReleaseVersion,
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
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Add(
					productSlug,
					releaseVersion,
					dependentProductSlug,
					dependentReleaseVersion,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error getting dependent release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleaseForVersionStub = func(productSlug string, releaseVersion string) (pivnet.Release, error) {
					if releaseVersion == dependentReleaseVersion {
						return pivnet.Release{}, expectedErr
					}
					return pivnet.Release{}, nil
				}
			})

			It("invokes the error handler", func() {
				err := client.Add(
					productSlug,
					releaseVersion,
					dependentProductSlug,
					dependentReleaseVersion,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("RemoveReleaseDependency", func() {
		var (
			productSlug             string
			releaseVersion          string
			dependentProductSlug    string
			dependentReleaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			releaseVersion = "some release version"
			dependentProductSlug = "dependent product slug"
			dependentReleaseVersion = "dependent release version"
		})

		It("removes ReleaseDependency", func() {
			err := client.Remove(
				productSlug,
				releaseVersion,
				dependentProductSlug,
				dependentReleaseVersion,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releaseDependencies error")
				fakePivnetClient.RemoveReleaseDependencyReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Remove(
					productSlug,
					releaseVersion,
					dependentProductSlug,
					dependentReleaseVersion,
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
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Remove(
					productSlug,
					releaseVersion,
					dependentProductSlug,
					dependentReleaseVersion,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error getting dependent release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleaseForVersionStub = func(productSlug string, releaseVersion string) (pivnet.Release, error) {
					if releaseVersion == dependentReleaseVersion {
						return pivnet.Release{}, expectedErr
					}
					return pivnet.Release{}, nil
				}
			})

			It("invokes the error handler", func() {
				err := client.Remove(
					productSlug,
					releaseVersion,
					dependentProductSlug,
					dependentReleaseVersion,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
