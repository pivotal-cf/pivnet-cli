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

var _ = Describe("PivnetRCWriter", func() {
	var (
		rcWriter *filesystem.PivnetRCWriter

		contents string

		tempDir        string
		configContents []byte
		configFilepath string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		contents = "some contents"

		configFilepath = filepath.Join(tempDir, ".pivnetrc")
		configContents = []byte(("some-existing-contents"))
	})

	JustBeforeEach(func() {
		err := ioutil.WriteFile(configFilepath, configContents, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		rcWriter = filesystem.NewPivnetRCWriter()
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("WriteToFile", func() {
		It("writes to file", func() {
			err := rcWriter.WriteToFile(
				configFilepath,
				contents,
			)
			Expect(err).NotTo(HaveOccurred())

			writtenContents := contentsFromRCFilepath(configFilepath)

			Expect(writtenContents).To(Equal(contents))
		})

		It("writes file with user-only read/write (i.e. 0600) permissions", func() {
			err := rcWriter.WriteToFile(
				configFilepath,
				contents,
			)
			Expect(err).NotTo(HaveOccurred())

			info, err := os.Stat(configFilepath)
			Expect(err).NotTo(HaveOccurred())

			Expect(info.Mode()).To(Equal(os.FileMode(0600)))
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
