package commands

import (
	pivnet "github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/pivnet-cli/commands/imagereference"
)

type ImageReferencesCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1 (Required for non-admins)"`
}

type CreateImageReferenceCommand struct {
	ProductSlug        string   `long:"product-slug" short:"p" description:"Product slug e.g. 'p-mysql'" required:"true"`
	Name               string   `long:"name" description:"Name e.g. 'p-mysql 1.7.13'" required:"true"`
	ImagePath          string   `long:"image-path" description:"Image path e.g. 'library/image:tag'" required:"true"`
	Digest             string   `long:"digest" description:"Digest of the file e.g. 'sha256:1234abcd'" required:"true"`
	Description        string   `long:"description" description:"Description of file"`
	DocsURL            string   `long:"docs-url" description:"URL of docs for file"`
	SystemRequirements []string `long:"system-requirement" description:"System-requirement of file"`
}

//go:generate counterfeiter . ImageReferenceClient
type ImageReferenceClient interface {
	List(productSlug string, releaseVersion string) error
	Create(config pivnet.CreateImageReferenceConfig) error
}

var NewImageReferenceClient = func(client imagereference.PivnetClient) ImageReferenceClient {
	return imagereference.NewImageReferenceClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		LogWriter,
		Printer,
		Pivnet.Logger,
	)
}

func (command *ImageReferencesCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewImageReferenceClient(client).List(command.ProductSlug, command.ReleaseVersion)
}


func (command *CreateImageReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	config := pivnet.CreateImageReferenceConfig{
		ProductSlug:        command.ProductSlug,
		Name:               command.Name,
		ImagePath:          command.ImagePath,
		Digest:             command.Digest,
		Description:        command.Description,
		DocsURL:            command.DocsURL,
		SystemRequirements: command.SystemRequirements,
	}

	return NewImageReferenceClient(client).Create(config)
}
