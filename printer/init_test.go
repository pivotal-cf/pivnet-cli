package printer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const ()

func TestPrinter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Printer Suite")
}
