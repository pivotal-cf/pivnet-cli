package commands_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/rc"
	"time"
)

var _ = Describe("Access Token Service", func() {
	var subject commands.SaveTokenDecorator
	var fakeRCHandler *commandsfakes.FakeRCHandler

	BeforeEach(func() {
		fakeRCHandler = &commandsfakes.FakeRCHandler{}

		subject = commands.SaveTokenDecorator{
			WrappedService: fakeAccessTokenService,
			ProfileName: "",
			RefreshToken: "",
			Host: "",
			Rc: fakeRCHandler,
		}
	})

	Context("when Profile does not have access token", func() {
		BeforeEach(func() {
			profile := &rc.PivnetProfile{AccessToken: ""}
			fakeRCHandler.ProfileForNameReturns(profile, nil)
		})

		It("saves profile with a new access token", func() {
			fakeAccessTokenService.AccessTokenReturns("token", nil)
			fakeRCHandler.SaveProfileReturns(nil)

			accessToken, err := subject.AccessToken()

			Expect(err).NotTo(HaveOccurred())
			Expect(accessToken).To(Equal("token"))
		})

		It("returns error if failed to fetch access token", func() {
			fetchAccessTokenError := fmt.Errorf("fetch access token error")
			fakeAccessTokenService.AccessTokenReturns("", fetchAccessTokenError)

			_, err := subject.AccessToken()

			Expect(err).To(MatchError(fmt.Errorf("could not get access token %s", fetchAccessTokenError)))
		})

		It("returns error if failed to save profile", func() {
			fakeAccessTokenService.AccessTokenReturns("token", nil)

			saveProfileError := fmt.Errorf("save profile error")
			fakeRCHandler.SaveProfileReturns(saveProfileError)

			_, err := subject.AccessToken()

			Expect(err).To(MatchError(fmt.Errorf("failed to save profile %s", saveProfileError)))
		})
	})

	Context("When Profile has access token", func() {
		var profile *rc.PivnetProfile
		BeforeEach(func() {
			profile = &rc.PivnetProfile{AccessToken: "newToken"}
			fakeRCHandler.ProfileForNameReturns(profile, nil)
		})

		Context("when access token is expired", func() {
			JustBeforeEach(func() {
				fakeRCHandler.ProfileForNameReturns(nil, nil)

				fakeAccessTokenService.AccessTokenReturns("token", nil)
				fakeRCHandler.SaveProfileReturns(nil)

				profile.AccessTokenExpiry = time.Now().Add(-time.Hour).Unix()
			})

			It("saves profile with new access token", func() {
				accessToken, err := subject.AccessToken()

				Expect(err).NotTo(HaveOccurred())
				Expect(accessToken).To(Equal("token"))
			})
		})

		Context("when access is not expired", func() {
			JustBeforeEach(func() {
				profile.AccessTokenExpiry = time.Now().Add(time.Hour).Unix()
			})

			It("returns existing access token", func() {
				accessToken, err := subject.AccessToken()

				Expect(err).NotTo(HaveOccurred())
				Expect(accessToken).To(Equal("newToken"))
			})
		})
	})

	Context("When failed to fetch Profile", func() {
		BeforeEach(func() {
			fetchProfileError := fmt.Errorf("fetch profile error")
			fakeRCHandler.ProfileForNameReturns(nil, fetchProfileError)

			fakeAccessTokenService.AccessTokenReturns("token", nil)
			fakeRCHandler.SaveProfileReturns(nil)
		})

		It("saves profile with a new access token", func() {

			accessToken, err := subject.AccessToken()

			Expect(err).NotTo(HaveOccurred())
			Expect(accessToken).To(Equal("token"))
		})
	})

	Context("When Profile is nil", func() {
		BeforeEach(func() {
			fakeRCHandler.ProfileForNameReturns(nil, nil)

			fakeAccessTokenService.AccessTokenReturns("token", nil)
			fakeRCHandler.SaveProfileReturns(nil)
		})

		It("saves profile with a new access token", func() {
			accessToken, err := subject.AccessToken()

			Expect(err).NotTo(HaveOccurred())
			Expect(accessToken).To(Equal("token"))
		})
	})
})
