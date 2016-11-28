package commands

import (
	"errors"

	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-cli/commands/productfile"
)

type ProductFilesCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1"`
}

type ProductFileCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1"`
	ProductFileID  int    `long:"product-file-id" short:"i" description:"Product file ID e.g. 1234" required:"true"`
}

type CreateProductFileCommand struct {
	ProductSlug        string   `long:"product-slug" short:"p" description:"Product slug e.g. 'p-mysql'" required:"true"`
	Name               string   `long:"name" description:"Name e.g. 'p-mysql 1.7.13'" required:"true"`
	AWSObjectKey       string   `long:"aws-object-key" description:"AWS Object Key e.g. 'product_files/P-MySQL/p-mysql-1.7.13.pivotal'" required:"true"`
	FileType           string   `long:"file-type" description:"File Type e.g. 'Software'" required:"true"`
	FileVersion        string   `long:"file-version" description:"File Version e.g. '1.7.13'" required:"true"`
	MD5                string   `long:"md5" description:"MD5 of file" required:"true"`
	Description        string   `long:"description" description:"Description of file"`
	DocsURL            string   `long:"docs-url" description:"URL of docs for file"`
	IncludedFiles      []string `long:"included-file" description:"Name of included file"`
	Platforms          []string `long:"platform" description:"Platform of file"`
	ReleasedAt         string   `long:"released-at" description:"When file is marked for release e.g. '2016/01/16'"`
	SystemRequirements []string `long:"system-requirement" description:"System-requirement of file"`
}

type UpdateProductFileCommand struct {
	ProductSlug   string  `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ProductFileID int     `long:"product-file-id" short:"i" description:"Product file ID e.g. 1234" required:"true"`
	Name          *string `long:"name" description:"Name e.g. p-mysql 1.7.13"`
	FileType      *string `long:"file-type" description:"File Type e.g. 'Software'"`
	FileVersion   *string `long:"file-version" description:"File Version e.g. '1.7.13'"`
	MD5           *string `long:"md5" description:"MD5 of file"`
	Description   *string `long:"description" description:"File description e.g. 'This is a file description.'"`
}

type AddProductFileCommand struct {
	ProductSlug    string  `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion *string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1"`
	ProductFileID  int     `long:"product-file-id" short:"i" description:"Product file ID e.g. 1234" required:"true"`
	FileGroupID    *int    `long:"file-group-id" short:"f" description:"File group ID e.g. 1234"`
}

type RemoveProductFileCommand struct {
	ProductSlug    string  `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion *string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1"`
	ProductFileID  int     `long:"product-file-id" short:"i" description:"Product file ID e.g. 1234" required:"true"`
	FileGroupID    *int    `long:"file-group-id" short:"f" description:"File group ID e.g. 1234"`
}

type DeleteProductFileCommand struct {
	ProductSlug   string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ProductFileID int    `long:"product-file-id" short:"i" description:"Product file ID e.g. 1234" required:"true"`
}

type DownloadProductFilesCommand struct {
	ProductSlug    string   `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string   `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	ProductFileIDs []int    `long:"product-file-id" short:"i" description:"Product file ID e.g. 1234"`
	Globs          []string `long:"glob" short:"g" description:"Glob to match product name e.g. *aws*"`
	DownloadDir    string   `long:"download-dir" short:"d" default:"." description:"Local directory to download files to e.g. /tmp/my-file/"`
	AcceptEULA     bool     `long:"accept-eula" description:"Automatically accept EULA if necessary"`
}

//go:generate counterfeiter . ProductFileClient
type ProductFileClient interface {
	List(productSlug string, releaseVersion string) error
	Get(productSlug string, releaseVersion string, productFileID int) error
	Create(config pivnet.CreateProductFileConfig) error
	Update(
		productFileID int,
		productSlug string,
		name *string,
		fileType *string,
		fileVersion *string,
		md5 *string,
		description *string,
	) error
	AddToRelease(productSlug string, releaseVersion string, productFileID int) error
	RemoveFromRelease(productSlug string, releaseVersion string, productFileID int) error
	AddToFileGroup(productSlug string, fileGroupID int, productFileID int) error
	RemoveFromFileGroup(productSlug string, fileGroupID int, productFileID int) error
	Delete(productSlug string, productFileID int) error
	Download(productSlug string, releaseVersion string, globs []string, productFileIDs []int, downloadDir string, acceptEULA bool) error
}

var NewProductFileClient = func(client productfile.PivnetClient) ProductFileClient {
	return productfile.NewProductFileClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		LogWriter,
		Printer,
		Pivnet.Logger,
		Filter,
	)
}

func (command *ProductFilesCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewProductFileClient(client).List(command.ProductSlug, command.ReleaseVersion)
}

func (command *ProductFileCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewProductFileClient(client).Get(command.ProductSlug, command.ReleaseVersion, command.ProductFileID)
}

func (command *CreateProductFileCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	config := pivnet.CreateProductFileConfig{
		ProductSlug:        command.ProductSlug,
		AWSObjectKey:       command.AWSObjectKey,
		Description:        command.Description,
		DocsURL:            command.DocsURL,
		FileType:           command.FileType,
		FileVersion:        command.FileVersion,
		IncludedFiles:      command.IncludedFiles,
		MD5:                command.MD5,
		Name:               command.Name,
		Platforms:          command.Platforms,
		ReleasedAt:         command.ReleasedAt,
		SystemRequirements: command.SystemRequirements,
	}

	return NewProductFileClient(client).Create(config)
}

func (command *UpdateProductFileCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewProductFileClient(client).Update(
		command.ProductFileID,
		command.ProductSlug,
		command.Name,
		command.FileType,
		command.FileVersion,
		command.MD5,
		command.Description,
	)
}

func (command *AddProductFileCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	if command.ReleaseVersion == nil && command.FileGroupID == nil {
		return errors.New("one of release-version or file-group-id must be provided")
	}
	if command.ReleaseVersion != nil && command.FileGroupID != nil {
		return errors.New("only one of release-version or file-group-id must be provided")
	}

	if command.ReleaseVersion != nil {
		return NewProductFileClient(client).AddToRelease(command.ProductSlug, *command.ReleaseVersion, command.ProductFileID)
	}

	return NewProductFileClient(client).AddToFileGroup(command.ProductSlug, *command.FileGroupID, command.ProductFileID)
}

func (command *RemoveProductFileCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	if command.ReleaseVersion == nil && command.FileGroupID == nil {
		return errors.New("one of release-version or file-group-id must be provided")
	}
	if command.ReleaseVersion != nil && command.FileGroupID != nil {
		return errors.New("only one of release-version or file-group-id must be provided")
	}

	if command.ReleaseVersion != nil {
		return NewProductFileClient(client).RemoveFromRelease(command.ProductSlug, *command.ReleaseVersion, command.ProductFileID)
	}

	return NewProductFileClient(client).RemoveFromFileGroup(command.ProductSlug, *command.FileGroupID, command.ProductFileID)
}

func (command *DeleteProductFileCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewProductFileClient(client).Delete(command.ProductSlug, command.ProductFileID)
}

func (command *DownloadProductFilesCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewProductFileClient(client).Download(
		command.ProductSlug,
		command.ReleaseVersion,
		command.Globs,
		command.ProductFileIDs,
		command.DownloadDir,
		command.AcceptEULA,
	)
}
