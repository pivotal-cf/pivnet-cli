package commands

import "github.com/pivotal-cf/go-pivnet/cmd/pivnet/commands/releaseupgradepath"

type ReleaseUpgradePathsCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type AddReleaseUpgradePathCommand struct {
	ProductSlug            string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion         string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	PreviousReleaseVersion string `long:"previous-release-version" description:"Previous release version e.g. 0.1.2-rc1" required:"true"`
}

type RemoveReleaseUpgradePathCommand struct {
	ProductSlug            string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion         string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	PreviousReleaseVersion string `long:"previous-release-version" description:"Previous release version e.g. 0.1.2-rc1" required:"true"`
}

//go:generate counterfeiter . ReleaseUpgradePathClient
type ReleaseUpgradePathClient interface {
	List(productSlug string, releaseVersion string) error
	Add(productSlug string, releaseVersion string, previousReleaseVersion string) error
	Remove(productSlug string, releaseVersion string, previousReleaseVersion string) error
}

var NewReleaseUpgradePathClient = func() ReleaseUpgradePathClient {
	return releaseupgradepath.NewReleaseUpgradePathClient(
		NewPivnetClient(),
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *ReleaseUpgradePathsCommand) Execute([]string) error {
	Init()

	return NewReleaseUpgradePathClient().List(command.ProductSlug, command.ReleaseVersion)
}

func (command *AddReleaseUpgradePathCommand) Execute([]string) error {
	Init()

	return NewReleaseUpgradePathClient().Add(
		command.ProductSlug,
		command.ReleaseVersion,
		command.PreviousReleaseVersion,
	)
}

func (command *RemoveReleaseUpgradePathCommand) Execute([]string) error {
	Init()

	return NewReleaseUpgradePathClient().Remove(
		command.ProductSlug,
		command.ReleaseVersion,
		command.PreviousReleaseVersion,
	)
}
