package commands

import "github.com/pivotal-cf/pivnet-cli/v2/commands/dependencyspecifier"

type DependencySpecifiersCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type DependencySpecifierCommand struct {
	ProductSlug           string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion        string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	DependencySpecifierID int    `long:"dependency-specifier-id" short:"i" description:"Dependency specifier ID e.g. 1234" required:"true"`
}

type CreateDependencySpecifierCommand struct {
	ProductSlug          string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion       string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	DependentProductSlug string `long:"dependent-product-slug" short:"s" description:"Dependent product slug e.g. p-mysql" required:"true"`
	Specifier            string `long:"specifier" short:"u" description:"Specifier e.g. 1.2.*" required:"true"`
}

type DeleteDependencySpecifierCommand struct {
	ProductSlug           string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion        string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	DependencySpecifierID int    `long:"dependency-specifier-id" short:"i" description:"Dependency specifier ID e.g. 1234" required:"true"`
}

//go:generate counterfeiter . DependencySpecifierClient
type DependencySpecifierClient interface {
	List(productSlug string, releaseVersion string) error
	Get(productSlug string, releaseVersion string, dependencySpecifierID int) error
	Create(productSlug string, releaseVersion string, dependentProductSlug string, specifier string) error
	Delete(productSlug string, releaseVersion string, dependencyspecifierID int) error
}

var NewDependencySpecifierClient = func(client dependencyspecifier.PivnetClient) DependencySpecifierClient {
	return dependencyspecifier.NewDependencySpecifierClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *DependencySpecifiersCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewDependencySpecifierClient(client).List(command.ProductSlug, command.ReleaseVersion)
}

func (command *DependencySpecifierCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewDependencySpecifierClient(client).Get(
		command.ProductSlug,
		command.ReleaseVersion,
		command.DependencySpecifierID,
	)
}

func (command *CreateDependencySpecifierCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewDependencySpecifierClient(client).Create(
		command.ProductSlug,
		command.ReleaseVersion,
		command.DependentProductSlug,
		command.Specifier,
	)
}

func (command *DeleteDependencySpecifierCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewDependencySpecifierClient(client).Delete(
		command.ProductSlug,
		command.ReleaseVersion,
		command.DependencySpecifierID,
	)
}
