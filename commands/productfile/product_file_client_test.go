package productfile_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logger/loggerfakes"
	"github.com/pivotal-cf/pivnet-cli/commands/productfile"
	"github.com/pivotal-cf/pivnet-cli/commands/productfile/productfilefakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("productfile commands", func() {
	var (
		fakeLogger       *loggerfakes.FakeLogger
		fakePivnetClient *productfilefakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer
		logBuffer bytes.Buffer

		productfiles []pivnet.ProductFile

		client *productfile.ProductFileClient
	)

	BeforeEach(func() {
		fakeLogger = &loggerfakes.FakeLogger{}
		fakePivnetClient = &productfilefakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}
		logBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		productfiles = []pivnet.ProductFile{
			{
				ID: 1234,
			},
			{
				ID: 2345,
			},
		}

		client = productfile.NewProductFileClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			&logBuffer,
			printer.NewPrinter(&outBuffer),
			fakeLogger,
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

			fakePivnetClient.GetProductFilesReturns(productfiles, nil)
		})

		It("lists all ProductFiles", func() {
			err := client.List(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())

			var returnedProductFiles []pivnet.ProductFile
			err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFiles)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedProductFiles).To(Equal(productfiles))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfiles error")
				fakePivnetClient.GetProductFilesReturns(nil, expectedErr)
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
				fakePivnetClient.GetProductFilesForReleaseReturns(productfiles, nil)
			})

			It("lists all ProductFiles", func() {
				err := client.List(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				var returnedProductFiles []pivnet.ProductFile
				err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFiles)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedProductFiles).To(Equal(productfiles))
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
					expectedErr = errors.New("productfiles error")
					fakePivnetClient.GetProductFilesForReleaseReturns(nil, expectedErr)
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

			fakePivnetClient.CreateProductFileReturns(productfiles[0], nil)
		})

		It("creates ProductFile", func() {
			err := client.Create(config)
			Expect(err).NotTo(HaveOccurred())

			var returnedProductFile pivnet.ProductFile
			err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFile)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedProductFile).To(Equal(productfiles[0]))
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
			productFileID = productfiles[0].ID

			fakePivnetClient.GetProductFileReturns(productfiles[0], nil)
		})

		It("gets ProductFile", func() {
			err := client.Get(productSlug, releaseVersion, productFileID)
			Expect(err).NotTo(HaveOccurred())

			var returnedProductFile pivnet.ProductFile
			err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFile)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedProductFile).To(Equal(productfiles[0]))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfile error")
				fakePivnetClient.GetProductFileReturns(pivnet.ProductFile{}, expectedErr)
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
				fakePivnetClient.GetProductFileForReleaseReturns(productfiles[0], nil)
			})

			It("gets ProductFile", func() {
				err := client.Get(productSlug, releaseVersion, productFileID)
				Expect(err).NotTo(HaveOccurred())

				var returnedProductFile pivnet.ProductFile
				err = json.Unmarshal(outBuffer.Bytes(), &returnedProductFile)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedProductFile).To(Equal(productfiles[0]))
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
					expectedErr = errors.New("productfiles error")
					fakePivnetClient.GetProductFileForReleaseReturns(pivnet.ProductFile{}, expectedErr)
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
			productFileID = productfiles[0].ID

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

			fakePivnetClient.GetProductFileReturns(existingProductFile, nil)
			fakePivnetClient.UpdateProductFileReturns(productfiles[0], nil)
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

			Expect(returnedProductFile).To(Equal(productfiles[0]))

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

				Expect(returnedProductFile).To(Equal(productfiles[0]))

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
			productFileID = productfiles[0].ID

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
			productFileID = productfiles[0].ID

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
			productFileID = productfiles[0].ID

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
			productFileID = productfiles[0].ID

			fakePivnetClient.DeleteProductFileReturns(productfiles[0], nil)
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
			productSlug      string
			releaseVersion   string
			productFileID    int
			providedFilepath string
			acceptEULA       bool

			tempDir   string
			filename  string
			releaseID int
		)

		BeforeEach(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			filename = "temp-file"

			productSlug = "some-product-slug"
			releaseVersion = "some-release-version"
			productFileID = productfiles[0].ID
			providedFilepath = filepath.Join(tempDir, filename)
			acceptEULA = false

			returnedRelease := pivnet.Release{
				ID:      releaseID,
				Version: releaseVersion,
			}

			fakePivnetClient.ReleaseForVersionReturns(returnedRelease, nil)
			fakePivnetClient.DownloadFileStub = func(writer io.Writer, downloadLink string) (err error, retryable bool) {
				_, err = fmt.Fprintf(writer, fileContents)
				return err, false
			}
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("downloads ProductFile", func() {
			err := client.Download(
				productSlug,
				releaseVersion,
				productFileID,
				providedFilepath,
				acceptEULA,
			)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePivnetClient.DownloadFileCallCount()).To(Equal(1))
			_, invokedLink := fakePivnetClient.DownloadFileArgsForCall(0)

			expectedLink := fmt.Sprintf(
				"/products/%s/releases/%d/product_files/%d/download",
				productSlug,
				releaseID,
				productFileID,
			)
			Expect(invokedLink).To(Equal(expectedLink))

			contents, err := ioutil.ReadFile(providedFilepath)
			Expect(err).NotTo(HaveOccurred())
			Expect(contents).To(Equal([]byte(fileContents)))
		})

		Context("when there is a non-retryable error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfile error")
				fakePivnetClient.DownloadFileReturns(expectedErr, false)
			})

			It("invokes the error handler", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					productFileID,
					providedFilepath,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is a retryable error", func() {
			var (
				expectedErr error
				triesCount  int
			)

			BeforeEach(func() {
				expectedErr = errors.New("productfile error")
				triesCount = 0
				fakePivnetClient.DownloadFileStub = func(writer io.Writer, downloadLink string) (err error, retryable bool) {
					failed := (triesCount < 1)
					triesCount++

					if failed {
						return expectedErr, true
					} else {
						return nil, true
					}
				}
			})

			It("retries the download", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					productFileID,
					providedFilepath,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.DownloadFileCallCount()).To(Equal(2))
				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(0))
			})

			Context("when the error has occurred three times", func() {
				BeforeEach(func() {
					fakePivnetClient.DownloadFileReturns(expectedErr, true)
				})

				It("invokes the error handler", func() {
					err := client.Download(
						productSlug,
						releaseVersion,
						productFileID,
						providedFilepath,
						acceptEULA,
					)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakePivnetClient.DownloadFileCallCount()).To(Equal(3))
					Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
					Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
				})
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
					productFileID,
					providedFilepath,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error creating file", func() {
			BeforeEach(func() {
				providedFilepath = "/not/a/valid/filepath"
			})

			It("invokes the error handler", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					productFileID,
					providedFilepath,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
			})
		})

		Context("when there is an error getting product file", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("product file error")
				fakePivnetClient.GetProductFileForReleaseReturns(pivnet.ProductFile{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Download(
					productSlug,
					releaseVersion,
					productFileID,
					providedFilepath,
					acceptEULA,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
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
					productFileID,
					providedFilepath,
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
						productFileID,
						providedFilepath,
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
			productFileID = productfiles[0].ID

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
