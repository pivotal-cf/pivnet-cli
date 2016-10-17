package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

const (
	apiPrefix = "/api/v2"
	apiToken  = "some-api-token"
)

var _ = Describe("pivnet cli", func() {
	var (
		server *ghttp.Server
		host   string

		product pivnet.Product
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		host = server.URL()

		product = pivnet.Product{
			ID:   1234,
			Slug: "some-product-slug",
			Name: "some-product-name",
		}
	})

	runMainWithArgs := func(args ...string) *gexec.Session {
		args = append(
			args,
			fmt.Sprintf("--api-token=%s", apiToken),
			fmt.Sprintf("--host=%s", host),
		)

		_, err := fmt.Fprintf(GinkgoWriter, "Running command: %v\n", args)
		Expect(err).NotTo(HaveOccurred())

		command := exec.Command(pivnetBinPath, args...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		return session
	}

	Describe("Displaying help", func() {
		It("displays help to stdout with '-h'", func() {
			session := runMainWithArgs("-h")

			Eventually(session, executableTimeout).Should(gexec.Exit())
			Expect(session.Out).Should(gbytes.Say("Usage"))
		})

		It("displays help to stdout with '--help'", func() {
			session := runMainWithArgs("--help")

			Eventually(session, executableTimeout).Should(gexec.Exit())
			Expect(session.Out).Should(gbytes.Say("Usage"))
		})

		It("displays help to stdout with 'help'", func() {
			session := runMainWithArgs("help")

			Eventually(session, executableTimeout).Should(gexec.Exit())
			Expect(session.Out).Should(gbytes.Say("Usage:"))
			Expect(session.Out).Should(gbytes.Say("network.pivotal.io"))
		})

		It("displays help of a command to stdout with '--help'", func() {
			session := runMainWithArgs("product", "--help")

			Eventually(session, executableTimeout).Should(gexec.Exit())
			Expect(session.Out).Should(gbytes.Say("Usage"))
			Expect(session.Out).Should(gbytes.Say("-product-slug"))
		})

		It("displays help of a command to stdout with '-h'", func() {
			session := runMainWithArgs("product", "-h")

			Eventually(session, executableTimeout).Should(gexec.Exit())
			Expect(session.Out).Should(gbytes.Say("Usage"))
			Expect(session.Out).Should(gbytes.Say("-product-slug"))
		})
	})

	Describe("Displaying version", func() {
		It("displays version with '-v'", func() {
			session := runMainWithArgs("-v")

			Eventually(session, executableTimeout).Should(gexec.Exit(0))
			Expect(session).Should(gbytes.Say("dev"))
		})

		It("displays version with '--version'", func() {
			session := runMainWithArgs("--version")

			Eventually(session, executableTimeout).Should(gexec.Exit(0))
			Expect(session).Should(gbytes.Say("dev"))
		})
	})

	Describe("printing as json", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						fmt.Sprintf("%s/products/%s", apiPrefix, product.Slug),
					),
					ghttp.RespondWithJSONEncoded(http.StatusOK, product),
				),
			)
		})

		It("prints as json", func() {
			session := runMainWithArgs(
				"--format=json",
				"product",
				"--product-slug", product.Slug)

			Eventually(session, executableTimeout).Should(gexec.Exit(0))

			var receivedProduct pivnet.Product
			err := json.Unmarshal(session.Out.Contents(), &receivedProduct)
			Expect(err).NotTo(HaveOccurred())

			Expect(receivedProduct.Slug).To(Equal(product.Slug))
		})
	})

	Describe("printing as yaml", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						fmt.Sprintf("%s/products/%s", apiPrefix, product.Slug),
					),
					ghttp.RespondWithJSONEncoded(http.StatusOK, product),
				),
			)
		})

		It("prints as yaml", func() {
			session := runMainWithArgs(
				"--format=yaml",
				"product",
				"--product-slug", product.Slug)

			Eventually(session, executableTimeout).Should(gexec.Exit(0))

			var receivedProduct pivnet.Product
			err := yaml.Unmarshal(session.Out.Contents(), &receivedProduct)
			Expect(err).NotTo(HaveOccurred())

			Expect(receivedProduct.Slug).To(Equal(product.Slug))
		})
	})
})
