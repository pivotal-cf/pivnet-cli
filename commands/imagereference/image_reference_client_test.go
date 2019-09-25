package imagereference_test

import (
	"bytes"
	"encoding/json"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/go-pivnet/v2/logger"
	"github.com/pivotal-cf/go-pivnet/v2/logshim"
	"github.com/pivotal-cf/pivnet-cli/commands/imagereference"
	"github.com/pivotal-cf/pivnet-cli/commands/imagereference/imagereferencefakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
	"log"
)

var _ = Describe("imagereference commands", func() {
	var (
		l                    logger.Logger
		fakePivnetClient     *imagereferencefakes.FakePivnetClient

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
})
