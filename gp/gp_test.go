package gp_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-cli/gp"
	"github.com/pivotal-cf/go-pivnet/logger/loggerfakes"
)

const (
	apiPrefix = "/api/v2"
)

var _ = Describe("Client", func() {
	var (
		server *ghttp.Server

		client     *gp.Client
		fakeLogger *loggerfakes.FakeLogger
	)

	BeforeEach(func() {
		server = ghttp.NewServer()

		fakeLogger = &loggerfakes.FakeLogger{}

		config := pivnet.ClientConfig{
			Host:      server.URL(),
			Token:     "some-token",
			UserAgent: "some-user-agent",
		}
		client = gp.NewClient(config, fakeLogger)
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
})
