package commands

import "github.com/pivotal-cf/go-pivnet/cmd/pivnet/commands/eula"

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

var NewEULAClient = func() EULAClient {
	return eula.NewEULAClient(
		NewPivnetClient(),
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *EULAsCommand) Execute(args []string) error {
	Init()

	return NewEULAClient().List()
}

func (command *EULACommand) Execute(args []string) error {
	Init()

	return NewEULAClient().Get(command.EULASlug)
}

func (command *AcceptEULACommand) Execute(args []string) error {
	Init()

	return NewEULAClient().AcceptEULA(command.ProductSlug, command.ReleaseVersion)
}
