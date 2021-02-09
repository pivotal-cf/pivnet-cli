package filesystem_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/v3/rc/filesystem"
)

var _ = Describe("PivnetRCReadWriter", func() {
	var (
		rcReadWriter *filesystem.PivnetRCReadWriter

		tempDir        string
		configFilepath string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		configFilepath = filepath.Join(tempDir, ".pivnetrc")

		rcReadWriter = filesystem.NewPivnetRCReadWriter(configFilepath)
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ReadFromFile", func() {
		var (
			contents []byte
		)

		BeforeEach(func() {
			contents = []byte("some contents")

			err := ioutil.WriteFile(configFilepath, contents, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("ReadFromFile", func() {
			It("reads from file", func() {
				b, err := rcReadWriter.ReadFromFile()
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(Equal(contents))
			})

			Context("when profile file does not exist", func() {
				It("returns empty profile without error", func() {
					otherFilepath := filepath.Join(tempDir, "other-file")
					rcReadWriter = filesystem.NewPivnetRCReadWriter(otherFilepath)

					b, err := rcReadWriter.ReadFromFile()

					Expect(err).NotTo(HaveOccurred())
					Expect(b).To(BeNil())
				})
			})
		})
	})

	Describe("Write", func() {
		var (
			contents []byte

			existingContents []byte
		)

		BeforeEach(func() {
			contents = []byte("some contents")

			existingContents = []byte("some-existing-contents")
		})

		JustBeforeEach(func() {
			err := ioutil.WriteFile(configFilepath, existingContents, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("WriteToFile", func() {
			It("creates new file", func() {
				otherFilepath := filepath.Join(tempDir, "other-file")
				rcReadWriter = filesystem.NewPivnetRCReadWriter(otherFilepath)

				err := rcReadWriter.WriteToFile(contents)
				Expect(err).NotTo(HaveOccurred())

				b, err := ioutil.ReadFile(otherFilepath)
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(Equal(contents))
			})

			It("overwrites existing file", func() {
				err := rcReadWriter.WriteToFile(contents)
				Expect(err).NotTo(HaveOccurred())

				b, err := ioutil.ReadFile(configFilepath)
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(Equal(contents))
			})
		})
	})
})
