package commands

import (
	"github.com/pivotal-cf/go-pivnet/v7"
	"github.com/pivotal-cf/pivnet-cli/v2/commands/artifactreference"
)

type ArtifactReferencesCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1 (Required for non-admins)"`
	ArtifactDigest string `long:"digest" short:"d" description:"Artifact digest e.g. sha256:9f86d0... (if provided ignores release-version value)"`
}

type ArtifactReferenceCommand struct {
	ProductSlug         string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion      string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1"`
	ArtifactReferenceID int    `long:"artifact-reference-id" short:"i" description:"Artifact reference ID e.g. 1234" required:"true"`
}

type CreateArtifactReferenceCommand struct {
	ProductSlug        string   `long:"product-slug" short:"p" description:"Product slug e.g. 'p-mysql'" required:"true"`
	Name               string   `long:"name" description:"Name e.g. 'p-mysql 1.7.13'" required:"true"`
	ArtifactPath       string   `long:"artifact-path" description:"Artifact path e.g. 'repo_name:tag'" required:"true"`
	Digest             string   `long:"digest" description:"Digest of the artifact e.g. 'sha256:1234abcd'" required:"true"`
	Description        string   `long:"description" description:"Description of the artifact"`
	DocsURL            string   `long:"docs-url" description:"URL of docs for the artifact"`
	SystemRequirements []string `long:"system-requirement" description:"System-requirement of the artifact"`
}

type UpdateArtifactReferenceCommand struct {
	ProductSlug         string    `long:"product-slug" short:"p" description:"Product slug e.g. 'p-mysql'" required:"true"`
	ArtifactReferenceID int       `long:"artifact-reference-id" short:"i" description:"Artifact reference ID e.g. 1234" required:"true"`
	Name                *string   `long:"name" description:"Name e.g. 'p-mysql 1.7.13'" `
	Description         *string   `long:"description" description:"Description of the artifact"`
	DocsURL             *string   `long:"docs-url" description:"URL of docs for the artifact"`
	SystemRequirements  *[]string `long:"system-requirement" description:"System-requirement of the artifact"`
}

type DeleteArtifactReferenceCommand struct {
	ProductSlug         string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ArtifactReferenceID int    `long:"artifact-reference-id" short:"i" description:"Artifact reference ID e.g. 1234" required:"true"`
}

type AddArtifactReferenceToReleaseCommand struct {
	ProductSlug         string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ArtifactReferenceID int    `long:"artifact-reference-id" short:"i" description:"Artifact reference ID e.g. 1234" required:"true"`
	ReleaseVersion      string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type RemoveArtifactReferenceFromReleaseCommand struct {
	ProductSlug         string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ArtifactReferenceID int    `long:"artifact-reference-id" short:"i" description:"Artifact reference ID e.g. 1234" required:"true"`
	ReleaseVersion      string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

//go:generate counterfeiter . ArtifactReferenceClient
type ArtifactReferenceClient interface {
	List(productSlug string, releaseVersion string, artifactDigest string) error
	Get(productSlug string, releaseVersion string, artifactReferenceID int) error
	Create(config pivnet.CreateArtifactReferenceConfig) error
	Update(
		productSlug string,
		artifactReferenceID int,
		name *string,
		description *string,
		docsURL *string,
		systemRequirements *[]string,
	) error
	Delete(productSlug string, artifactReferenceID int) error
	AddToRelease(productSlug string, artifactReferenceID int, releaseVersion string) error
	RemoveFromRelease(productSlug string, artifactReferenceID int, releaseVersion string) error
}

var NewArtifactReferenceClient = func(client artifactreference.PivnetClient) ArtifactReferenceClient {
	return artifactreference.NewArtifactReferenceClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		LogWriter,
		Printer,
		Pivnet.Logger,
	)
}

func (command *ArtifactReferencesCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewArtifactReferenceClient(client).List(command.ProductSlug, command.ReleaseVersion, command.ArtifactDigest)
}

func (command *ArtifactReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewArtifactReferenceClient(client).Get(command.ProductSlug, command.ReleaseVersion, command.ArtifactReferenceID)
}

func (command *CreateArtifactReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	config := pivnet.CreateArtifactReferenceConfig{
		ProductSlug:        command.ProductSlug,
		Name:               command.Name,
		ArtifactPath:       command.ArtifactPath,
		Digest:             command.Digest,
		Description:        command.Description,
		DocsURL:            command.DocsURL,
		SystemRequirements: command.SystemRequirements,
	}

	return NewArtifactReferenceClient(client).Create(config)
}

func (command *DeleteArtifactReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewArtifactReferenceClient(client).Delete(command.ProductSlug, command.ArtifactReferenceID)
}

func (command *AddArtifactReferenceToReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewArtifactReferenceClient(client).AddToRelease(
		command.ProductSlug,
		command.ArtifactReferenceID,
		command.ReleaseVersion,
	)
}

func (command *RemoveArtifactReferenceFromReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewArtifactReferenceClient(client).RemoveFromRelease(
		command.ProductSlug,
		command.ArtifactReferenceID,
		command.ReleaseVersion,
	)
}

func (command *UpdateArtifactReferenceCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewArtifactReferenceClient(client).Update(
		command.ProductSlug,
		command.ArtifactReferenceID,
		command.Name,
		command.Description,
		command.DocsURL,
		command.SystemRequirements,
	)
}
