package product

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf/go-pivnet/v6"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	Products() ([]pivnet.Product, error)
	FindProductForSlug(productSlug string) (pivnet.Product, error)
	SlugAlias(productSlug string) (pivnet.SlugAliasResponse, error)
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

func (c *ProductClient) SlugAlias(productSlug string) error {
	slugAliasResponse, err := c.pivnetClient.SlugAlias(productSlug)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printSlugAlias(slugAliasResponse)
}

func (c *ProductClient) printSlugAlias(response pivnet.SlugAliasResponse) error {
	switch c.format {

	case printer.PrintAsTable:
		b := color.New(color.Bold).SprintFunc()

		fmt.Println("")
		fmt.Println("Current Product Slug: ", b(response.CurrentSlug))
		fmt.Println("")

		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{"Slugs"})

		for _, slug := range response.Slugs {
			table.Append([]string{slug})
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(response)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(response)
	}

	return nil
}
