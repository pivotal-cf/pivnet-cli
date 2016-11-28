package rc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRC(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RC Suite")
}
