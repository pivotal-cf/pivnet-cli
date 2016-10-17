package commands

import "github.com/pivotal-cf/go-pivnet/cmd/pivnet/commands/release"

type ReleasesCommand struct {
	ProductSlug string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
}

type ReleaseCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type DeleteReleaseCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type CreateReleaseCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	ReleaseType    string `long:"release-type" short:"t" description:"Release type e.g. 'Minor Release'" required:"true"`
	EULASlug       string `long:"eula-slug" short:"e" description:"EULA slug e.g. pivotal_software_eula" required:"true"`
}

//go:generate counterfeiter . ReleaseClient
type ReleaseClient interface {
	List(productSlug string) error
	Get(productSlug string, releaseVersion string) error
	Create(productSlug string, releaseVersion string, releaseType string, eulaSlug string) error
	Delete(productSlug string, releaseVersion string) error
}

var NewReleaseClient = func() ReleaseClient {
	return release.NewReleaseClient(
		NewPivnetClient(),
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *ReleasesCommand) Execute([]string) error {
	Init()

	return NewReleaseClient().List(command.ProductSlug)
}

func (command *ReleaseCommand) Execute([]string) error {
	Init()

	return NewReleaseClient().Get(command.ProductSlug, command.ReleaseVersion)
}

func (command *CreateReleaseCommand) Execute([]string) error {
	Init()

	return NewReleaseClient().Create(
		command.ProductSlug,
		command.ReleaseVersion,
		command.ReleaseType,
		command.EULASlug,
	)
}

func (command *DeleteReleaseCommand) Execute([]string) error {
	Init()

	return NewReleaseClient().Delete(command.ProductSlug, command.ReleaseVersion)
}
