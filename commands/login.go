package commands

import (
	"github.com/pivotal-cf/pivnet-cli/commands/login"
)

type LoginCommand struct {
	APIToken string `long:"api-token" description:"Pivnet API Token (Pivnet legacy token or UAA refresh token)" required:"true"`
	Host     string `long:"host" description:"Pivnet API Host" default:"https://network.pivotal.io"`
}

//go:generate counterfeiter . LoginClient
type LoginClient interface {
	Login(name string, apiToken string, host string) error
}

var NewLoginClient = func(client login.PivnetClient) LoginClient {
	return login.NewLoginClient(
		client,
		RC,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *LoginCommand) Execute([]string) error {
	err := Init(false)
	if err != nil {
		return err
	}

	sanitizeWriters(command.APIToken)

	accessTokenService := CreateAccessTokenService(RC, Pivnet.ProfileName, command.APIToken, command.Host)

	client := NewPivnetClientWithToken(accessTokenService, command.Host)

	return NewLoginClient(client).Login(
		Pivnet.ProfileName,
		command.APIToken,
		command.Host,
	)
}
