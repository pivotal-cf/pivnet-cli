package commands

import "github.com/pivotal-cf/pivnet-cli/commands/pivnetversions"

type PivnetVersionsCommand struct {
}

//go:generate counterfeiter . PivnetVersionsClient
type PivnetVersionsClient interface {
	List() error
}

var NewPivnetVersionsClient = func(client pivnetversions.PivnetClient) PivnetVersionsClient {
	return pivnetversions.NewPivnetVersionsClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *PivnetVersionsCommand) Execute(args []string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewPivnetVersionsClient(client).List()
}
