package filter_test

import (
	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logger/loggerfakes"
	"github.com/pivotal-cf/pivnet-cli/filter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	var (
		fakeLogger *loggerfakes.FakeLogger
		f          *filter.Filter

		releases []pivnet.Release
	)

	BeforeEach(func() {
		fakeLogger = &loggerfakes.FakeLogger{}

		f = filter.NewFilter(fakeLogger)

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

	Describe("ProductFileNamesByGlobs", func() {
		var (
			productFiles []pivnet.ProductFile
			globs        []string
		)

		BeforeEach(func() {
			productFiles = []pivnet.ProductFile{
				{
					ID:   1234,
					Name: "name0",
				},
				{
					ID:   2345,
					Name: "name1",
				},
				{
					ID:   3456,
					Name: "name2",
				},
			}

			globs = []string{"*name1*", "*name2*"}
		})

		It("returns the download links that match the glob filters", func() {
			filtered, err := f.ProductFileNamesByGlobs(
				productFiles,
				globs,
			)

			Expect(err).NotTo(HaveOccurred())
			Expect(filtered).To(HaveLen(2))
			Expect(filtered).To(Equal([]pivnet.ProductFile{productFiles[1], productFiles[2]}))
		})

		Context("when a bad pattern is passed", func() {
			BeforeEach(func() {
				globs = []string{"["}
			})

			It("returns an error", func() {
				_, err := f.ProductFileNamesByGlobs(
					productFiles,
					globs,
				)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("syntax error in pattern"))
			})
		})

		Describe("Passed a glob that matches no files", func() {
			BeforeEach(func() {
				globs = []string{"*will-not-match*"}
			})

			It("returns empty slice", func() {
				filtered, err := f.ProductFileNamesByGlobs(
					productFiles,
					globs,
				)
				Expect(err).To(HaveOccurred())

				Expect(filtered).To(HaveLen(0))
			})
		})

		Describe("When a glob that matches a file and glob that does not match a file", func() {
			BeforeEach(func() {
				globs = []string{"name1", "does-not-exist.txt"}
			})

			It("returns an error", func() {
				_, err := f.ProductFileNamesByGlobs(
					productFiles,
					globs,
				)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("no product files match glob: does-not-exist.txt"))
			})
		})
	})
})
