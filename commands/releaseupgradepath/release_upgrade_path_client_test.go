package releaseupgradepath_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-cli/commands/releaseupgradepath"
	"github.com/pivotal-cf/pivnet-cli/commands/releaseupgradepath/releaseupgradepathfakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("releaseupgradepath commands", func() {
	var (
		fakePivnetClient *releaseupgradepathfakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		releaseUpgradePaths []pivnet.ReleaseUpgradePath

		client *releaseupgradepath.ReleaseUpgradePathClient
	)

	BeforeEach(func() {
		fakePivnetClient = &releaseupgradepathfakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		releaseUpgradePaths = []pivnet.ReleaseUpgradePath{
			{
				Release: pivnet.UpgradePathRelease{
					ID: 1234,
				},
			},
			{
				Release: pivnet.UpgradePathRelease{
					ID: 2345,
				},
			},
		}

		fakePivnetClient.ReleaseUpgradePathsReturns(releaseUpgradePaths, nil)

		client = releaseupgradepath.NewReleaseUpgradePathClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("ReleaseUpgradePaths", func() {
		var (
			productSlug    string
			releaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			releaseVersion = "some release version"
		})

		It("lists all ReleaseUpgradePaths", func() {
			err := client.List(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())

			var returnedReleaseUpgradePaths []pivnet.ReleaseUpgradePath
			err = json.Unmarshal(outBuffer.Bytes(), &returnedReleaseUpgradePaths)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedReleaseUpgradePaths).To(Equal(releaseUpgradePaths))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releaseUpgradePaths error")
				fakePivnetClient.ReleaseUpgradePathsReturns(nil, expectedErr)
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

	Describe("AddReleaseUpgradePath", func() {
		var (
			productSlug            string
			releaseVersion         string
			previousReleaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			releaseVersion = "some release version"
			previousReleaseVersion = "previous release version"
		})

		It("adds ReleaseUpgradePath", func() {
			err := client.Add(productSlug, releaseVersion, previousReleaseVersion)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releaseUpgradePaths error")
				fakePivnetClient.AddReleaseUpgradePathReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Add(productSlug, releaseVersion, previousReleaseVersion)
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
				err := client.Add(productSlug, releaseVersion, previousReleaseVersion)
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
				fakePivnetClient.ReleaseForVersionStub = func(productSlug string, releaseVersion string) (pivnet.Release, error) {
					if releaseVersion == previousReleaseVersion {
						return pivnet.Release{}, expectedErr
					}
					return pivnet.Release{}, nil
				}
			})

			It("invokes the error handler", func() {
				err := client.Add(productSlug, releaseVersion, previousReleaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("RemoveReleaseUpgradePath", func() {
		var (
			productSlug            string
			releaseVersion         string
			previousReleaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			releaseVersion = "some release version"
			previousReleaseVersion = "previous release version"
		})

		It("adds ReleaseUpgradePath", func() {
			err := client.Remove(productSlug, releaseVersion, previousReleaseVersion)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releaseUpgradePaths error")
				fakePivnetClient.RemoveReleaseUpgradePathReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Remove(productSlug, releaseVersion, previousReleaseVersion)
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
				err := client.Remove(productSlug, releaseVersion, previousReleaseVersion)
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
				fakePivnetClient.ReleaseForVersionStub = func(productSlug string, releaseVersion string) (pivnet.Release, error) {
					if releaseVersion == previousReleaseVersion {
						return pivnet.Release{}, expectedErr
					}
					return pivnet.Release{}, nil
				}
			})

			It("invokes the error handler", func() {
				err := client.Remove(productSlug, releaseVersion, previousReleaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
