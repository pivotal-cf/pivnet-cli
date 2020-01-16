package commands

import (
	"github.com/pivotal-cf/go-pivnet/v3"
	"github.com/pivotal-cf/pivnet-cli/commands/helmchartreference"
)

type HelmChartReferencesCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1 (Required for non-admins)"`
}

type HelmChartReferenceCommand struct {
	ProductSlug          string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion       string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1"`
	HelmChartReferenceID int    `long:"helm-chart-reference-id" short:"i" description:"Helm chart reference ID e.g. 1234" required:"true"`
}

type CreateHelmChartReferenceCommand struct {
	ProductSlug        string   `long:"product-slug" short:"p" description:"Product slug e.g. 'p-mysql'" required:"true"`
	Name               string   `long:"name" description:"Name e.g. 'p-mysql'" required:"true"`
	Version            string   `long:"version" description:"Version path e.g. '1.2.3'" required:"true"`
	Description        string   `long:"description" description:"Description of the helm chart"`
	DocsURL            string   `long:"docs-url" description:"URL of docs for the helm chart"`
	SystemRequirements []string `long:"system-requirement" description:"System-requirement of the helm chart"`
}

type UpdateHelmChartReferenceCommand struct {
	ProductSlug          string    `long:"product-slug" short:"p" description:"Product slug e.g. 'p-mysql'" required:"true"`
	HelmChartReferenceID int       `long:"helm-chart-reference-id" short:"i" description:"Helm chart reference ID e.g. 1234" required:"true"`
	Description          *string   `long:"description" description:"Description of the helm chart"`
	DocsURL              *string   `long:"docs-url" description:"URL of docs for the helm chart"`
	SystemRequirements   *[]string `long:"system-requirement" description:"System-requirement of the helm chart"`
}

type DeleteHelmChartReferenceCommand struct {
	ProductSlug          string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	HelmChartReferenceID int    `long:"helm-chart-reference-id" short:"i" description:"Helm chart reference ID e.g. 1234" required:"true"`
}

type AddHelmChartReferenceToReleaseCommand struct {
	ProductSlug          string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	HelmChartReferenceID int    `long:"helm-chart-reference-id" short:"i" description:"Helm chart reference ID e.g. 1234" required:"true"`
	ReleaseVersion       string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type RemoveHelmChartReferenceFromReleaseCommand struct {
	ProductSlug          string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	HelmChartReferenceID int    `long:"helm-chart-reference-id" short:"i" description:"Helm chart reference ID e.g. 1234" required:"true"`
	ReleaseVersion       string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

//go:generate counterfeiter . HelmChartReferenceClient
type HelmChartReferenceClient interface {
	List(productSlug string, releaseVersion string) error
	Get(productSlug string, releaseVersion string, helmChartReferenceID int) error
	Create(config pivnet.CreateHelmChartReferenceConfig) error
	Update(
		productSlug string,
		helmChartReferenceID int,
		description *string,
		docsURL *string,
		systemRequirements *[]string,
	) error
	Delete(productSlug string, helmChartReferenceID int) error
	AddToRelease(productSlug string, helmChartReferenceID int, releaseVersion string) error
	RemoveFromRelease(productSlug string, helmChartReferenceID int, releaseVersion string) error
}

var NewHelmChartReferenceClient = func(client helmchartreference.PivnetClient) HelmChartReferenceClient {
	return helmchartreference.NewHelmChartReferenceClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		LogWriter,
		Printer,
		Pivnet.Logger,
	)
}

func (command *HelmChartReferencesCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewHelmChartReferenceClient(client).List(command.ProductSlug, command.ReleaseVersion)
}

func (command *HelmChartReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewHelmChartReferenceClient(client).Get(command.ProductSlug, command.ReleaseVersion, command.HelmChartReferenceID)
}

func (command *CreateHelmChartReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	config := pivnet.CreateHelmChartReferenceConfig{
		ProductSlug:        command.ProductSlug,
		Name:               command.Name,
		Version:            command.Version,
		Description:        command.Description,
		DocsURL:            command.DocsURL,
		SystemRequirements: command.SystemRequirements,
	}

	return NewHelmChartReferenceClient(client).Create(config)
}

func (command *DeleteHelmChartReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewHelmChartReferenceClient(client).Delete(command.ProductSlug, command.HelmChartReferenceID)
}

func (command *AddHelmChartReferenceToReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewHelmChartReferenceClient(client).AddToRelease(
		command.ProductSlug,
		command.HelmChartReferenceID,
		command.ReleaseVersion,
	)
}

func (command *RemoveHelmChartReferenceFromReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewHelmChartReferenceClient(client).RemoveFromRelease(
		command.ProductSlug,
		command.HelmChartReferenceID,
		command.ReleaseVersion,
	)
}

func (command *UpdateHelmChartReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewHelmChartReferenceClient(client).Update(
		command.ProductSlug,
		command.HelmChartReferenceID,
		command.Description,
		command.DocsURL,
		command.SystemRequirements,
	)
}
