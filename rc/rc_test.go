package rc_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/rc"
	"github.com/pivotal-cf/pivnet-cli/rc/rcfakes"
)

var _ = Describe("RCHandler", func() {
	var (
		fakePivnetRCWriter *rcfakes.FakePivnetRCWriter
		rcHandler          *rc.RCHandler

		profile rc.PivnetProfile

		tempDir        string
		configContents []byte
		configFilepath string
	)

	BeforeEach(func() {
		fakePivnetRCWriter = &rcfakes.FakePivnetRCWriter{}

		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		profile = rc.PivnetProfile{
			Name:     "some-profile",
			APIToken: "some-api-token",
			Host:     "some-host",
		}

		configFilepath = filepath.Join(tempDir, ".pivnetrc")
		configContents = []byte(fmt.Sprintf(
			`
---
profiles:
- name: %s
  api_token: %s
  host: %s
`,
			profile.Name,
			profile.APIToken,
			profile.Host,
		),
		)
	})

	JustBeforeEach(func() {
		err := ioutil.WriteFile(configFilepath, configContents, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		rcHandler = rc.NewRCHandler(configFilepath, fakePivnetRCWriter)
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
			Expect(returnedProfile.Host).To(Equal(profile.Host))
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
				rcHandler = rc.NewRCHandler(otherFilepath, fakePivnetRCWriter)
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
			updatedAPIToken := "updatedAPIToken"

			err := rcHandler.SaveProfile(
				profile.Name,
				updatedAPIToken,
				profile.Host,
			)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePivnetRCWriter.WriteToFileCallCount()).To(Equal(1))

			invokedFilepath, invokedContents := fakePivnetRCWriter.WriteToFileArgsForCall(0)
			Expect(invokedFilepath).To(Equal(configFilepath))

			expectedPivnetRC := rc.PivnetRC{
				Profiles: []rc.PivnetProfile{
					{
						Name:     profile.Name,
						APIToken: updatedAPIToken,
						Host:     profile.Host,
					},
				},
			}
			Expect(invokedContents).To(Equal(expectedPivnetRC))
		})

		Context("when profile does not yet exist", func() {
			var (
				newName     string
				newAPIToken string
				newHost     string
			)

			BeforeEach(func() {
				newName = "some-other-profile"
				newAPIToken = "some-other-api-token"
				newHost = "some-other-host"
			})

			It("creates new profile without error", func() {
				err := rcHandler.SaveProfile(
					newName,
					newAPIToken,
					newHost,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetRCWriter.WriteToFileCallCount()).To(Equal(1))

				invokedFilepath, invokedContents := fakePivnetRCWriter.WriteToFileArgsForCall(0)
				Expect(invokedFilepath).To(Equal(configFilepath))

				expectedPivnetRC := rc.PivnetRC{
					Profiles: []rc.PivnetProfile{
						profile,
						{
							Name:     newName,
							APIToken: newAPIToken,
							Host:     newHost,
						},
					},
				}
				Expect(invokedContents).To(Equal(expectedPivnetRC))
			})
		})

		Context("when profile file does not exist", func() {
			var (
				newFilepath string
			)

			JustBeforeEach(func() {
				newFilepath = filepath.Join(tempDir, "other-file")
				rcHandler = rc.NewRCHandler(newFilepath, fakePivnetRCWriter)
			})

			It("saves new profile without error", func() {
				err := rcHandler.SaveProfile(
					profile.Name,
					profile.APIToken,
					profile.Host,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetRCWriter.WriteToFileCallCount()).To(Equal(1))

				invokedFilepath, invokedContents := fakePivnetRCWriter.WriteToFileArgsForCall(0)
				Expect(invokedFilepath).To(Equal(newFilepath))

				expectedPivnetRC := rc.PivnetRC{
					Profiles: []rc.PivnetProfile{
						profile,
					},
				}
				Expect(invokedContents).To(Equal(expectedPivnetRC))
			})
		})

		Context("when profile file exists but cannot be read", func() {
			JustBeforeEach(func() {
				err := os.Chmod(configFilepath, 0)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error", func() {
				err := rcHandler.SaveProfile(
					profile.Name,
					profile.APIToken,
					profile.Host,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when profile file cannot be unmarshalled", func() {
			BeforeEach(func() {
				configContents = []byte("[*invalid-yaml!")
			})

			It("returns an error", func() {
				err := rcHandler.SaveProfile(
					profile.Name,
					profile.APIToken,
					profile.Host,
				)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("RemoveProfileWithName", func() {
		It("removes profile", func() {
			err := rcHandler.RemoveProfileWithName(profile.Name)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePivnetRCWriter.WriteToFileCallCount()).To(Equal(1))

			invokedFilepath, invokedContents := fakePivnetRCWriter.WriteToFileArgsForCall(0)
			Expect(invokedFilepath).To(Equal(configFilepath))

			expectedPivnetRC := rc.PivnetRC{
				Profiles: []rc.PivnetProfile{},
			}
			Expect(invokedContents).To(Equal(expectedPivnetRC))
		})

		Context("when profile does not yet exist", func() {
			var (
				newName string
			)

			BeforeEach(func() {
				newName = "some-other-profile"
			})

			It("writes existing profiles without error", func() {
				err := rcHandler.RemoveProfileWithName(newName)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetRCWriter.WriteToFileCallCount()).To(Equal(1))

				invokedFilepath, invokedContents := fakePivnetRCWriter.WriteToFileArgsForCall(0)
				Expect(invokedFilepath).To(Equal(configFilepath))

				expectedPivnetRC := rc.PivnetRC{
					Profiles: []rc.PivnetProfile{
						profile,
					},
				}
				Expect(invokedContents).To(Equal(expectedPivnetRC))
			})
		})

		Context("when profile file does not exist", func() {
			var (
				newFilepath string
			)

			JustBeforeEach(func() {
				newFilepath = filepath.Join(tempDir, "other-file")
				rcHandler = rc.NewRCHandler(newFilepath, fakePivnetRCWriter)
			})

			It("does not write a file", func() {
				err := rcHandler.RemoveProfileWithName("unused name")
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetRCWriter.WriteToFileCallCount()).To(Equal(0))
			})
		})

		Context("when profile file exists but cannot be read", func() {
			JustBeforeEach(func() {
				err := os.Chmod(configFilepath, 0)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error", func() {
				err := rcHandler.RemoveProfileWithName("unused name")

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when profile file cannot be unmarshalled", func() {
			BeforeEach(func() {
				configContents = []byte("[*invalid-yaml!")
			})

			It("returns an error", func() {
				err := rcHandler.RemoveProfileWithName("unused name")

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
