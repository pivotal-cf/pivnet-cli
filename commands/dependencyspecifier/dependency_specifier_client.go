package dependencyspecifier

import (
	"fmt"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
	"github.com/pivotal-cf/pivnet-cli/ui"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
	DependencySpecifiers(productSlug string, releaseID int) ([]pivnet.DependencySpecifier, error)
	DependencySpecifier(productSlug string, releaseID int, dependencySpecifierID int) (pivnet.DependencySpecifier, error)
	CreateDependencySpecifier(productSlug string, releaseID int, dependentProductSlug string, specifier string) (pivnet.DependencySpecifier, error)
	DeleteDependencySpecifier(productSlug string, releaseID int, dependencySpecifierID int) error
}

type DependencySpecifierClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewDependencySpecifierClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *DependencySpecifierClient {
	return &DependencySpecifierClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *DependencySpecifierClient) List(productSlug string, releaseVersion string) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	dependencySpecifiers, err := c.pivnetClient.DependencySpecifiers(
		productSlug,
		release.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printDependencySpecifiers(dependencySpecifiers)
}

func (c *DependencySpecifierClient) Get(
	productSlug string,
	releaseVersion string,
	dependencySpecifierID int,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	dependencySpecifier, err := c.pivnetClient.DependencySpecifier(
		productSlug,
		release.ID,
		dependencySpecifierID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printDependencySpecifier(dependencySpecifier)
}

func (c *DependencySpecifierClient) printDependencySpecifier(dependencySpecifier pivnet.DependencySpecifier) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Specifier",
			"Product Name",
			"Product ID",
			"Product Slug",
		})

		table.Append([]string{
			strconv.Itoa(dependencySpecifier.ID),
			dependencySpecifier.Specifier,
			dependencySpecifier.Product.Name,
			strconv.Itoa(dependencySpecifier.Product.ID),
			dependencySpecifier.Product.Slug,
		})
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(dependencySpecifier)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(dependencySpecifier)
	}

	return nil
}

func (c *DependencySpecifierClient) printDependencySpecifiers(dependencySpecifiers []pivnet.DependencySpecifier) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Specifier",
			"Product Name",
			"Product ID",
			"Product Slug",
		})

		for _, dependencySpecifier := range dependencySpecifiers {
			table.Append([]string{
				strconv.Itoa(dependencySpecifier.ID),
				dependencySpecifier.Specifier,
				dependencySpecifier.Product.Name,
				strconv.Itoa(dependencySpecifier.Product.ID),
				dependencySpecifier.Product.Slug,
			})
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(dependencySpecifiers)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(dependencySpecifiers)
	}

	return nil
}

func (c *DependencySpecifierClient) Create(
	productSlug string,
	releaseVersion string,
	dependentProductSlug string,
	specifier string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	dependencySpecifier, err := c.pivnetClient.CreateDependencySpecifier(
		productSlug,
		release.ID,
		dependentProductSlug,
		specifier,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printDependencySpecifier(dependencySpecifier)
}

func (c *DependencySpecifierClient) Delete(
	productSlug string,
	releaseVersion string,
	dependencySpecifierID int,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.DeleteDependencySpecifier(
		productSlug,
		release.ID,
		dependencySpecifierID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		message := fmt.Sprintf(
			"Dependency specifier '%d' deleted successfully",
			dependencySpecifierID,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err := fmt.Fprintln(c.outputWriter, coloredMessage)

		return err
	}

	return nil
}
