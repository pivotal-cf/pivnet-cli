package semver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf/pivnet-cli/semver"
)

var _ = Describe("semantic versioning comparison", func() {
	DescribeTable("", func(s1, s2 string, expected int){
		Expect(semver.Compare(s1, s2)).To(Equal(expected))
	},
		Entry("when equal", "1.2.3", "1.2.3", 0),
		Entry("when greater", "1.2.4", "1.2.3", 1),
		Entry("when smaller", "1.2.2", "1.2.3", -1),
		Entry("when both empty", "", "", 0),
		Entry("when left side empty", "", "1.2.3", -1),
		Entry("when right side empty", "1.2.3", "", 1),
		Entry("when both zero", "0", "0", 0),
		Entry("when left side zero", "0", "1.2.3", -1),
		Entry("when right side zero", "1.2.3", "0", 1),
		Entry("when letters on both sides", "v1.2.3", "v1.2.3", 0),
		Entry("when letter on right side", "v1.2.3", "1.2.3", 0),
		Entry("when letter on left side", "1.2.3", "v1.2.3", 0),
		Entry("when consecutive dots on the both sides", "1..2", "1..2", 0),
		Entry("when consecutive dots on the left side", "1..2", "1.2", 0),
		Entry("when consecutive dots on the right side", "1.2", "1..2", 0),
	)
})
