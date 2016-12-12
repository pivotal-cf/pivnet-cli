package productfile_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logger"
	"github.com/pivotal-cf/go-pivnet/logshim"
	"github.com/pivotal-cf/pivnet-cli/commands/productfile"
	"github.com/pivotal-cf/pivnet-cli/commands/productfile/productfilefakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("productfile commands", func() {
	var (
		l                logger.Logger
		fakeFilter       *productfilefakes.FakeFilter
		fakePivnetClient *productfilefakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer
		logBuffer bytes.Buffer

		productFiles []pivnet.ProductFile

		client *productfile.ProductFileClient
	)

	BeforeEach(func() {
		infoLogger := log.New(GinkgoWriter, "", 0)
		debugLogger := log.New(GinkgoWriter, "", 0)
		l = logshim.NewLogShim(infoLogger, debugLogger, true)

		fakeFilter = &productfilefakes.FakeFilter{}
		fakePivnetClient = &productfilefakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}
		logBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		productFiles = []pivnet.ProductFile{
			{
				ID:           1234,
				Name:         "Some File",
				AWSObjectKey: "/remote/path/some-file",
				Links: &pivnet.Links{
					Download: map[string]string{"href": "download-link-0"},
				},
			},
			{
				ID:           2345,
				Name:         "Some Other File",
				AWSObjectKey: "/remote/path/some-other-file",
				Links: &pivnet.Links{
					Download: map[string]string{"href": "download-link-1"},
				},
			},
		}

		client = productfile.NewProductFileClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			&logBuffer,
			printer.NewPrinter(&outBuffer),
			l,
			fakeFilter,
		)
	})

	Describe("List", func() {
		var (
			productSlug    string
			releaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = ""

			fakePivnetClient.ProductFilesReturns(productFiles, nil)
		})

		It("lists all ProductFiles", func() {
			err := client.List(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())

			var returnedProductFiles []pivnet.ProductFile
			err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFiles)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedProductFiles).To(Equal(productFiles))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productFiles error")
				fakePivnetClient.ProductFilesReturns(nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.List(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when release version is not empty", func() {
			BeforeEach(func() {
				releaseVersion = "some-release-version"
				fakePivnetClient.ProductFilesForReleaseReturns(productFiles, nil)
			})

			It("lists all ProductFiles", func() {
				err := client.List(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				var returnedProductFiles []pivnet.ProductFile
				err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFiles)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedProductFiles).To(Equal(productFiles))
			})

			Context("when there is an error getting release", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("releases error")
					fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.List(productSlug, releaseVersion)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})

			Context("when there is an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("productFiles error")
					fakePivnetClient.ProductFilesForReleaseReturns(nil, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.List(productSlug, releaseVersion)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})
		})
	})

	Describe("Create", func() {
		var (
			config pivnet.CreateProductFileConfig
		)

		BeforeEach(func() {
			config = pivnet.CreateProductFileConfig{
				Name: "some-name",
			}

			fakePivnetClient.CreateProductFileReturns(productFiles[0], nil)
		})

		It("creates ProductFile", func() {
			err := client.Create(config)
			Expect(err).NotTo(HaveOccurred())

			var returnedProductFile pivnet.ProductFile
			err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFile)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedProductFile).To(Equal(productFiles[0]))
			Expect(fakePivnetClient.CreateProductFileArgsForCall(0)).To(Equal(config))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfile error")
				fakePivnetClient.CreateProductFileReturns(pivnet.ProductFile{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Create(config)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Get", func() {
		var (
			productSlug    string
			releaseVersion string
			productFileID  int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = ""
			productFileID = productFiles[0].ID

			fakePivnetClient.ProductFileReturns(productFiles[0], nil)
		})

		It("gets ProductFile", func() {
			err := client.Get(productSlug, releaseVersion, productFileID)
			Expect(err).NotTo(HaveOccurred())

			var returnedProductFile pivnet.ProductFile
			err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFile)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedProductFile).To(Equal(productFiles[0]))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfile error")
				fakePivnetClient.ProductFileReturns(pivnet.ProductFile{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(productSlug, releaseVersion, productFileID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when release version is not empty", func() {
			BeforeEach(func() {
				releaseVersion = "some-release-version"
				fakePivnetClient.ProductFileForReleaseReturns(productFiles[0], nil)
			})

			It("gets ProductFile", func() {
				err := client.Get(productSlug, releaseVersion, productFileID)
				Expect(err).NotTo(HaveOccurred())

				var returnedProductFile pivnet.ProductFile
				err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFile)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedProductFile).To(Equal(productFiles[0]))
			})

			Context("when there is an error getting release", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("releases error")
					fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.Get(productSlug, releaseVersion, productFileID)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})

			Context("when there is an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("productFiles error")
					fakePivnetClient.ProductFileForReleaseReturns(pivnet.ProductFile{}, expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.Get(productSlug, releaseVersion, productFileID)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})
		})
	})

	Describe("Update", func() {
		var (
			productFileID int
			productSlug   string

			existingName        string
			existingFileType    string
			existingFileVersion string
			existingMD5         string
			existingDescription string

			name        string
			fileType    string
			fileVersion string
			md5         string
			description string

			existingProductFile pivnet.ProductFile
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			productFileID = productFiles[0].ID

			existingName = "some-name"
			existingFileType = "some-file-type"
			existingFileVersion = "some-file-type"
			existingMD5 = "some-md5"
			existingDescription = "some-description"

			name = "some-new-name"
			fileType = "some-new-file-type"
			fileVersion = "some-new-file-type"
			md5 = "some-new-md5"
			description = "some-new-description"

			existingProductFile = pivnet.ProductFile{
				ID:          productFileID,
				Name:        existingName,
				FileType:    existingFileType,
				FileVersion: existingFileVersion,
				MD5:         existingMD5,
				Description: existingDescription,
			}

			fakePivnetClient.ProductFileReturns(existingProductFile, nil)
			fakePivnetClient.UpdateProductFileReturns(productFiles[0], nil)
		})

		It("updates ProductFile", func() {
			err := client.Update(
				productFileID,
				productSlug,
				&name,
				&fileType,
				&fileVersion,
				&md5,
				&description,
			)
			Expect(err).NotTo(HaveOccurred())

			var returnedProductFile pivnet.ProductFile
			err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFile)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedProductFile).To(Equal(productFiles[0]))

			invokedProductSlug, invokedProductFile := fakePivnetClient.UpdateProductFileArgsForCall(0)
			Expect(invokedProductSlug).To(Equal(productSlug))
			Expect(invokedProductFile.ID).To(Equal(productFileID))
			Expect(invokedProductFile.Name).To(Equal(name))
			Expect(invokedProductFile.FileType).To(Equal(fileType))
			Expect(invokedProductFile.FileVersion).To(Equal(fileVersion))
			Expect(invokedProductFile.MD5).To(Equal(md5))
			Expect(invokedProductFile.Description).To(Equal(description))
		})

		Context("when optional fields are nil", func() {
			It("updates ProductFile with previous values", func() {
				err := client.Update(
					productFileID,
					productSlug,
					nil,
					nil,
					nil,
					nil,
					nil,
				)
				Expect(err).NotTo(HaveOccurred())

				var returnedProductFile pivnet.ProductFile
				err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFile)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedProductFile).To(Equal(productFiles[0]))

				invokedProductSlug, invokedProductFile := fakePivnetClient.UpdateProductFileArgsForCall(0)
				Expect(invokedProductSlug).To(Equal(productSlug))
				Expect(invokedProductFile.ID).To(Equal(productFileID))
				Expect(invokedProductFile.Name).To(Equal(existingName))
				Expect(invokedProductFile.FileType).To(Equal(existingFileType))
				Expect(invokedProductFile.FileVersion).To(Equal(existingFileVersion))
				Expect(invokedProductFile.MD5).To(Equal(existingMD5))
				Expect(invokedProductFile.Description).To(Equal(existingDescription))
			})
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfile error")
				fakePivnetClient.UpdateProductFileReturns(pivnet.ProductFile{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Update(
					productFileID,
					productSlug,
					&name,
					&fileType,
					&fileVersion,
					&md5,
					&description,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("AddToRelease", func() {
		var (
			productSlug    string
			releaseVersion string
			productFileID  int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			productFileID = productFiles[0].ID

			fakePivnetClient.AddProductFileToReleaseReturns(nil)
		})

		It("adds ProductFile", func() {
			err := client.AddToRelease(productSlug, releaseVersion, productFileID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error getting release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddToRelease(productSlug, releaseVersion, productFileID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfile error")
				fakePivnetClient.AddProductFileToReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddToRelease(productSlug, releaseVersion, productFileID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("RemoveFromRelease", func() {
		var (
			productSlug    string
			releaseVersion string
			productFileID  int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			productFileID = productFiles[0].ID

			fakePivnetClient.RemoveProductFileFromReleaseReturns(nil)
		})

		It("removes ProductFile", func() {
			err := client.RemoveFromRelease(productSlug, releaseVersion, productFileID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error getting release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveFromRelease(productSlug, releaseVersion, productFileID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfile error")
				fakePivnetClient.RemoveProductFileFromReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveFromRelease(productSlug, releaseVersion, productFileID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("AddToFileGroup", func() {
		var (
			productSlug   string
			fileGroupID   int
			productFileID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			fileGroupID = 5432
			productFileID = productFiles[0].ID

			fakePivnetClient.AddProductFileToFileGroupReturns(nil)
		})

		It("adds ProductFile", func() {
			err := client.AddToFileGroup(productSlug, fileGroupID, productFileID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfile error")
				fakePivnetClient.AddProductFileToFileGroupReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddToFileGroup(productSlug, fileGroupID, productFileID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Delete", func() {
		var (
			productSlug   string
			productFileID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			productFileID = productFiles[0].ID

			fakePivnetClient.DeleteProductFileReturns(productFiles[0], nil)
		})

		It("deletes ProductFile", func() {
			err := client.Delete(productSlug, productFileID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfile error")
				fakePivnetClient.DeleteProductFileReturns(pivnet.ProductFile{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Delete(productSlug, productFileID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Download", func() {
		const (
			fileContents = "file-contents"
		)

		var (
			productSlug    string
			releaseVersion string
			globs          []string
			productFileIDs []int
			downloadDir    string
			acceptEULA     bool

			tempDir   string
			filename  string
			releaseID int

			productFilesForReleaseErr error
			filterErr                 error
			downloadErr               error
		)

		BeforeEach(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			filename = "temp-file"

			productSlug = "some-product-slug"
			releaseVersion = "some-release-version"
			globs = []string{}
			productFileIDs = []int{productFiles[0].ID, productFiles[1].ID}
			downloadDir = tempDir
			acceptEULA = false

			returnedRelease := pivnet.Release{
				ID:      releaseID,
				Version: releaseVersion,
			}

			productFilesForReleaseErr = nil
			filterErr = nil
			downloadErr = nil

			fakePivnetClient.ReleaseForVersionReturns(returnedRelease, nil)
		})

		JustBeforeEach(func() {
			fakePivnetClient.ProductFilesForReleaseReturns(productFiles, productFilesForReleaseErr)
			fakeFilter.ProductFileKeysByGlobsReturns(productFiles, filterErr)

			fakePivnetClient.DownloadProductFileReturns(downloadErr)
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("downloads ProductFile", func() {
			err := client.Download(
				productSlug,
				releaseVersion,
				globs,
				productFileIDs,
				downloadDir,
				acceptEULA,
			)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePivnetClient.DownloadProductFileCallCount()).To(Equal(2))

			for i, pf := range productFiles {
				_, invokedProductSlug, invokedReleaseID, invokedProductFileID :=
					fakePivnetClient.DownloadProductFileArgsForCall(i)

				Expect(invokedProductSlug).To(Equal(productSlug))
				Expect(invokedReleaseID).To(Equal(releaseID))
				Expect(invokedProductFileID).To(Equal(pf.ID))
			}
		})

		Context("when globs are provided", func() {
			BeforeEach(func() {
				globs = []string{"glob1", "glob2"}
				productFileIDs = []int{}
			})

			It("downloads matching files", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					globs,
					productFileIDs,
					downloadDir,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.DownloadProductFileCallCount()).To(Equal(2))

				for i, pf := range productFiles {
					_, invokedProductSlug, invokedReleaseID, invokedProductFileID :=
						fakePivnetClient.DownloadProductFileArgsForCall(i)

					Expect(invokedProductSlug).To(Equal(productSlug))
					Expect(invokedReleaseID).To(Equal(releaseID))
					Expect(invokedProductFileID).To(Equal(pf.ID))
				}
			})

			Context("when filter returns an error", func() {
				BeforeEach(func() {
					filterErr = errors.New("filter error")
				})

				It("invokes the error handler", func() {
					err := client.Download(
						productSlug,
						releaseVersion,
						globs,
						productFileIDs,
						downloadDir,
						acceptEULA,
					)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(filterErr))
				})
			})
		})

		Context("when neither globs nor ids are provided", func() {
			BeforeEach(func() {
				globs = []string{}
				productFileIDs = []int{}
			})

			It("invokes the error handler", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					globs,
					productFileIDs,
					downloadDir,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
			})
		})

		Context("when both globs and ids are provided", func() {
			BeforeEach(func() {
				globs = []string{"glob1", "glob2"}
				productFileIDs = []int{1234, 2345}
			})

			It("invokes the error handler", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					globs,
					productFileIDs,
					downloadDir,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
			})
		})

		Context("when there is an error", func() {
			BeforeEach(func() {
				downloadErr = errors.New("download error")
			})

			It("invokes the error handler", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					globs,
					productFileIDs,
					downloadDir,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(downloadErr))
			})
		})

		Context("when there is an error getting all product files", func() {
			BeforeEach(func() {
				productFiles = nil
			})

			It("invokes the error handler", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					globs,
					productFileIDs,
					downloadDir,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0).Error()).To(ContainSubstring("No product files found"))
			})
		})

		Context("when no product files match filter", func() {
			BeforeEach(func() {
				productFilesForReleaseErr = errors.New("product files for release error")
			})

			It("invokes the error handler", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					globs,
					productFileIDs,
					downloadDir,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(productFilesForReleaseErr))
			})
		})

		Context("when there is an error getting release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleaseForVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					globs,
					productFileIDs,
					downloadDir,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error creating file", func() {
			BeforeEach(func() {
				downloadDir = "/not/a/valid/filepath"
			})

			It("invokes the error handler", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					globs,
					productFileIDs,
					downloadDir,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
			})
		})

		Context("when acceptEULA is true", func() {
			BeforeEach(func() {
				acceptEULA = true
			})

			It("accepts the EULA", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					globs,
					productFileIDs,
					downloadDir,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.AcceptEULACallCount()).To(Equal(1))
			})

			Context("when accepting the EULA returns an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = errors.New("product file error")
					fakePivnetClient.AcceptEULAReturns(expectedErr)
				})

				It("invokes the error handler", func() {
					err := client.Download(
						productSlug,
						releaseVersion,
						globs,
						productFileIDs,
						downloadDir,
						acceptEULA,
					)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
			})
		})
	})

	Describe("RemoveFromFileGroup", func() {
		var (
			productSlug   string
			fileGroupID   int
			productFileID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			fileGroupID = 1234
			productFileID = productFiles[0].ID

			fakePivnetClient.RemoveProductFileFromFileGroupReturns(nil)
		})

		It("removes ProductFile", func() {
			err := client.RemoveFromFileGroup(productSlug, fileGroupID, productFileID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfile error")
				fakePivnetClient.RemoveProductFileFromFileGroupReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveFromFileGroup(productSlug, fileGroupID, productFileID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
