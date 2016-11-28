package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

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

		tempDir        string
		configFilepath string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		host = server.URL()

		product = pivnet.Product{
			ID:   1234,
			Slug: "some-product-slug",
			Name: "some-product-name",
		}

		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		configFilepath = filepath.Join(tempDir, ".pivnetrc")
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDir)
		Expect(err).NotTo(HaveOccurred())
	})

	runMainWithArgs := func(args ...string) *gexec.Session {
		args = append(
			args,
			fmt.Sprintf("--config=%s", configFilepath),
			fmt.Sprintf("--host=%s", host),
			"--verbose",
		)

		_, err := fmt.Fprintf(GinkgoWriter, "Running command: %v\n", args)
		Expect(err).NotTo(HaveOccurred())

		command := exec.Command(pivnetBinPath, args...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		return session
	}

	login := func() {
		session := runMainWithArgs(
			"login",
			"--api-token", apiToken,
		)
		Eventually(session, executableTimeout).Should(gexec.Exit(0))
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

	It("exits with error when not logged in", func() {
		session := runMainWithArgs(
			"products",
		)

		Eventually(session, executableTimeout).Should(gexec.Exit(1))
		Eventually(session).Should(gbytes.Say("login"))
	})

	Context("when logged in", func() {
		BeforeEach(func() {
			// we hit the login endpoint twice:
			// once for the login command and once for the actual command
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						fmt.Sprintf("%s/authentication", apiPrefix),
					),
					ghttp.RespondWith(http.StatusOK, ""),
				),
			)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						fmt.Sprintf("%s/authentication", apiPrefix),
					),
					ghttp.RespondWith(http.StatusOK, ""),
				),
			)
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
				login()

				session := runMainWithArgs(
					"--format=json",
					"product",
					"--product-slug", product.Slug,
				)

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
				login()

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
})
