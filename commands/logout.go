package commands

import "github.com/pivotal-cf/pivnet-cli/commands/logout"

type LogoutCommand struct {
}

//go:generate counterfeiter . LogoutClient
type LogoutClient interface {
	Logout(profileName string) error
}

var NewLogoutClient = func() LogoutClient {
	return logout.NewLogoutClient(
		RC,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *LogoutCommand) Execute([]string) error {
	err := Init(false)
	if err != nil {
		return err
	}

	return NewLogoutClient().Logout(
		Pivnet.ProfileName,
	)
}
