package rc_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/rc"
)

var _ = Describe("RCHandler", func() {
	var (
		rcHandler *rc.RCHandler

		profile rc.PivnetProfile

		tempDir        string
		configContents []byte
		configFilepath string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		profile = rc.PivnetProfile{
			Name:     "some-profile",
			APIToken: "some-api-token",
		}

		configFilepath = filepath.Join(tempDir, ".pivnetrc")
		configContents = []byte(fmt.Sprintf(
			`
---
profiles:
- name: %s
  api_token: %s
`,
			profile.Name,
			profile.APIToken,
		),
		)
	})

	JustBeforeEach(func() {
		err := ioutil.WriteFile(configFilepath, configContents, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		rcHandler = rc.NewRCHandler(configFilepath)
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ProfileForName", func() {
		It("returns located profile", func() {
			returnedProfile, err := rcHandler.ProfileForName(profile.Name)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedProfile).NotTo(BeNil())
			Expect(returnedProfile.Name).To(Equal(profile.Name))
			Expect(returnedProfile.APIToken).To(Equal(profile.APIToken))
		})

		Context("when profile cannot be found", func() {
			It("returns nil profile without error", func() {
				returnedProfile, err := rcHandler.ProfileForName("some-other-profile")
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedProfile).To(BeNil())
			})
		})

		Context("when profile file does not exist", func() {
			JustBeforeEach(func() {
				otherFilepath := filepath.Join(tempDir, "other-file")
				rcHandler = rc.NewRCHandler(otherFilepath)
			})

			It("returns empty profile without error", func() {
				_, err := rcHandler.ProfileForName("some-profile")

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when profile file cannot be read", func() {
			JustBeforeEach(func() {
				err := os.Chmod(configFilepath, 0)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error", func() {
				_, err := rcHandler.ProfileForName("some-profile")

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when profile file cannot be unmarshalled", func() {
			BeforeEach(func() {
				configContents = []byte("[*invalid-yaml!")
			})

			It("returns an error", func() {
				_, err := rcHandler.ProfileForName("some-profile")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("SaveProfile", func() {
		It("saves profile", func() {
			err := rcHandler.SaveProfile(profile.Name, profile.APIToken)
			Expect(err).NotTo(HaveOccurred())

			profiles := profilesFromRCFilepath(configFilepath)

			Expect(profiles).To(HaveLen(1))

			Expect(profiles[0]).NotTo(BeNil())
			Expect(profiles[0].Name).To(Equal(profile.Name))
			Expect(profiles[0].APIToken).To(Equal(profile.APIToken))
		})

		It("updates existing file with user-only read/write (i.e. 0600) permissions", func() {
			err := rcHandler.SaveProfile(profile.Name, profile.APIToken)
			Expect(err).NotTo(HaveOccurred())

			info, err := os.Stat(configFilepath)
			Expect(err).NotTo(HaveOccurred())

			Expect(info.Mode()).To(Equal(os.FileMode(0600)))
		})

		Context("when profile does not yet exist", func() {
			var (
				newName     string
				newAPIToken string
			)

			BeforeEach(func() {
				newName = "some-other-profile"
				newAPIToken = "some-other-api-token"
			})

			It("creates new profile without error", func() {
				err := rcHandler.SaveProfile(newName, newAPIToken)
				Expect(err).NotTo(HaveOccurred())

				profiles := profilesFromRCFilepath(configFilepath)

				Expect(profiles).To(HaveLen(2))

				Expect(profiles[0]).NotTo(BeNil())
				Expect(profiles[0].Name).To(Equal(profile.Name))
				Expect(profiles[0].APIToken).To(Equal(profile.APIToken))

				Expect(profiles[1]).NotTo(BeNil())
				Expect(profiles[1].Name).To(Equal(newName))
				Expect(profiles[1].APIToken).To(Equal(newAPIToken))
			})
		})

		Context("when profile file does not exist", func() {
			var (
				newFilepath string
			)

			JustBeforeEach(func() {
				newFilepath = filepath.Join(tempDir, "other-file")
				rcHandler = rc.NewRCHandler(newFilepath)
			})

			It("create file and saves new profile without error", func() {
				err := rcHandler.SaveProfile(profile.Name, profile.APIToken)
				Expect(err).NotTo(HaveOccurred())

				profiles := profilesFromRCFilepath(newFilepath)

				Expect(profiles).To(HaveLen(1))

				Expect(profiles[0]).NotTo(BeNil())
				Expect(profiles[0].Name).To(Equal(profile.Name))
				Expect(profiles[0].APIToken).To(Equal(profile.APIToken))
			})

			It("creates new file with user-only read/write (i.e. 0600) permissions", func() {
				err := rcHandler.SaveProfile(profile.Name, profile.APIToken)
				Expect(err).NotTo(HaveOccurred())

				info, err := os.Stat(newFilepath)
				Expect(err).NotTo(HaveOccurred())

				Expect(info.Mode()).To(Equal(os.FileMode(0600)))
			})
		})

		Context("when profile file exists but cannot be read", func() {
			JustBeforeEach(func() {
				err := os.Chmod(configFilepath, 0)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error", func() {
				err := rcHandler.SaveProfile(profile.Name, profile.APIToken)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when profile file cannot be unmarshalled", func() {
			BeforeEach(func() {
				configContents = []byte("[*invalid-yaml!")
			})

			It("returns an error", func() {
				err := rcHandler.SaveProfile(profile.Name, profile.APIToken)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})

func profilesFromRCFilepath(filepath string) []rc.PivnetProfile {
	pivnetRCBytes, err := ioutil.ReadFile(filepath)
	Expect(err).NotTo(HaveOccurred())

	var pivnetRC rc.PivnetRC
	err = yaml.Unmarshal(pivnetRCBytes, &pivnetRC)
	Expect(err).NotTo(HaveOccurred())

	Expect(pivnetRC).NotTo(BeNil())

	return pivnetRC.Profiles
}
