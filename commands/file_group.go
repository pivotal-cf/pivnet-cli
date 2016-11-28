package commands

import "github.com/pivotal-cf/pivnet-cli/commands/filegroup"

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

var NewFileGroupClient = func(client filegroup.PivnetClient) FileGroupClient {
	return filegroup.NewFileGroupClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *FileGroupsCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewFileGroupClient(client).List(command.ProductSlug, command.ReleaseVersion)
}

func (command *FileGroupCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewFileGroupClient(client).Get(command.ProductSlug, command.FileGroupID)
}

func (command *CreateFileGroupCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewFileGroupClient(client).Create(command.ProductSlug, command.Name)
}

func (command *UpdateFileGroupCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewFileGroupClient(client).Update(
		command.ProductSlug,
		command.FileGroupID,
		command.Name,
	)
}

func (command *DeleteFileGroupCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewFileGroupClient(client).Delete(command.ProductSlug, command.FileGroupID)
}

func (command *AddFileGroupToReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewFileGroupClient(client).AddToRelease(
		command.ProductSlug,
		command.FileGroupID,
		command.ReleaseVersion,
	)
}

func (command *RemoveFileGroupFromReleaseCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewFileGroupClient(client).RemoveFromRelease(
		command.ProductSlug,
		command.FileGroupID,
		command.ReleaseVersion,
	)
}
