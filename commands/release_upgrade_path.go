package commands

import "github.com/pivotal-cf/pivnet-cli/v3/commands/releaseupgradepath"

type ReleaseUpgradePathsCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type AddReleaseUpgradePathCommand struct {
	ProductSlug            string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion         string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	PreviousReleaseVersion string `long:"previous-release-version" short:"u" description:"Regex for previous release version e.g. 0.1.*" required:"true"`
}

type RemoveReleaseUpgradePathCommand struct {
	ProductSlug            string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion         string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	PreviousReleaseVersion string `long:"previous-release-version" short:"u" description:"Regex for previous release version e.g. 0.1.*" required:"true"`
}

//go:generate counterfeiter . ReleaseUpgradePathClient
type ReleaseUpgradePathClient interface {
	List(productSlug string, releaseVersion string) error
	Add(productSlug string, releaseVersion string, previousReleaseVersion string) error
	Remove(productSlug string, releaseVersion string, previousReleaseVersion string) error
}

var NewReleaseUpgradePathClient = func(client releaseupgradepath.PivnetClient) ReleaseUpgradePathClient {
	return releaseupgradepath.NewReleaseUpgradePathClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
		Filter,
		Pivnet.Logger,
	)
}

func (command *ReleaseUpgradePathsCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseUpgradePathClient(client).List(command.ProductSlug, command.ReleaseVersion)
}

func (command *AddReleaseUpgradePathCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseUpgradePathClient(client).Add(
		command.ProductSlug,
		command.ReleaseVersion,
		command.PreviousReleaseVersion,
	)
}

func (command *RemoveReleaseUpgradePathCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseUpgradePathClient(client).Remove(
		command.ProductSlug,
		command.ReleaseVersion,
		command.PreviousReleaseVersion,
	)
}
