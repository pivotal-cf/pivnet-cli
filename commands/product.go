package commands

import "github.com/pivotal-cf/pivnet-cli/commands/product"

type ProductCommand struct {
	ProductSlug string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
}

type ProductsCommand struct {
}

type ProductSlugsCommand struct {
	ProductSlug string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
}

//go:generate counterfeiter . ProductClient
type ProductClient interface {
	List() error
	Get(productSlug string) error
	SlugAlias(productSlug string) error
}

var NewProductClient = func(client product.PivnetClient) ProductClient {
	return product.NewProductClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *ProductsCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewProductClient(client).List()
}

func (command *ProductCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewProductClient(client).Get(command.ProductSlug)
}

func (command *ProductSlugsCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewProductClient(client).SlugAlias(command.ProductSlug)
}
