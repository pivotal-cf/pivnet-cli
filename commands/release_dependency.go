package commands

import "github.com/pivotal-cf/pivnet-cli/v3/commands/releasedependency"

type ReleaseDependenciesCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type AddReleaseDependencyCommand struct {
	ProductSlug             string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion          string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	DependentProductSlug    string `long:"dependent-product-slug" short:"s" description:"Dependent product slug e.g. p-mysql" required:"true"`
	DependentReleaseVersion string `long:"dependent-release-version" short:"u" description:"Dependent release version e.g. 0.1.2-rc1" required:"true"`
}

type RemoveReleaseDependencyCommand struct {
	ProductSlug             string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion          string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	DependentProductSlug    string `long:"dependent-product-slug" short:"s" description:"Dependent product slug e.g. p-mysql" required:"true"`
	DependentReleaseVersion string `long:"dependent-release-version" short:"u" description:"Dependent release version e.g. 0.1.2-rc1" required:"true"`
}

//go:generate counterfeiter . ReleaseDependencyClient
type ReleaseDependencyClient interface {
	List(productSlug string, releaseVersion string) error
	Add(productSlug string, releaseVersion string, dependentProductSlug string, dependentReleaseVersion string) error
	Remove(productSlug string, releaseVersion string, dependentProductSlug string, dependentReleaseVersion string) error
}

var NewReleaseDependencyClient = func(client releasedependency.PivnetClient) ReleaseDependencyClient {
	return releasedependency.NewReleaseDependencyClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *ReleaseDependenciesCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseDependencyClient(client).List(command.ProductSlug, command.ReleaseVersion)
}

func (command *AddReleaseDependencyCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseDependencyClient(client).Add(
		command.ProductSlug,
		command.ReleaseVersion,
		command.DependentProductSlug,
		command.DependentReleaseVersion,
	)
}

func (command *RemoveReleaseDependencyCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseDependencyClient(client).Remove(
		command.ProductSlug,
		command.ReleaseVersion,
		command.DependentProductSlug,
		command.DependentReleaseVersion,
	)
}
