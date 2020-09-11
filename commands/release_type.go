package commands

import "github.com/pivotal-cf/pivnet-cli/v2/commands/releasetype"

type ReleaseTypesCommand struct {
}

//go:generate counterfeiter . ReleaseTypeClient
type ReleaseTypeClient interface {
	List() error
}

var NewReleaseTypeClient = func(client releasetype.PivnetClient) ReleaseTypeClient {
	return releasetype.NewReleaseTypeClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *ReleaseTypesCommand) Execute(args []string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseTypeClient(client).List()
}
