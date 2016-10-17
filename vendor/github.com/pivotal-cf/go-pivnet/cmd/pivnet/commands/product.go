package commands

import "github.com/pivotal-cf/go-pivnet/cmd/pivnet/commands/product"

type ProductCommand struct {
	ProductSlug string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
}

type ProductsCommand struct {
}

//go:generate counterfeiter . ProductClient
type ProductClient interface {
	List() error
	Get(productSlug string) error
}

var NewProductClient = func() ProductClient {
	return product.NewProductClient(
		NewPivnetClient(),
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *ProductsCommand) Execute([]string) error {
	Init()

	return NewProductClient().List()
}

func (command *ProductCommand) Execute([]string) error {
	Init()

	return NewProductClient().Get(command.ProductSlug)
}
