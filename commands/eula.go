package commands

import "github.com/pivotal-cf/pivnet-cli/v2/commands/eula"

type EULAsCommand struct {
}

type EULACommand struct {
	EULASlug string `long:"eula-slug" description:"EULA slug e.g. pivotal_software_eula" required:"true"`
}

type AcceptEULACommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

//go:generate counterfeiter . EULAClient
type EULAClient interface {
	List() error
	Get(eulaSlug string) error
	AcceptEULA(productSlug string, releaseVersion string) error
}

var NewEULAClient = func(client eula.PivnetClient) EULAClient {
	return eula.NewEULAClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *EULAsCommand) Execute(args []string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewEULAClient(client).List()
}

func (command *EULACommand) Execute(args []string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewEULAClient(client).Get(command.EULASlug)
}

func (command *AcceptEULACommand) Execute(args []string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewEULAClient(client).AcceptEULA(command.ProductSlug, command.ReleaseVersion)
}
