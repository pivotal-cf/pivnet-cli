package commands

import "github.com/pivotal-cf/go-pivnet/cmd/pivnet/commands/filegroup"

type FileGroupsCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1"`
}

type FileGroupCommand struct {
	ProductSlug string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	FileGroupID int    `long:"file-group-id" short:"i" description:"Filegroup ID e.g. 1234" required:"true"`
}

type CreateFileGroupCommand struct {
	ProductSlug string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	Name        string `long:"name" description:"Name e.g. my_file_group" required:"true"`
}

type UpdateFileGroupCommand struct {
	ProductSlug string  `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	FileGroupID int     `long:"file-group-id" short:"i" description:"Filegroup ID e.g. 1234" required:"true"`
	Name        *string `long:"name" description:"Name e.g. my_file_group"`
}

type DeleteFileGroupCommand struct {
	ProductSlug string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	FileGroupID int    `long:"file-group-id" short:"i" description:"File group ID e.g. 1234" required:"true"`
}

type AddFileGroupToReleaseCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	FileGroupID    int    `long:"file-group-id" short:"i" description:"Filegroup ID e.g. 1234" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

type RemoveFileGroupFromReleaseCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	FileGroupID    int    `long:"file-group-id" short:"i" description:"Filegroup ID e.g. 1234" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
}

//go:generate counterfeiter . FileGroupClient
type FileGroupClient interface {
	List(productSlug string, releaseVersion string) error
	Get(productSlug string, productFileID int) error
	Create(productSlug string, name string) error
	Update(productSlug string, productFileID int, name *string) error
	Delete(productSlug string, productFileID int) error
	AddToRelease(productSlug string, productFileID int, releaseVersion string) error
	RemoveFromRelease(productSlug string, productFileID int, releaseVersion string) error
}

var NewFileGroupClient = func() FileGroupClient {
	return filegroup.NewFileGroupClient(
		NewPivnetClient(),
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *FileGroupsCommand) Execute([]string) error {
	Init()
	return NewFileGroupClient().List(command.ProductSlug, command.ReleaseVersion)
}

func (command *FileGroupCommand) Execute([]string) error {
	Init()
	return NewFileGroupClient().Get(command.ProductSlug, command.FileGroupID)
}

func (command *CreateFileGroupCommand) Execute([]string) error {
	Init()
	return NewFileGroupClient().Create(command.ProductSlug, command.Name)
}

func (command *UpdateFileGroupCommand) Execute([]string) error {
	Init()

	return NewFileGroupClient().Update(
		command.ProductSlug,
		command.FileGroupID,
		command.Name,
	)
}

func (command *DeleteFileGroupCommand) Execute([]string) error {
	Init()
	return NewFileGroupClient().Delete(command.ProductSlug, command.FileGroupID)
}

func (command *AddFileGroupToReleaseCommand) Execute([]string) error {
	Init()

	return NewFileGroupClient().AddToRelease(
		command.ProductSlug,
		command.FileGroupID,
		command.ReleaseVersion,
	)
}

func (command *RemoveFileGroupFromReleaseCommand) Execute([]string) error {
	Init()

	return NewFileGroupClient().RemoveFromRelease(
		command.ProductSlug,
		command.FileGroupID,
		command.ReleaseVersion,
	)
}
