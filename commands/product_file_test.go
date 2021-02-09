package commands_test

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet/v7"
	"github.com/pivotal-cf/pivnet-cli/v2/commands"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/productfile"
)

var _ = Describe("product file commands", func() {
	var (
		field reflect.StructField

		fakeProductFileClient *commandsfakes.FakeProductFileClient
	)

	BeforeEach(func() {
		fakeProductFileClient = &commandsfakes.FakeProductFileClient{}

		commands.NewProductFileClient = func(productfile.PivnetClient) commands.ProductFileClient {
			return fakeProductFileClient
		}
	})

	Describe("ProductFilesCommand", func() {
		var (
			cmd commands.ProductFilesCommand
		)

		BeforeEach(func() {
			cmd = commands.ProductFilesCommand{}
		})

		It("invokes the ProductFile client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeProductFileClient.ListCallCount()).To(Equal(1))
		})

		Context("when the ProductFile client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeProductFileClient.ListReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Context("when Init returns an error", func() {
			BeforeEach(func() {
				initErr = fmt.Errorf("init error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(initErr))
			})
		})

		Context("when Authentication returns an error", func() {
			BeforeEach(func() {
				authErr = fmt.Errorf("auth error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(authErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ProductFilesCommand{}, "ProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("p"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-slug"))
			})
		})

		Describe("ReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ProductFilesCommand{}, "ReleaseVersion")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("r"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-version"))
			})
		})
	})

	Describe("ProductFileCommand", func() {
		var (
			cmd commands.ProductFileCommand
		)

		BeforeEach(func() {
			cmd = commands.ProductFileCommand{}
		})

		It("invokes the ProductFile client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeProductFileClient.GetCallCount()).To(Equal(1))
		})

		Context("when the ProductFile client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeProductFileClient.GetReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Context("when Init returns an error", func() {
			BeforeEach(func() {
				initErr = fmt.Errorf("init error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(initErr))
			})
		})

		Context("when Authentication returns an error", func() {
			BeforeEach(func() {
				authErr = fmt.Errorf("auth error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(authErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ProductFileCommand{}, "ProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("p"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-slug"))
			})
		})

		Describe("ReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ProductFileCommand{}, "ReleaseVersion")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("r"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-version"))
			})
		})

		Describe("ProductFileID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ProductFileCommand{}, "ProductFileID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-file-id"))
			})
		})
	})

	Describe("CreateProductFileCommand", func() {
		var (
			productSlug  string
			name         string
			awsObjectKey string
			fileType     string
			fileVersion  string
			sha256       string
			md5          string

			description        string
			docsURL            string
			includedFiles      []string
			platforms          []string
			releasedAt         string
			systemRequirements []string

			cmd commands.CreateProductFileCommand
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			name = "some product file"
			awsObjectKey = "some aws object key"
			fileType = "some file type"
			fileVersion = "some file version"
			sha256 = "some sha256"
			md5 = "some md5"

			description = "some description"
			docsURL = "some-docs-url"
			includedFiles = []string{"file1", "file2"}
			platforms = []string{"platform1", "platform2"}
			releasedAt = "released-at"
			systemRequirements = []string{"system1", "system2"}

			cmd = commands.CreateProductFileCommand{
				ProductSlug:        productSlug,
				Name:               name,
				AWSObjectKey:       awsObjectKey,
				FileType:           fileType,
				FileVersion:        fileVersion,
				SHA256:             sha256,
				MD5:                md5,
				Description:        description,
				DocsURL:            docsURL,
				IncludedFiles:      includedFiles,
				Platforms:          platforms,
				ReleasedAt:         releasedAt,
				SystemRequirements: systemRequirements,
			}
		})

		It("invokes the ProductFile client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			config := pivnet.CreateProductFileConfig{
				ProductSlug:        productSlug,
				Name:               name,
				AWSObjectKey:       awsObjectKey,
				FileType:           fileType,
				FileVersion:        fileVersion,
				SHA256:             sha256,
				MD5:                md5,
				Description:        description,
				DocsURL:            docsURL,
				IncludedFiles:      includedFiles,
				Platforms:          platforms,
				ReleasedAt:         releasedAt,
				SystemRequirements: systemRequirements,
			}

			Expect(fakeProductFileClient.CreateCallCount()).To(Equal(1))
			Expect(fakeProductFileClient.CreateArgsForCall(0)).To(Equal(config))
		})

		Context("when the ProductFile client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeProductFileClient.CreateReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Context("when Init returns an error", func() {
			BeforeEach(func() {
				initErr = fmt.Errorf("init error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(initErr))
			})
		})

		Context("when Authentication returns an error", func() {
			BeforeEach(func() {
				authErr = fmt.Errorf("auth error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(authErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("p"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-slug"))
			})
		})

		Describe("Name flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Name")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("name"))
			})
		})

		Describe("AWSObjectKey flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "AWSObjectKey")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("aws-object-key"))
			})
		})

		Describe("FileType flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "FileType")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("file-type"))
			})
		})

		Describe("FileVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "FileVersion")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("file-version"))
			})
		})

		Describe("SHA256 flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "SHA256")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("sha256"))
			})
		})

		Describe("MD5 flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "MD5")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("md5"))
			})
		})

		Describe("Description flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Description")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("description"))
			})
		})

		Describe("DocsURL flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "DocsURL")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("docs-url"))
			})
		})

		Describe("IncludedFiles flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "IncludedFiles")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("included-file"))
			})
		})

		Describe("Platforms flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Platforms")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("platform"))
			})
		})

		Describe("ReleasedAt flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ReleasedAt")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("released-at"))
			})
		})

		Describe("SystemRequirements flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "SystemRequirements")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("system-requirement"))
			})
		})

	})

	Describe("UpdateProductFileCommand", func() {
		var (
			productSlug   string
			productFileID int

			description        string
			fileVersion        string
			sha256             string
			md5                string
			name               string
			docsURL            string
			systemRequirements []string

			cmd commands.UpdateProductFileCommand
		)

		BeforeEach(func() {
			productSlug = "some product slug"
			productFileID = 1234

			description = "some description"
			fileVersion = "some file version"
			sha256 = "some sha256"
			md5 = "some md5"
			name = "some product file"
			docsURL = "some-docs-url"
			systemRequirements = []string{"system1", "system2"}

			cmd = commands.UpdateProductFileCommand{
				ProductSlug:        productSlug,
				ProductFileID:      productFileID,
				Name:               &name,
				Description:        &description,
				FileVersion:        &fileVersion,
				SHA256:             &sha256,
				MD5:                &md5,
				DocsURL:            &docsURL,
				SystemRequirements: &systemRequirements,
			}
		})

		It("invokes the ProductFile client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeProductFileClient.UpdateCallCount()).To(Equal(1))

			invokedProductFileID,
			invokedProductSlug,
			invokedName,
			invokedFileVersion,
			invokedSHA256,
			invokedMD5,
			invokedDescription,
			invokedDocsURL,
			invokedSystemRequirements := fakeProductFileClient.UpdateArgsForCall(0)

			Expect(invokedProductFileID).To(Equal(productFileID))
			Expect(invokedProductSlug).To(Equal(productSlug))
			Expect(*invokedName).To(Equal(name))
			Expect(*invokedFileVersion).To(Equal(fileVersion))
			Expect(*invokedSHA256).To(Equal(sha256))
			Expect(*invokedMD5).To(Equal(md5))
			Expect(*invokedDescription).To(Equal(description))
			Expect(*invokedDocsURL).To(Equal(docsURL))
			Expect(*invokedSystemRequirements).To(Equal(systemRequirements))
		})

		Context("when the ProductFile client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeProductFileClient.UpdateReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Context("when Init returns an error", func() {
			BeforeEach(func() {
				initErr = fmt.Errorf("init error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(initErr))
			})
		})

		Context("when Authentication returns an error", func() {
			BeforeEach(func() {
				authErr = fmt.Errorf("auth error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(authErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("p"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-slug"))
			})
		})

		Describe("ProductFileID flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "ProductFileID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-file-id"))
			})
		})

		Describe("Name flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Name")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("name"))
			})
		})

		Describe("Description flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "Description")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("description"))
			})
		})

		Describe("FileVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "FileVersion")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("file-version"))
			})
		})

		Describe("SHA256 flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "SHA256")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("sha256"))
			})
		})

		Describe("MD5 flag", func() {
			BeforeEach(func() {
				field = fieldFor(cmd, "MD5")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("md5"))
			})
		})
	})

	Describe("AddProductFileCommand", func() {
		var (
			cmd commands.AddProductFileCommand

			releaseVersion string
			fileGroupID    int
		)

		BeforeEach(func() {
			releaseVersion = "some release version"
			fileGroupID = 5432

			cmd = commands.AddProductFileCommand{ReleaseVersion: &releaseVersion}
		})

		Context("when neither releaseVersion nor fileGroupID are provided", func() {
			BeforeEach(func() {
				cmd = commands.AddProductFileCommand{}
			})

			It("returns an error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("one of release-version or file-group-id"))
			})
		})

		Context("when both releaseVersion and fileGroupID are provided", func() {
			BeforeEach(func() {
				cmd = commands.AddProductFileCommand{
					ReleaseVersion: &releaseVersion,
					FileGroupID:    &fileGroupID,
				}
			})

			It("returns an error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("only one of release-version or file-group-id"))
			})
		})

		Context("when release-version is provided", func() {
			BeforeEach(func() {
				cmd = commands.AddProductFileCommand{ReleaseVersion: &releaseVersion}
			})

			It("invokes the ProductFile client", func() {
				err := cmd.Execute(nil)

				Expect(err).NotTo(HaveOccurred())

				Expect(fakeProductFileClient.AddToReleaseCallCount()).To(Equal(1))
			})

			Context("when the ProductFile client returns an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("expected error")
					fakeProductFileClient.AddToReleaseReturns(expectedErr)
				})

				It("forwards the error", func() {
					err := cmd.Execute(nil)

					Expect(err).To(Equal(expectedErr))
				})
			})
		})

		Context("when file-group-id is provided", func() {
			BeforeEach(func() {
				cmd = commands.AddProductFileCommand{FileGroupID: &fileGroupID}
			})

			It("invokes the ProductFile client", func() {
				err := cmd.Execute(nil)

				Expect(err).NotTo(HaveOccurred())

				Expect(fakeProductFileClient.AddToFileGroupCallCount()).To(Equal(1))
			})

			Context("when the ProductFile client returns an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("expected error")
					fakeProductFileClient.AddToFileGroupReturns(expectedErr)
				})

				It("forwards the error", func() {
					err := cmd.Execute(nil)

					Expect(err).To(Equal(expectedErr))
				})
			})
		})

		Context("when Init returns an error", func() {
			BeforeEach(func() {
				initErr = fmt.Errorf("init error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(initErr))
			})
		})

		Context("when Authentication returns an error", func() {
			BeforeEach(func() {
				authErr = fmt.Errorf("auth error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(authErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddProductFileCommand{}, "ProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("p"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-slug"))
			})
		})

		Describe("ProductFileID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddProductFileCommand{}, "ProductFileID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-file-id"))
			})
		})

		Describe("ReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddProductFileCommand{}, "ReleaseVersion")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("r"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-version"))
			})
		})

		Describe("FileGroupID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.AddProductFileCommand{}, "FileGroupID")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("f"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("file-group-id"))
			})
		})
	})

	Describe("RemoveProductFileCommand", func() {
		var (
			cmd commands.RemoveProductFileCommand

			releaseVersion string
			fileGroupID    int
		)

		BeforeEach(func() {
			releaseVersion = "some release version"
			fileGroupID = 5432

			cmd = commands.RemoveProductFileCommand{ReleaseVersion: &releaseVersion}
		})

		Context("when neither releaseVersion nor fileGroupID are provided", func() {
			BeforeEach(func() {
				cmd = commands.RemoveProductFileCommand{}
			})

			It("returns an error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("one of release-version or file-group-id"))
			})
		})

		Context("when both releaseVersion and fileGroupID are provided", func() {
			BeforeEach(func() {
				cmd = commands.RemoveProductFileCommand{
					ReleaseVersion: &releaseVersion,
					FileGroupID:    &fileGroupID,
				}
			})

			It("returns an error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("only one of release-version or file-group-id"))
			})
		})

		Context("when release-version is provided", func() {
			BeforeEach(func() {
				cmd = commands.RemoveProductFileCommand{ReleaseVersion: &releaseVersion}
			})

			It("invokes the ProductFile client", func() {
				err := cmd.Execute(nil)

				Expect(err).NotTo(HaveOccurred())

				Expect(fakeProductFileClient.RemoveFromReleaseCallCount()).To(Equal(1))
			})

			Context("when the ProductFile client returns an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("expected error")
					fakeProductFileClient.RemoveFromReleaseReturns(expectedErr)
				})

				It("forwards the error", func() {
					err := cmd.Execute(nil)

					Expect(err).To(Equal(expectedErr))
				})
			})
		})

		Context("when file-group-id is provided", func() {
			BeforeEach(func() {
				cmd = commands.RemoveProductFileCommand{FileGroupID: &fileGroupID}
			})

			It("invokes the ProductFile client", func() {
				err := cmd.Execute(nil)

				Expect(err).NotTo(HaveOccurred())

				Expect(fakeProductFileClient.RemoveFromFileGroupCallCount()).To(Equal(1))
			})

			Context("when the ProductFile client returns an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("expected error")
					fakeProductFileClient.RemoveFromFileGroupReturns(expectedErr)
				})

				It("forwards the error", func() {
					err := cmd.Execute(nil)

					Expect(err).To(Equal(expectedErr))
				})
			})
		})

		Context("when Init returns an error", func() {
			BeforeEach(func() {
				initErr = fmt.Errorf("init error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(initErr))
			})
		})

		Context("when Authentication returns an error", func() {
			BeforeEach(func() {
				authErr = fmt.Errorf("auth error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(authErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveProductFileCommand{}, "ProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("p"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-slug"))
			})
		})

		Describe("ProductFileID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveProductFileCommand{}, "ProductFileID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-file-id"))
			})
		})

		Describe("ReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveProductFileCommand{}, "ReleaseVersion")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("r"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-version"))
			})
		})

		Describe("FileGroupID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.RemoveProductFileCommand{}, "FileGroupID")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("file-group-id"))
			})
		})
	})

	Describe("DeleteProductFileCommand", func() {
		var (
			cmd commands.DeleteProductFileCommand
		)

		BeforeEach(func() {
			cmd = commands.DeleteProductFileCommand{}
		})

		It("invokes the ProductFile client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeProductFileClient.DeleteCallCount()).To(Equal(1))
		})

		Context("when the ProductFile client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeProductFileClient.DeleteReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Context("when Init returns an error", func() {
			BeforeEach(func() {
				initErr = fmt.Errorf("init error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(initErr))
			})
		})

		Context("when Authentication returns an error", func() {
			BeforeEach(func() {
				authErr = fmt.Errorf("auth error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(authErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DeleteProductFileCommand{}, "ProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("p"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-slug"))
			})
		})

		Describe("ProductFileID flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DeleteProductFileCommand{}, "ProductFileID")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-file-id"))
			})
		})
	})

	Describe("DownloadProductFileCommand", func() {
		var (
			cmd commands.DownloadProductFilesCommand
		)

		BeforeEach(func() {
			cmd = commands.DownloadProductFilesCommand{}
		})

		It("invokes the ProductFile client", func() {
			err := cmd.Execute(nil)

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeProductFileClient.DownloadCallCount()).To(Equal(1))
		})

		Context("when the ProductFile client returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("expected error")
				fakeProductFileClient.DownloadReturns(expectedErr)
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(expectedErr))
			})
		})

		Context("when Init returns an error", func() {
			BeforeEach(func() {
				initErr = fmt.Errorf("init error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(initErr))
			})
		})

		Context("when Authentication returns an error", func() {
			BeforeEach(func() {
				authErr = fmt.Errorf("auth error")
			})

			It("forwards the error", func() {
				err := cmd.Execute(nil)

				Expect(err).To(Equal(authErr))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DownloadProductFilesCommand{}, "ProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("p"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-slug"))
			})
		})

		Describe("ReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DownloadProductFilesCommand{}, "ReleaseVersion")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-version"))
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("r"))
			})
		})

		Describe("ProductFileIDs flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DownloadProductFilesCommand{}, "ProductFileIDs")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("i"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-file-id"))
			})
		})

		Describe("Globs flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DownloadProductFilesCommand{}, "Globs")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("g"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("glob"))
			})
		})

		Describe("DownloadDir flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DownloadProductFilesCommand{}, "DownloadDir")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("d"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("download-dir"))
			})
		})

		Describe("AcceptEULA flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.DownloadProductFilesCommand{}, "AcceptEULA")
			})

			It("is not required", func() {
				Expect(isRequired(field)).To(BeFalse())
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("accept-eula"))
			})
		})
	})
})
