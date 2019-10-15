package errorhandler_test

import (
	"bytes"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("ErrorHandler", func() {
	var (
		errorHandler errorhandler.ErrorHandler

		format    string
		outWriter *bytes.Buffer
		logWriter *bytes.Buffer

		inputErr error
	)

	BeforeEach(func() {
		outWriter = &bytes.Buffer{}
		logWriter = &bytes.Buffer{}

		format = printer.PrintAsTable

		inputErr = fmt.Errorf("some error")
	})

	JustBeforeEach(func() {
		errorHandler = errorhandler.NewErrorHandler(
			format,
			outWriter,
			logWriter,
		)
	})

	It("returns ErrAlreadyHandled", func() {
		err := errorHandler.HandleError(inputErr)

		Expect(err).To(Equal(errorhandler.ErrAlreadyHandled))
	})

	It("writes to outWriter", func() {
		_ = errorHandler.HandleError(inputErr)

		Expect(outWriter.String()).To(ContainSubstring(fmt.Sprint("some error")))
	})

	Context("when the error is nil", func() {
		BeforeEach(func() {
			inputErr = nil
		})

		It("returns nil", func() {
			err := errorHandler.HandleError(inputErr)

			Expect(err).NotTo(HaveOccurred())
		})

		It("does not write to printer", func() {
			_ = errorHandler.HandleError(nil)

			Expect(outWriter.String()).To(BeEmpty())
			Expect(logWriter.String()).To(BeEmpty())
		})
	})

	Describe("print as JSON", func() {
		BeforeEach(func() {
			format = printer.PrintAsJSON
		})

		It("writes to logWriter", func() {
			_ = errorHandler.HandleError(inputErr)

			Expect(logWriter.String()).To(ContainSubstring(fmt.Sprint("some error")))
		})
	})

	Describe("print as YAML", func() {
		BeforeEach(func() {
			format = printer.PrintAsYAML
		})

		It("writes to logWriter", func() {
			_ = errorHandler.HandleError(inputErr)

			Expect(logWriter.String()).To(ContainSubstring(fmt.Sprint("some error")))
		})
	})

	Describe("Handling specific Pivnet errors", func() {
		Describe("pivnet.ErrUnauthorized", func() {
			BeforeEach(func() {
				inputErr = pivnet.ErrUnauthorized{}
			})

			It("returns custom message", func() {
				_ = errorHandler.HandleError(inputErr)

				Expect(outWriter.String()).To(ContainSubstring(fmt.Sprint("Failed to authenticate - please provide valid API token")))
			})
		})

		Describe("pivnet.ErrNotFound", func() {
			BeforeEach(func() {
				inputErr = pivnet.ErrNotFound{
					ResponseCode: 404,
					Message:      "something not found",
				}
			})

			It("returns custom message", func() {
				_ = errorHandler.HandleError(inputErr)

				Expect(outWriter.String()).To(ContainSubstring(fmt.Sprint("Pivnet error: something not found")))
			})
		})

		Describe("pivnet.ErrPivnetOther", func() {
			BeforeEach(func() {
				inputErr = pivnet.ErrPivnetOther{
					ResponseCode: http.StatusTeapot,
					Message:      "something is wrong",
					Errors:       []string{"err1", "err2"},
				}
			})

			It("returns custom message", func() {
				_ = errorHandler.HandleError(inputErr)

				Expect(outWriter.String()).To(ContainSubstring("418 - something is wrong"))
				Expect(outWriter.String()).To(ContainSubstring("err1"))
				Expect(outWriter.String()).To(ContainSubstring("err2"))
			})
		})
	})
})
