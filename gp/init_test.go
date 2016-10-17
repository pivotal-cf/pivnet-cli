package gp_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const ()

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GP Suite")
}
