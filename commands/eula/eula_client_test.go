package eula_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet/v6"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/eula"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/eula/eulafakes"
	"github.com/pivotal-cf/pivnet-cli/v2/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/v2/printer"
)

var _ = Describe("eula commands", func() {
	var (
		fakePivnetClient *eulafakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		eulas []pivnet.EULA

		client *eula.EULAClient
	)

	BeforeEach(func() {
		fakePivnetClient = &eulafakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		eulas = []pivnet.EULA{
			{
				ID:   1234,
				Name: "some eula",
				Slug: "some-eula",
			},
			{
				ID:   2345,
				Name: "another eula",
				Slug: "another-eula",
			},
		}

		fakePivnetClient.EULAsReturns(eulas, nil)
		fakePivnetClient.EULAReturns(eulas[0], nil)
		fakePivnetClient.AcceptEULAReturns(nil)

		client = eula.NewEULAClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("EULAs", func() {
		It("lists all EULAs", func() {
			err := client.List()
			Expect(err).NotTo(HaveOccurred())

			var returnedEULAs []pivnet.EULA
			err = json.Unmarshal(outBuffer.Bytes(), &returnedEULAs)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedEULAs).To(Equal(eulas))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("eulas error")
				fakePivnetClient.EULAsReturns(nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.List()
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("EULACommand", func() {
		It("gets EULA", func() {
			err := client.Get(eulas[0].Slug)
			Expect(err).NotTo(HaveOccurred())

			var returnedEULA pivnet.EULA
			err = json.Unmarshal(outBuffer.Bytes(), &returnedEULA)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedEULA).To(Equal(eulas[0]))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("eulas error")
				fakePivnetClient.EULAReturns(pivnet.EULA{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(eulas[0].Slug)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("AcceptEULACommand", func() {
		const (
			productSlug = "some-product-slug"
		)

		var (
			release pivnet.Release
		)

		BeforeEach(func() {
			release = pivnet.Release{
				ID:          1234,
				Version:     "version 0.2.3",
				Description: "Some release with some description.",
			}
		})

		It("accepts EULA", func() {
			err := client.AcceptEULA(productSlug, release.Version)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("eulas error")
				fakePivnetClient.AcceptEULAReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AcceptEULA(productSlug, release.Version)
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
				err := client.AcceptEULA(productSlug, release.Version)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
