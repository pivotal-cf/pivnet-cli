package product

import (
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/errorhandler"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	Products() ([]pivnet.Product, error)
	FindProductForSlug(productSlug string) (pivnet.Product, error)
}

type ProductClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewProductClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *ProductClient {
	return &ProductClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *ProductClient) List() error {
	products, err := c.pivnetClient.Products()
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printProducts(products)
}

func (c *ProductClient) printProducts(products []pivnet.Product) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Slug",
			"Name",
		})

		for _, product := range products {
			productAsString := []string{
				strconv.Itoa(product.ID),
				product.Slug,
				product.Name,
			}
			table.Append(productAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(products)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(products)
	}

	return nil
}

func (c *ProductClient) Get(productSlug string) error {
	product, err := c.pivnetClient.FindProductForSlug(productSlug)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printProduct(product)
}

func (c *ProductClient) printProduct(product pivnet.Product) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Slug",
			"Name",
		})

		productAsString := []string{
			strconv.Itoa(product.ID),
			product.Slug,
			product.Name,
		}
		table.Append(productAsString)
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(product)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(product)
	}

	return nil
}
