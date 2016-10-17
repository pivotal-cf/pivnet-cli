package commands

import "github.com/pivotal-cf/go-pivnet/cmd/pivnet/commands/releasedependency"

type ReleaseDependenciesCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type AddReleaseDependencyCommand struct {
	ProductSlug             string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion          string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	DependentProductSlug    string `long:"dependent-product-slug" description:"Dependent product slug e.g. p-mysql" required:"true"`
	DependentReleaseVersion string `long:"dependent-release-version" description:"Dependent release version e.g. 0.1.2-rc1" required:"true"`
}

type RemoveReleaseDependencyCommand struct {
	ProductSlug             string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion          string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	DependentProductSlug    string `long:"dependent-product-slug" description:"Dependent product slug e.g. p-mysql" required:"true"`
	DependentReleaseVersion string `long:"dependent-release-version" description:"Dependent release version e.g. 0.1.2-rc1" required:"true"`
}

//go:generate counterfeiter . ReleaseDependencyClient
type ReleaseDependencyClient interface {
	List(productSlug string, releaseVersion string) error
	Add(productSlug string, releaseVersion string, dependentProductSlug string, dependentReleaseVersion string) error
	Remove(productSlug string, releaseVersion string, dependentProductSlug string, dependentReleaseVersion string) error
}

var NewReleaseDependencyClient = func() ReleaseDependencyClient {
	return releasedependency.NewReleaseDependencyClient(
		NewPivnetClient(),
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *ReleaseDependenciesCommand) Execute([]string) error {
	Init()

	return NewReleaseDependencyClient().List(command.ProductSlug, command.ReleaseVersion)
}

func (command *AddReleaseDependencyCommand) Execute([]string) error {
	Init()

	return NewReleaseDependencyClient().Add(
		command.ProductSlug,
		command.ReleaseVersion,
		command.DependentProductSlug,
		command.DependentReleaseVersion,
	)
}

func (command *RemoveReleaseDependencyCommand) Execute([]string) error {
	Init()

	return NewReleaseDependencyClient().Remove(
		command.ProductSlug,
		command.ReleaseVersion,
		command.DependentProductSlug,
		command.DependentReleaseVersion,
	)
}
