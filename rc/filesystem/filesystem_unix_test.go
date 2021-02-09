// +build !windows

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

