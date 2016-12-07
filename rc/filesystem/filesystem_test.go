package filesystem_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/rc/filesystem"
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

			Context("when profile file cannot be read", func() {
				JustBeforeEach(func() {
					err := os.Chmod(configFilepath, 0)
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns an error", func() {
					_, err := rcReadWriter.ReadFromFile()

					Expect(err).To(HaveOccurred())
				})
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
			contents string

			existingContents []byte
		)

		BeforeEach(func() {
			contents = "some contents"

			existingContents = []byte(("some-existing-contents"))
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

				writtenContents := contentsFromRCFilepath(otherFilepath)

				Expect(writtenContents).To(Equal(contents))
			})

			It("overwrites existing file", func() {
				err := rcReadWriter.WriteToFile(contents)
				Expect(err).NotTo(HaveOccurred())

				writtenContents := contentsFromRCFilepath(configFilepath)

				Expect(writtenContents).To(Equal(contents))
			})

			It("writes file with user-only read/write (i.e. 0600) permissions", func() {
				err := rcReadWriter.WriteToFile(contents)
				Expect(err).NotTo(HaveOccurred())

				info, err := os.Stat(configFilepath)
				Expect(err).NotTo(HaveOccurred())

				Expect(info.Mode()).To(Equal(os.FileMode(0600)))
			})
		})
	})
})

func contentsFromRCFilepath(filepath string) string {
	b, err := ioutil.ReadFile(filepath)
	Expect(err).NotTo(HaveOccurred())

	var contents string
	err = yaml.Unmarshal(b, &contents)
	Expect(err).NotTo(HaveOccurred())

	Expect(contents).NotTo(BeEmpty())

	return contents
}
