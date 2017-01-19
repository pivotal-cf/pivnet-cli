package commands

import "github.com/pivotal-cf/pivnet-cli/commands/release"

type ReleasesCommand struct {
	ProductSlug string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
}

type ReleaseCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type UpdateReleaseCommand struct {
	ProductSlug    string  `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string  `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	Availability   *string `long:"availability" description:"Release availability. Optional." choice:"admins" choice:"selected-user-groups" choice:"all"`
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
	Update(productSlug string, releaseVersion string, availability *string) error
	Delete(productSlug string, releaseVersion string) error
}

var NewReleaseClient = func(client release.PivnetClient) ReleaseClient {
	return release.NewReleaseClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *ReleasesCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseClient(client).List(command.ProductSlug)
}

func (command *ReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseClient(client).Get(command.ProductSlug, command.ReleaseVersion)
}

func (command *CreateReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseClient(client).Create(
		command.ProductSlug,
		command.ReleaseVersion,
		command.ReleaseType,
		command.EULASlug,
	)
}

func (command *UpdateReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseClient(client).Update(
		command.ProductSlug,
		command.ReleaseVersion,
		command.Availability,
	)
}

func (command *DeleteReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewReleaseClient(client).Delete(command.ProductSlug, command.ReleaseVersion)
}
