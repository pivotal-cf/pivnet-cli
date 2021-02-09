package dependencyspecifier_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet/v7"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/dependencyspecifier"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/dependencyspecifier/dependencyspecifierfakes"
	"github.com/pivotal-cf/pivnet-cli/v2/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/v2/printer"
)

var _ = Describe("dependencyspecifier commands", func() {
	var (
		fakePivnetClient *dependencyspecifierfakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		dependencySpecifiers []pivnet.DependencySpecifier

		client *dependencyspecifier.DependencySpecifierClient
	)

	BeforeEach(func() {
		fakePivnetClient = &dependencyspecifierfakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		dependencySpecifiers = []pivnet.DependencySpecifier{
			{
				ID: 1234,
			},
			{
				ID: 2345,
			},
		}

		fakePivnetClient.DependencySpecifiersReturns(dependencySpecifiers, nil)

		client = dependencyspecifier.NewDependencySpecifierClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("DependencySpecifiers", func() {
		var (
			productSlug    string
			releaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			releaseVersion = "some release version"
		})

		It("lists all DependencySpecifiers", func() {
			err := client.List(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())

			var returnedDependencySpecifiers []pivnet.DependencySpecifier
			err = json.Unmarshal(outBuffer.Bytes(), &returnedDependencySpecifiers)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedDependencySpecifiers).To(Equal(dependencySpecifiers))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("dependencySpecifiers error")
				fakePivnetClient.DependencySpecifiersReturns(nil, expectedErr)
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

	Describe("CreateDependencySpecifier", func() {
		var (
			productSlug          string
			releaseVersion       string
			dependentProductSlug string
			specifier            string
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			releaseVersion = "some release version"
			dependentProductSlug = "dependent product slug"
			specifier = "specifier"
		})

		It("creates DependencySpecifier", func() {
			err := client.Create(
				productSlug,
				releaseVersion,
				dependentProductSlug,
				specifier,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releaseDependencies error")
				fakePivnetClient.CreateDependencySpecifierReturns(pivnet.DependencySpecifier{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Create(
					productSlug,
					releaseVersion,
					dependentProductSlug,
					specifier,
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
				err := client.Create(
					productSlug,
					releaseVersion,
					dependentProductSlug,
					specifier,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("DeleteDependencySpecifier", func() {
		var (
			productSlug           string
			releaseVersion        string
			dependencySpecifierID int
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			releaseVersion = "some release version"
			dependencySpecifierID = 1234
		})

		It("deletes DependencySpecifier", func() {
			err := client.Delete(
				productSlug,
				releaseVersion,
				dependencySpecifierID,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releaseDependencies error")
				fakePivnetClient.DeleteDependencySpecifierReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Delete(
					productSlug,
					releaseVersion,
					dependencySpecifierID,
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
				err := client.Delete(
					productSlug,
					releaseVersion,
					dependencySpecifierID,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
