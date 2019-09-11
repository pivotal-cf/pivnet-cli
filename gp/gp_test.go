package gp_test

import (
	"fmt"
	"github.com/pivotal-cf/pivnet-cli/gp/gpfakes"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/go-pivnet/v2/logger/loggerfakes"
	"github.com/pivotal-cf/pivnet-cli/gp"
)

const (
	apiPrefix = "/api/v2"
)

var _ = Describe("Client", func() {
	var (
		server *ghttp.Server

		client                 *gp.Client
		fakeLogger             *loggerfakes.FakeLogger
		fakeAccessTokenService *gpfakes.FakeAccessTokenService
	)

	BeforeEach(func() {
		server = ghttp.NewServer()

		fakeLogger = &loggerfakes.FakeLogger{}
		fakeAccessTokenService = &gpfakes.FakeAccessTokenService{}

		config := pivnet.ClientConfig{
			Host:      server.URL(),
			UserAgent: "some-user-agent",
		}
		client = gp.NewClient(fakeAccessTokenService, config, fakeLogger)
	})

	Describe("ReleaseForVersion", func() {
		var (
			productSlug    string
			releaseVersion string

			release pivnet.Release

			releasesResponseStatusCode int
			releasesResponse           pivnet.ReleasesResponse

			releaseResponseStatusCode int
			releaseResponse           pivnet.Release
		)

		BeforeEach(func() {
			productSlug = "product-slug"
			releaseVersion = "release-version"

			release = pivnet.Release{
				ID:      1234,
				Version: releaseVersion,
			}

			releasesResponseStatusCode = http.StatusOK
			releasesResponse = pivnet.ReleasesResponse{[]pivnet.Release{release}}

			releaseResponseStatusCode = http.StatusOK
			releaseResponse = release
		})

		JustBeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						fmt.Sprintf(
							"%s/products/%s/releases",
							apiPrefix,
							productSlug,
						),
					),
					ghttp.RespondWithJSONEncoded(releasesResponseStatusCode, releasesResponse),
				),
			)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						fmt.Sprintf(
							"%s/products/%s/releases/%d",
							apiPrefix,
							productSlug,
							release.ID,
						),
					),
					ghttp.RespondWithJSONEncoded(releaseResponseStatusCode, releaseResponse),
				),
			)
		})

		It("returns release", func() {
			returnedRelease, err := client.ReleaseForVersion(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedRelease.ID).To(Equal(release.ID))
			Expect(returnedRelease.Version).To(Equal(releaseVersion))
		})

		Context("When getting releases returns an error", func() {
			BeforeEach(func() {
				releasesResponseStatusCode = http.StatusTeapot
			})

			It("returns an error", func() {
				_, err := client.ReleaseForVersion(productSlug, releaseVersion)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("When release is not found", func() {
			BeforeEach(func() {
				releasesResponse = pivnet.ReleasesResponse{[]pivnet.Release{{
					ID:      2345,
					Version: "some-other-version",
				}}}
			})

			It("returns an error", func() {
				_, err := client.ReleaseForVersion(productSlug, releaseVersion)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("When getting release returns an error", func() {
			BeforeEach(func() {
				releaseResponseStatusCode = http.StatusTeapot
			})

			It("returns an error", func() {
				_, err := client.ReleaseForVersion(productSlug, releaseVersion)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ProductFilesForRelease", func() {
		var (
			productSlug string
			releaseID   int

			productFile pivnet.ProductFile

			productFilesResponseStatusCode int
			productFilesResponse           pivnet.ProductFilesResponse

			fileGroupProductFile pivnet.ProductFile
			fileGroup            pivnet.FileGroup

			fileGroupsResponseStatusCode int
			fileGroupsResponse           pivnet.FileGroupsResponse
		)

		BeforeEach(func() {
			productSlug = "product-slug"
			releaseID = 1234

			productFile = pivnet.ProductFile{
				ID: 5678,
			}

			productFilesResponseStatusCode = http.StatusOK
			productFilesResponse = pivnet.ProductFilesResponse{[]pivnet.ProductFile{productFile}}

			fileGroupProductFile = pivnet.ProductFile{
				ID: 80,
			}

			fileGroup = pivnet.FileGroup{
				ID:           8,
				ProductFiles: []pivnet.ProductFile{fileGroupProductFile},
			}

			fileGroupsResponseStatusCode = http.StatusOK
			fileGroupsResponse = pivnet.FileGroupsResponse{[]pivnet.FileGroup{fileGroup}}
		})

		JustBeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						fmt.Sprintf(
							"%s/products/%s/releases/%d/product_files",
							apiPrefix,
							productSlug,
							releaseID,
						),
					),
					ghttp.RespondWithJSONEncoded(productFilesResponseStatusCode, productFilesResponse),
				),
			)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						fmt.Sprintf(
							"%s/products/%s/releases/%d/file_groups",
							apiPrefix,
							productSlug,
							releaseID,
						),
					),
					ghttp.RespondWithJSONEncoded(fileGroupsResponseStatusCode, fileGroupsResponse),
				),
			)
		})

		It("return a list of product files from release and file groups of that release", func() {
			returnedProductFiles, err := client.ProductFilesForRelease(productSlug, releaseID)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(returnedProductFiles)).To(Equal(2))
			Expect(returnedProductFiles[0].ID).To(Equal(productFile.ID))
			Expect(returnedProductFiles[1].ID).To(Equal(fileGroupProductFile.ID))
		})
	})
})
