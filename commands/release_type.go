package commands

import "github.com/pivotal-cf/pivnet-cli/commands/releasetype"

type ReleaseTypesCommand struct {
}

//go:generate counterfeiter . ReleaseTypeClient
type ReleaseTypeClient interface {
	List() error
}

var NewReleaseTypeClient = func() ReleaseTypeClient {
	return releasetype.NewReleaseTypeClient(
		NewPivnetClient(),
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *ReleaseTypesCommand) Execute(args []string) error {
	Init()

	return NewReleaseTypeClient().List()
}
