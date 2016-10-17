package commands

import (
	"errors"

	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/commands/productfile"
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
	ProductSlug  string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	Name         string `long:"name" description:"Name e.g. p-mysql 1.7.13" required:"true"`
	AWSObjectKey string `long:"aws-object-key" description:"AWS Object Key e.g. product_files/P-MySQL/p-mysql-1.7.13.pivotal" required:"true"`
	FileType     string `long:"file-type" description:"File Type e.g. 'Software'" required:"true"`
	FileVersion  string `long:"file-version" description:"File Version e.g. '1.7.13'" required:"true"`
	MD5          string `long:"md5" description:"MD5 of file" required:"true"`
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

type DownloadProductFileCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	ProductFileID  int    `long:"product-file-id" short:"i" description:"Product file ID e.g. 1234" required:"true"`
	Filepath       string `long:"filepath" description:"Local filepath to download file to e.g. /tmp/my-file" required:"true"`
	AcceptEULA     bool   `long:"accept-eula" description:"Automatically accept EULA if necessary"`
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
	Download(productSlug string, releaseVersion string, productFileID int, filepath string, acceptEULA bool) error
}

var NewProductFileClient = func() ProductFileClient {
	return productfile.NewProductFileClient(
		NewPivnetClient(),
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		LogWriter,
		Printer,
		Pivnet.Logger,
	)
}

func (command *ProductFilesCommand) Execute([]string) error {
	Init()
	return NewProductFileClient().List(command.ProductSlug, command.ReleaseVersion)
}

func (command *ProductFileCommand) Execute([]string) error {
	Init()
	return NewProductFileClient().Get(command.ProductSlug, command.ReleaseVersion, command.ProductFileID)
}

func (command *CreateProductFileCommand) Execute([]string) error {
	Init()

	config := pivnet.CreateProductFileConfig{
		ProductSlug:  command.ProductSlug,
		Name:         command.Name,
		AWSObjectKey: command.AWSObjectKey,
		FileType:     command.FileType,
		FileVersion:  command.FileVersion,
		MD5:          command.MD5,
	}

	return NewProductFileClient().Create(config)
}

func (command *UpdateProductFileCommand) Execute([]string) error {
	Init()

	return NewProductFileClient().Update(
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
	Init()

	if command.ReleaseVersion == nil && command.FileGroupID == nil {
		return errors.New("one of release-version or file-group-id must be provided")
	}
	if command.ReleaseVersion != nil && command.FileGroupID != nil {
		return errors.New("only one of release-version or file-group-id must be provided")
	}

	if command.ReleaseVersion != nil {
		return NewProductFileClient().AddToRelease(command.ProductSlug, *command.ReleaseVersion, command.ProductFileID)
	}

	return NewProductFileClient().AddToFileGroup(command.ProductSlug, *command.FileGroupID, command.ProductFileID)
}

func (command *RemoveProductFileCommand) Execute([]string) error {
	Init()

	if command.ReleaseVersion == nil && command.FileGroupID == nil {
		return errors.New("one of release-version or file-group-id must be provided")
	}
	if command.ReleaseVersion != nil && command.FileGroupID != nil {
		return errors.New("only one of release-version or file-group-id must be provided")
	}

	if command.ReleaseVersion != nil {
		return NewProductFileClient().RemoveFromRelease(command.ProductSlug, *command.ReleaseVersion, command.ProductFileID)
	}

	return NewProductFileClient().RemoveFromFileGroup(command.ProductSlug, *command.FileGroupID, command.ProductFileID)
}

func (command *DeleteProductFileCommand) Execute([]string) error {
	Init()
	return NewProductFileClient().Delete(command.ProductSlug, command.ProductFileID)
}

func (command *DownloadProductFileCommand) Execute([]string) error {
	Init()

	return NewProductFileClient().Download(
		command.ProductSlug,
		command.ReleaseVersion,
		command.ProductFileID,
		command.Filepath,
		command.AcceptEULA,
	)
}
