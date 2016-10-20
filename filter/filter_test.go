package filter_test

import (
	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-cli/filter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	var (
		f *filter.Filter

		releases []pivnet.Release
	)

	BeforeEach(func() {
		f = filter.NewFilter()

		releases = []pivnet.Release{
			{
				ID:          1,
				Version:     "version1",
				ReleaseType: pivnet.ReleaseType("foo"),
			},
			{
				ID:          2,
				Version:     "version2",
				ReleaseType: pivnet.ReleaseType("bar"),
			},
			{
				ID:          3,
				Version:     "version3",
				ReleaseType: pivnet.ReleaseType("foo"),
			},
		}
	})

	Describe("ReleasesByVersion", func() {
		var (
			version  string
			releases []pivnet.Release
		)

		BeforeEach(func() {
			version = "version2"

			releases = []pivnet.Release{
				{
					ID:          1,
					Version:     "version1",
					ReleaseType: "foo",
				},
				{
					ID:          2,
					Version:     "version2",
					ReleaseType: "bar",
				},
				{
					ID:          3,
					Version:     "version3",
					ReleaseType: "foo",
				},
				{
					ID:          4,
					Version:     "version3.2",
					ReleaseType: "foo-minor",
				},
				{
					ID:          5,
					Version:     "version3.1.2",
					ReleaseType: "foo-patch",
				},
			}
		})

		It("filters releases by version without error", func() {
			filteredReleases, err := f.ReleasesByVersion(releases, version)

			Expect(err).NotTo(HaveOccurred())

			Expect(filteredReleases).To(HaveLen(1))
			Expect(filteredReleases).To(ContainElement(releases[1]))
		})

		Context("when the input releases are nil", func() {
			BeforeEach(func() {
				releases = nil
			})

			It("returns empty slice without error", func() {
				filteredReleases, err := f.ReleasesByVersion(releases, version)

				Expect(err).NotTo(HaveOccurred())

				Expect(filteredReleases).NotTo(BeNil())
				Expect(filteredReleases).To(HaveLen(0))
			})
		})

		Describe("Matching on regex", func() {
			Context("when the regex matches one release versions", func() {
				BeforeEach(func() {
					version = `version3\.1\..*`
				})

				It("returns all releases that match the regex without error", func() {
					filteredReleases, err := f.ReleasesByVersion(releases, version)

					Expect(err).NotTo(HaveOccurred())

					Expect(filteredReleases).To(HaveLen(1))
					Expect(filteredReleases).To(ContainElement(releases[4]))
				})
			})

			Context("when the regex matches multiple release versions", func() {
				BeforeEach(func() {
					version = `version3\..*`
				})

				It("returns all releases that match the regex without error", func() {
					filteredReleases, err := f.ReleasesByVersion(releases, version)

					Expect(err).NotTo(HaveOccurred())

					Expect(filteredReleases).To(HaveLen(2))
					Expect(filteredReleases).To(ContainElement(releases[3]))
					Expect(filteredReleases).To(ContainElement(releases[4]))
				})
			})

			Context("when the regex is invalid", func() {
				BeforeEach(func() {
					version = "some(invalid^regex"
				})

				It("returns error", func() {
					_, err := f.ReleasesByVersion(releases, version)

					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
