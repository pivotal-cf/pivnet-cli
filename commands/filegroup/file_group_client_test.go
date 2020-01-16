package filegroup_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet/v4"
	"github.com/pivotal-cf/pivnet-cli/commands/filegroup"
	"github.com/pivotal-cf/pivnet-cli/commands/filegroup/filegroupfakes"
	"github.com/pivotal-cf/pivnet-cli/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

var _ = Describe("filegroup commands", func() {
	var (
		fakePivnetClient *filegroupfakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		fileGroups []pivnet.FileGroup

		client *filegroup.FileGroupClient
	)

	BeforeEach(func() {
		fakePivnetClient = &filegroupfakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		fileGroups = []pivnet.FileGroup{
			{
				ID: 1234,
			},
			{
				ID: 2345,
			},
		}

		client = filegroup.NewFileGroupClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
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

			fakePivnetClient.FileGroupsReturns(fileGroups, nil)
		})

		It("lists all FileGroups", func() {
			err := client.List(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())

			var returnedFileGroups []pivnet.FileGroup
			err = json.Unmarshal(outBuffer.Bytes(), &returnedFileGroups)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedFileGroups).To(Equal(fileGroups))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("fileGroups error")
				fakePivnetClient.FileGroupsReturns(nil, expectedErr)
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
				fakePivnetClient.FileGroupsForReleaseReturns(fileGroups, nil)
			})

			It("lists all FileGroups", func() {
				err := client.List(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				var returnedFileGroups []pivnet.FileGroup
				err = json.Unmarshal(outBuffer.Bytes(), &returnedFileGroups)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedFileGroups).To(Equal(fileGroups))
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
					expectedErr = errors.New("fileGroups error")
					fakePivnetClient.FileGroupsForReleaseReturns(nil, expectedErr)
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

	Describe("Get", func() {
		var (
			productSlug string
			fileGroupID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			fileGroupID = fileGroups[0].ID

			fakePivnetClient.FileGroupReturns(fileGroups[0], nil)
		})

		It("gets FileGroup", func() {
			err := client.Get(productSlug, fileGroupID)
			Expect(err).NotTo(HaveOccurred())

			var returnedFileGroup pivnet.FileGroup
			err = json.Unmarshal(outBuffer.Bytes(), &returnedFileGroup)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedFileGroup).To(Equal(fileGroups[0]))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("filegroup error")
				fakePivnetClient.FileGroupReturns(pivnet.FileGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(productSlug, fileGroupID)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Create", func() {
		var (
			productSlug string
			name        string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			name = "new-name"

			fakePivnetClient.CreateFileGroupReturns(fileGroups[0], nil)
		})

		It("creates FileGroup", func() {
			err := client.Create(productSlug, name)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("filegroup error")
				fakePivnetClient.CreateFileGroupReturns(pivnet.FileGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Create(productSlug, name)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Update", func() {
		var (
			productSlug string
			fileGroupID int
			name        *string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			fileGroupID = fileGroups[0].ID

			fakePivnetClient.FileGroupReturns(fileGroups[0], nil)
			fakePivnetClient.UpdateFileGroupReturns(fileGroups[0], nil)
		})

		It("updates FileGroup", func() {
			err := client.Update(productSlug, fileGroupID, name)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("filegroup error")
				fakePivnetClient.UpdateFileGroupReturns(pivnet.FileGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Update(productSlug, fileGroupID, name)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when the name is non-nil", func() {
			var (
				nameVal string
			)

			BeforeEach(func() {
				nameVal = "some-name"
				name = &nameVal
			})

			It("updates FileGroup with the provided name", func() {
				err := client.Update(productSlug, fileGroupID, name)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetClient.UpdateFileGroupCallCount()).To(Equal(1))
				_, invokedFileGroup := fakePivnetClient.UpdateFileGroupArgsForCall(0)

				Expect(invokedFileGroup.ID).To(Equal(fileGroupID))
				Expect(invokedFileGroup.Name).To(Equal(nameVal))
			})
		})
	})

	Describe("Delete", func() {
		var (
			productSlug string
			fileGroupID int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			fileGroupID = fileGroups[0].ID

			fakePivnetClient.DeleteFileGroupReturns(fileGroups[0], nil)
		})

		It("deletes FileGroup", func() {
			err := client.Delete(productSlug, fileGroupID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("filegroup error")
				fakePivnetClient.DeleteFileGroupReturns(pivnet.FileGroup{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Delete(productSlug, fileGroupID)
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
			fileGroupID    int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			fileGroupID = fileGroups[0].ID

			fakePivnetClient.AddFileGroupToReleaseReturns(nil)
		})

		It("adds FileGroup to release", func() {
			err := client.AddToRelease(
				productSlug,
				fileGroupID,
				releaseVersion,
			)
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
				err := client.AddToRelease(
					productSlug,
					fileGroupID,
					releaseVersion,
				)
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
				expectedErr = errors.New("file group error")
				fakePivnetClient.AddFileGroupToReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.AddToRelease(
					productSlug,
					fileGroupID,
					releaseVersion,
				)
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
			fileGroupID    int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = "release-version"
			fileGroupID = fileGroups[0].ID

			fakePivnetClient.RemoveFileGroupFromReleaseReturns(nil)
		})

		It("removes FileGroup from release", func() {
			err := client.RemoveFromRelease(
				productSlug,
				fileGroupID,
				releaseVersion,
			)
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
				err := client.RemoveFromRelease(
					productSlug,
					fileGroupID,
					releaseVersion,
				)
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
				expectedErr = errors.New("file group error")
				fakePivnetClient.RemoveFileGroupFromReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.RemoveFromRelease(
					productSlug,
					fileGroupID,
					releaseVersion,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})
