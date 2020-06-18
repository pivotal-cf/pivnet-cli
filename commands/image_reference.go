package commands

import (
	"github.com/pivotal-cf/go-pivnet/v5"
	"github.com/pivotal-cf/pivnet-cli/commands/imagereference"
)

type ImageReferencesCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1 (Required for non-admins)"`
	ImageDigest    string `long:"digest" short:"d" description:"Image digest e.g. sha256:9f86d0... (if provided ignores release-version value)"`
}

type ImageReferenceCommand struct {
	ProductSlug      string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion   string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1"`
	ImageReferenceID int    `long:"image-reference-id" short:"i" description:"Image reference ID e.g. 1234" required:"true"`
}

type CreateImageReferenceCommand struct {
	ProductSlug        string   `long:"product-slug" short:"p" description:"Product slug e.g. 'p-mysql'" required:"true"`
	Name               string   `long:"name" description:"Name e.g. 'p-mysql 1.7.13'" required:"true"`
	ImagePath          string   `long:"image-path" description:"Image path e.g. 'library/image:tag'" required:"true"`
	Digest             string   `long:"digest" description:"Digest of the image e.g. 'sha256:1234abcd'" required:"true"`
	Description        string   `long:"description" description:"Description of the image"`
	DocsURL            string   `long:"docs-url" description:"URL of docs for the image"`
	SystemRequirements []string `long:"system-requirement" description:"System-requirement of the image"`
}

type UpdateImageReferenceCommand struct {
	ProductSlug        string    `long:"product-slug" short:"p" description:"Product slug e.g. 'p-mysql'" required:"true"`
	ImageReferenceID   int       `long:"image-reference-id" short:"i" description:"Image reference ID e.g. 1234" required:"true"`
	Name               *string   `long:"name" description:"Name e.g. 'p-mysql 1.7.13'" `
	Description        *string   `long:"description" description:"Description of the image"`
	DocsURL            *string   `long:"docs-url" description:"URL of docs for the image"`
	SystemRequirements *[]string `long:"system-requirement" description:"System-requirement of the image"`
}

type DeleteImageReferenceCommand struct {
	ProductSlug      string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ImageReferenceID int    `long:"image-reference-id" short:"i" description:"Image reference ID e.g. 1234" required:"true"`
}

type AddImageReferenceToReleaseCommand struct {
	ProductSlug      string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ImageReferenceID int    `long:"image-reference-id" short:"i" description:"Image reference ID e.g. 1234" required:"true"`
	ReleaseVersion   string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type RemoveImageReferenceFromReleaseCommand struct {
	ProductSlug      string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ImageReferenceID int    `long:"image-reference-id" short:"i" description:"Image reference ID e.g. 1234" required:"true"`
	ReleaseVersion   string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

//go:generate counterfeiter . ImageReferenceClient
type ImageReferenceClient interface {
	List(productSlug string, releaseVersion string, imageDigest string) error
	Get(productSlug string, releaseVersion string, imageReferenceID int) error
	Create(config pivnet.CreateImageReferenceConfig) error
	Update(
		productSlug string,
		imageReferenceID int,
		name *string,
		description *string,
		docsURL *string,
		systemRequirements *[]string,
	) error
	Delete(productSlug string, imageReferenceID int) error
	AddToRelease(productSlug string, imageReferenceID int, releaseVersion string) error
	RemoveFromRelease(productSlug string, imageReferenceID int, releaseVersion string) error
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

	return NewImageReferenceClient(client).List(command.ProductSlug, command.ReleaseVersion, command.ImageDigest)
}

func (command *ImageReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewImageReferenceClient(client).Get(command.ProductSlug, command.ReleaseVersion, command.ImageReferenceID)
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

func (command *DeleteImageReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewImageReferenceClient(client).Delete(command.ProductSlug, command.ImageReferenceID)
}

func (command *AddImageReferenceToReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewImageReferenceClient(client).AddToRelease(
		command.ProductSlug,
		command.ImageReferenceID,
		command.ReleaseVersion,
	)
}

func (command *RemoveImageReferenceFromReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewImageReferenceClient(client).RemoveFromRelease(
		command.ProductSlug,
		command.ImageReferenceID,
		command.ReleaseVersion,
	)
}

func (command *UpdateImageReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewImageReferenceClient(client).Update(
		command.ProductSlug,
		command.ImageReferenceID,
		command.Name,
		command.Description,
		command.DocsURL,
		command.SystemRequirements,
	)
}
