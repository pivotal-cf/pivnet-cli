package rc_test

import (
	"fmt"
	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/v2/rc"
	"github.com/pivotal-cf/pivnet-cli/v2/rc/rcfakes"
)

var _ = Describe("RCHandler", func() {
	var (
		fakePivnetRCReadWriter *rcfakes.FakePivnetRCReadWriter

		rcHandler *rc.RCHandler

		profile rc.PivnetProfile

		configContents []byte

		readErr error
	)

	BeforeEach(func() {
		fakePivnetRCReadWriter = &rcfakes.FakePivnetRCReadWriter{}

		readErr = nil

		profile = rc.PivnetProfile{
			Name:              "some-profile",
			APIToken:          "some-api-token",
			Host:              "some-host",
			AccessToken:       "some-access-token",
			AccessTokenExpiry: 12345,
		}

		configContents = []byte(fmt.Sprintf(
			`
---
profiles:
- name: %s
  api_token: %s
  host: %s
  access_token: %s
  access_token_expiry: %d
`,
			profile.Name,
			profile.APIToken,
			profile.Host,
			profile.AccessToken,
			profile.AccessTokenExpiry,
		),
		)
	})

	JustBeforeEach(func() {
		fakePivnetRCReadWriter.ReadFromFileReturns(configContents, readErr)
		rcHandler = rc.NewRCHandler(fakePivnetRCReadWriter)
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

		Context("when rc file contents cannot be unmarshalled", func() {
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
		It("successfully saves profile", func() {
			updatedAPIToken := "updatedAPIToken"

			err := rcHandler.SaveProfile(
				profile.Name,
				updatedAPIToken,
				profile.Host,
				profile.AccessToken,
				profile.AccessTokenExpiry,
			)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePivnetRCReadWriter.WriteToFileCallCount()).To(Equal(1))

			invokedContents := fakePivnetRCReadWriter.WriteToFileArgsForCall(0)

			expectedPivnetRC := rc.PivnetRC{
				Profiles: []rc.PivnetProfile{
					{
						Name:              profile.Name,
						APIToken:          updatedAPIToken,
						Host:              profile.Host,
						AccessToken:       profile.AccessToken,
						AccessTokenExpiry: profile.AccessTokenExpiry,
					},
				},
			}

			expectedBytes, err := yaml.Marshal(expectedPivnetRC)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(invokedContents)).To(Equal(string(expectedBytes)))
		})

		Context("when profile does not yet exist", func() {
			var (
				newName              string
				newAPIToken          string
				newHost              string
				newAccessToken       string
				newAccessTokenExpiry int64
			)

			BeforeEach(func() {
				newName = "some-other-profile"
				newAPIToken = "some-other-api-token"
				newHost = "some-other-host"
				newAccessToken = "new_access_token"
				newAccessTokenExpiry = 1554219433
			})

			It("creates new profile without error", func() {
				err := rcHandler.SaveProfile(
					newName,
					newAPIToken,
					newHost,
					newAccessToken,
					newAccessTokenExpiry,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetRCReadWriter.WriteToFileCallCount()).To(Equal(1))

				invokedContents := fakePivnetRCReadWriter.WriteToFileArgsForCall(0)

				expectedPivnetRC := rc.PivnetRC{
					Profiles: []rc.PivnetProfile{
						profile,
						{
							Name:              newName,
							APIToken:          newAPIToken,
							Host:              newHost,
							AccessToken:       newAccessToken,
							AccessTokenExpiry: newAccessTokenExpiry,
						},
					},
				}

				expectedBytes, err := yaml.Marshal(expectedPivnetRC)
				Expect(err).NotTo(HaveOccurred())
				Expect(invokedContents).To(Equal(expectedBytes))
			})
		})

		Context("when profile file does not exist", func() {
			BeforeEach(func() {
				configContents = nil
			})

			It("saves new profile without error", func() {
				err := rcHandler.SaveProfile(
					profile.Name,
					profile.APIToken,
					profile.Host,
					profile.AccessToken,
					profile.AccessTokenExpiry,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetRCReadWriter.WriteToFileCallCount()).To(Equal(1))

				invokedContents := fakePivnetRCReadWriter.WriteToFileArgsForCall(0)

				expectedPivnetRC := rc.PivnetRC{
					Profiles: []rc.PivnetProfile{
						profile,
					},
				}

				expectedBytes, err := yaml.Marshal(expectedPivnetRC)
				Expect(err).NotTo(HaveOccurred())
				Expect(invokedContents).To(Equal(expectedBytes))
			})
		})

		Context("when rc file exists but cannot be read", func() {
			BeforeEach(func() {
				readErr = fmt.Errorf("some read error")
			})

			It("returns an error", func() {
				err := rcHandler.SaveProfile(
					profile.Name,
					profile.APIToken,
					profile.Host,
					profile.AccessToken,
					profile.AccessTokenExpiry,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when rc file contents cannot be unmarshalled", func() {
			BeforeEach(func() {
				configContents = []byte("[*invalid-yaml!")
			})

			It("returns an error", func() {
				err := rcHandler.SaveProfile(
					profile.Name,
					profile.APIToken,
					profile.Host,
					profile.AccessToken,
					profile.AccessTokenExpiry,
				)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("RemoveProfileWithName", func() {
		It("removes profile", func() {
			err := rcHandler.RemoveProfileWithName(profile.Name)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePivnetRCReadWriter.WriteToFileCallCount()).To(Equal(1))

			invokedContents := fakePivnetRCReadWriter.WriteToFileArgsForCall(0)

			expectedPivnetRC := rc.PivnetRC{
				Profiles: []rc.PivnetProfile{},
			}

			expectedBytes, err := yaml.Marshal(expectedPivnetRC)
			Expect(err).NotTo(HaveOccurred())
			Expect(invokedContents).To(Equal(expectedBytes))
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

				Expect(fakePivnetRCReadWriter.WriteToFileCallCount()).To(Equal(1))

				invokedContents := fakePivnetRCReadWriter.WriteToFileArgsForCall(0)

				expectedPivnetRC := rc.PivnetRC{
					Profiles: []rc.PivnetProfile{
						profile,
					},
				}

				expectedBytes, err := yaml.Marshal(expectedPivnetRC)
				Expect(err).NotTo(HaveOccurred())
				Expect(invokedContents).To(Equal(expectedBytes))
			})
		})

		Context("when rc file does not exist", func() {
			BeforeEach(func() {
				configContents = nil
			})

			It("does not write a file", func() {
				err := rcHandler.RemoveProfileWithName("unused name")
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePivnetRCReadWriter.WriteToFileCallCount()).To(Equal(0))
			})
		})

		Context("when reading rc file returns an error", func() {
			BeforeEach(func() {
				readErr = fmt.Errorf("some read error")
			})

			It("returns an error", func() {
				err := rcHandler.RemoveProfileWithName("unused name")

				Expect(err).To(Equal(readErr))
			})
		})

		Context("when rc file contents cannot be unmarshalled", func() {
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
