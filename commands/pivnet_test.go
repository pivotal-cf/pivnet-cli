package commands_test

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
	"github.com/pivotal-cf/pivnet-cli/rc"
)

var _ = Describe("Pivnet commands", func() {
	var (
		field reflect.StructField
	)

	Describe("Init", func() {
		var (
			server *ghttp.Server

			outBuffer bytes.Buffer

			fakeRCHandler    *commandsfakes.FakeRCHandler
			fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

			profile         *rc.PivnetProfile
			profileErr      error
			profileRequired bool

			apiToken string
		)

		BeforeEach(func() {
			fakeRCHandler = &commandsfakes.FakeRCHandler{}
			fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

			server = ghttp.NewServer()

			commands.ErrorHandler = fakeErrorHandler
			commands.RC = fakeRCHandler

			commands.Pivnet.Verbose = true

			outBuffer = bytes.Buffer{}
			commands.LogWriter = &outBuffer
			commands.OutputWriter = &outBuffer
			commands.Printer = printer.NewPrinter(commands.OutputWriter)

			products := []pivnet.Product{
				{
					ID:   2345,
					Slug: "another-product-slug",
					Name: "another-product-name",
				},
			}

			productsResponse := pivnet.ProductsResponse{
				Products: products,
			}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("%s/authentication", apiPrefix)),
					ghttp.RespondWithJSONEncoded(http.StatusOK, productsResponse),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("%s/products", apiPrefix)),
					ghttp.RespondWithJSONEncoded(http.StatusOK, productsResponse),
				),
			)

			profileRequired = true

			apiToken = "some-api-token"
			profile = &rc.PivnetProfile{
				Name:     "some-profile",
				APIToken: apiToken,
				Host:     server.URL(),
			}
			profileErr = nil
		})

		JustBeforeEach(func() {
			fakeRCHandler.ProfileForNameReturns(profile, profileErr)

			commands.Init = origInitFunc
		})

		It("redacts api token from profile", func() {
			err := commands.Init(profileRequired)
			Expect(err).NotTo(HaveOccurred())

			_, err = fmt.Fprintf(commands.OutputWriter, apiToken)
			Expect(err).NotTo(HaveOccurred())

			Expect(outBuffer.String()).Should(ContainSubstring("*** redacted api token ***"))
			Expect(outBuffer.String()).ShouldNot(ContainSubstring(apiToken))
		})

		Context("when profile validation returns an error", func() {
			BeforeEach(func() {
				profile.APIToken = ""
			})

			It("Invokes the error handler", func() {
				err := commands.Init(profileRequired)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0).Error()).To(ContainSubstring("login"))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0).Error()).To(ContainSubstring("API"))
			})
		})

		Context("when getting profile returns an error", func() {
			BeforeEach(func() {
				profileErr = fmt.Errorf("some profile error")
			})

			It("Invokes the error handler", func() {
				err := commands.Init(profileRequired)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(profileErr))
			})
		})

		Context("when profile is nil", func() {
			BeforeEach(func() {
				profile = nil
			})

			It("Invokes the error handler", func() {
				err := commands.Init(profileRequired)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0).Error()).To(ContainSubstring("login"))
			})

			Context("when profile is not required", func() {
				BeforeEach(func() {
					profileRequired = false
				})

				It("returns without error", func() {
					err := commands.Init(profileRequired)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(0))
				})
			})
		})

		AfterEach(func() {
			server.Close()

			commands.Pivnet.Verbose = false
		})
	})

	Describe("VersionFunc", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "VersionFunc")
		})

		It("contains short flag", func() {
			Expect(shortTag(field)).To(Equal("v"))
		})

		It("contains long flag", func() {
			Expect(longTag(field)).To(Equal("version"))
		})
	})

	Describe("Version command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "Version")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("version"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("v"))
		})
	})

	Describe("Help command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "Help")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("help"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("h"))
		})
	})

	Describe("Verbose flag", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "Verbose")
		})

		It("contains long flag", func() {
			Expect(longTag(field)).To(Equal("verbose"))
		})
	})

	Describe("Format flag", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "Format")
		})

		It("contains long flag", func() {
			Expect(longTag(field)).To(Equal("format"))
		})

		It("defaults to table", func() {
			Expect(field.Tag.Get("default")).To(Equal("table"))
		})

		It("contains choice", func() {
			Expect(string(field.Tag)).To(
				MatchRegexp(`choice:"table".*choice:"json".*choice:"yaml"`))
		})

		It("is not required", func() {
			Expect(isRequired(field)).To(BeFalse())
		})
	})

	Describe("Login command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "Login")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("login"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("l"))
		})
	})

	Describe("Logout command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "Logout")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("logout"))
		})
	})

	Describe("ReleaseTypes command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "ReleaseTypes")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("release-types"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("rts"))
		})
	})

	Describe("EULAs command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "EULAs")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("eulas"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("es"))
		})
	})

	Describe("EULA command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "EULA")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("eula"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("e"))
		})
	})

	Describe("AcceptEULA command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "AcceptEULA")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("accept-eula"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("ae"))
		})
	})

	Describe("Products command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "Products")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("products"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("ps"))
		})
	})

	Describe("Product command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "Product")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("product"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("p"))
		})
	})

	Describe("ProductFiles command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "ProductFiles")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("product-files"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("pfs"))
		})
	})

	Describe("ProductFile command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "ProductFile")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("product-file"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("pf"))
		})
	})

	Describe("CreateProductFile command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "CreateProductFile")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("create-product-file"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("cpf"))
		})
	})

	Describe("UpdateProductFile command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "UpdateProductFile")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("update-product-file"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("upf"))
		})
	})

	Describe("AddProductFile command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "AddProductFile")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("add-product-file"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("apf"))
		})
	})

	Describe("RemoveProductFile command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "RemoveProductFile")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("remove-product-file"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("rpf"))
		})
	})

	Describe("DeleteProductFile command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "DeleteProductFile")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("delete-product-file"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("dpf"))
		})
	})

	Describe("FileGroups command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "FileGroups")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("file-groups"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("fgs"))
		})
	})

	Describe("FileGroup command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "FileGroup")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("file-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("fg"))
		})
	})

	Describe("CreateFileGroup command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "CreateFileGroup")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("create-file-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("cfg"))
		})
	})

	Describe("UpdateFileGroup command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "UpdateFileGroup")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("update-file-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("ufg"))
		})
	})

	Describe("DeleteFileGroup command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "DeleteFileGroup")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("delete-file-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("dfg"))
		})
	})

	Describe("AddFileGroupToRelease command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "AddFileGroupToRelease")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("add-file-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("afg"))
		})
	})

	Describe("RemoveFileGroupFromRelease command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "RemoveFileGroupFromRelease")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("remove-file-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("rfg"))
		})
	})

	Describe("Releases command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "Releases")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("releases"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("rs"))
		})
	})

	Describe("Release command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "Release")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("release"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("r"))
		})
	})

	Describe("CreateRelease command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "CreateRelease")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("create-release"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("cr"))
		})
	})

	Describe("DeleteRelease command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "DeleteRelease")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("delete-release"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("dr"))
		})
	})

	Describe("UserGroups command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "UserGroups")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("user-groups"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("ugs"))
		})
	})

	Describe("UserGroup command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "UserGroup")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("user-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("ug"))
		})
	})

	Describe("AddUserGroup command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "AddUserGroup")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("add-user-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("aug"))
		})
	})

	Describe("RemoveUserGroup command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "RemoveUserGroup")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("remove-user-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("rug"))
		})
	})

	Describe("CreateUserGroup command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "CreateUserGroup")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("create-user-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("cug"))
		})
	})

	Describe("UpdateUserGroup command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "UpdateUserGroup")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("update-user-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("uug"))
		})
	})

	Describe("DeleteUserGroup command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "DeleteUserGroup")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("delete-user-group"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("dug"))
		})
	})

	Describe("AddUserGroupMemberCommand command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "AddUserGroupMember")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("add-user-group-member"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("augm"))
		})
	})

	Describe("RemoveUserGroupMemberCommand command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "RemoveUserGroupMember")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("remove-user-group-member"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("rugm"))
		})
	})

	Describe("ReleaseDependencies command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "ReleaseDependencies")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("release-dependencies"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("rds"))
		})
	})

	Describe("AddReleaseDependency command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "AddReleaseDependency")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("add-release-dependency"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("ard"))
		})
	})

	Describe("RemoveReleaseDependency command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "RemoveReleaseDependency")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("remove-release-dependency"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("rrd"))
		})
	})

	Describe("ReleaseUpgradePaths command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "ReleaseUpgradePaths")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("release-upgrade-paths"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("rups"))
		})
	})

	Describe("AddReleaseUpgradePath command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "AddReleaseUpgradePath")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("add-release-upgrade-path"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("arup"))
		})
	})

	Describe("RemoveReleaseUpgradePath command", func() {
		BeforeEach(func() {
			field = fieldFor(commands.Pivnet, "RemoveReleaseUpgradePath")
		})

		It("contains command", func() {
			Expect(command(field)).To(Equal("remove-release-upgrade-path"))
		})

		It("contains alias", func() {
			Expect(alias(field)).To(Equal("rrup"))
		})
	})
})
