package commands

import (
	"fmt"
	"github.com/pivotal-cf/go-pivnet/v5"
	"github.com/pivotal-cf/pivnet-cli/gp"
	"github.com/pivotal-cf/pivnet-cli/rc"
	"time"
)

var CreateAccessTokenService = func(
	rc RCHandler,
	profileName string,
	refreshToken string,
	host string,
	skipSSLValidation bool,
) gp.AccessTokenService {
	tokenService := pivnet.NewAccessTokenOrLegacyToken(refreshToken, host, skipSSLValidation, "Pivnet CLI")
	serviceThatSavesRc := CreateSaveTokenDecorator(rc, tokenService, profileName, refreshToken, host)
	return serviceThatSavesRc
}

type SaveTokenDecorator struct {
	WrappedService gp.AccessTokenService
	ProfileName string
	RefreshToken string
	Host string
	Rc RCHandler
}

func (o SaveTokenDecorator) AccessToken() (string, error) {
	pivnetProfile, err := o.Rc.ProfileForName(o.ProfileName)

	if err == nil && pivnetProfile != nil && validAccessToken(pivnetProfile) && o.notDuringLogin(pivnetProfile) {
		return pivnetProfile.AccessToken, nil
	}

	accessToken, err := o.WrappedService.AccessToken()
	if err != nil {
		return "", fmt.Errorf("could not get access token %s", err)
	}

	err = o.Rc.SaveProfile(o.ProfileName, o.RefreshToken, o.Host, accessToken, o.accessTokenExpiry())
	if err != nil {
		return "", fmt.Errorf("failed to save profile %s", err)
	}

	return accessToken, nil
}

func (o SaveTokenDecorator) accessTokenExpiry() int64 {
	return time.Now().Add(time.Hour).Unix()
}

func (o SaveTokenDecorator) notDuringLogin(pivnetProfile *rc.PivnetProfile) bool {
	return pivnetProfile.Host == o.Host && pivnetProfile.APIToken == o.RefreshToken
}

func CreateSaveTokenDecorator(rc RCHandler, accessTokenService gp.AccessTokenService, profileName string, refreshToken string, host string) gp.AccessTokenService {
	return SaveTokenDecorator {
		WrappedService: accessTokenService,
		ProfileName: profileName,
		RefreshToken: refreshToken,
		Host: host,
		Rc: rc,
	}
}

func validAccessToken(pivnetProfile *rc.PivnetProfile) bool {
	return pivnetProfile.AccessToken != "" && pivnetProfile.AccessTokenExpiry > time.Now().Unix()
}
