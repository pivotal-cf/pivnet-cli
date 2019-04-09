package commands

import (
	"fmt"
	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-cli/gp"
	"time"
)

var CreateAccessTokenService = func(
	rc RCHandler,
	profileName string,
	refreshToken string,
	host string,
) gp.AccessTokenService {
	tokenService := pivnet.NewAccessTokenOrLegacyToken(refreshToken, host)
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
	if err == nil && pivnetProfile != nil && pivnetProfile.AccessToken != "" && pivnetProfile.AccessTokenExpiry > time.Now().Unix() {
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

func CreateSaveTokenDecorator(rc RCHandler, accessTokenService gp.AccessTokenService, profileName string, refreshToken string, host string) gp.AccessTokenService {
	return SaveTokenDecorator {
		WrappedService: accessTokenService,
		ProfileName: profileName,
		RefreshToken: refreshToken,
		Host: host,
		Rc: rc,
	}
}