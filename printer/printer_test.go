package printer_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("Printer", func() {
	var (
		p printer.Printer

		outputWriter *bytes.Buffer
	)

	BeforeEach(func() {
		outputWriter = &bytes.Buffer{}

		p = printer.NewPrinter(outputWriter)
	})

	Describe("Println", func() {
		It("Prints a line", func() {
			err := p.Println("some message")

			Expect(err).NotTo(HaveOccurred())

			Expect(outputWriter.String()).To(Equal("some message\n"))
		})

		Context("when writing fails", func() {
			BeforeEach(func() {
				writer := errWriter{}
				p = printer.NewPrinter(writer)
			})

			It("returns an error", func() {
				err := p.Println("")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("PrintJSON", func() {
		It("Prints object as JSON", func() {
			object := map[string]interface{}{
				"bar": 1234,
				"foo": "foo val",
			}
			err := p.PrintJSON(object)

			Expect(err).NotTo(HaveOccurred())

			expectedString := `{
"foo": "foo val",
"bar": 1234
}
`

			Expect(outputWriter.String()).To(MatchJSON(expectedString))
		})

		Context("when marshalling the object fails", func() {
			It("returns an error", func() {
				object := map[string]interface{}{
					"foo": make(chan string),
				}
				err := p.PrintJSON(object)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when writing fails", func() {
			BeforeEach(func() {
				writer := errWriter{}
				p = printer.NewPrinter(writer)
			})

			It("returns an error", func() {
				err := p.PrintJSON(struct{}{})

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("PrintYAML", func() {
		It("Prints object as YAML", func() {
			object := map[string]interface{}{
				"bar": 1234,
				"foo": "foo val",
			}
			err := p.PrintYAML(object)

			Expect(err).NotTo(HaveOccurred())

			expectedString := `---
foo: "foo val"
bar: 1234
`

			Expect(outputWriter.String()).To(MatchYAML(expectedString))
		})

		Context("when marshalling the object fails", func() {
			It("returns an error", func() {
				object := map[string]interface{}{
					"foo": make(chan string),
				}
				err := p.PrintYAML(object)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when writing fails", func() {
			BeforeEach(func() {
				writer := errWriter{}
				p = printer.NewPrinter(writer)
			})

			It("returns an error", func() {
				err := p.PrintYAML(struct{}{})

				Expect(err).To(HaveOccurred())
			})
		})
	})
})

type errWriter struct {
}

func (w errWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("Error writer erroring out")
}
