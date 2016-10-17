package releasedependency

import (
	"fmt"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/errorhandler"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
	ReleaseDependencies(productSlug string, releaseID int) ([]pivnet.ReleaseDependency, error)
	AddReleaseDependency(productSlug string, releaseID int, dependentReleaseID int) error
	RemoveReleaseDependency(productSlug string, releaseID int, dependentReleaseID int) error
}

type ReleaseDependencyClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewReleaseDependencyClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *ReleaseDependencyClient {
	return &ReleaseDependencyClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *ReleaseDependencyClient) List(productSlug string, releaseVersion string) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	releaseDependencies, err := c.pivnetClient.ReleaseDependencies(
		productSlug,
		release.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Version",
			"Product ID",
			"Product Name",
		})

		for _, r := range releaseDependencies {
			table.Append([]string{
				strconv.Itoa(r.Release.ID),
				r.Release.Version,
				strconv.Itoa(r.Release.Product.ID),
				r.Release.Product.Name,
			})
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(releaseDependencies)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(releaseDependencies)
	}

	return nil
}

func (c *ReleaseDependencyClient) Add(
	productSlug string,
	releaseVersion string,
	dependentProductSlug string,
	dependentReleaseVersion string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	dependentRelease, err := c.pivnetClient.ReleaseForVersion(dependentProductSlug, dependentReleaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.AddReleaseDependency(
		productSlug,
		release.ID,
		dependentRelease.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"release dependency added successfully to %s/%s\n",
			productSlug,
			releaseVersion,
		)
	}

	return nil
}

func (c *ReleaseDependencyClient) Remove(
	productSlug string,
	releaseVersion string,
	dependentProductSlug string,
	dependentReleaseVersion string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	dependentRelease, err := c.pivnetClient.ReleaseForVersion(dependentProductSlug, dependentReleaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.RemoveReleaseDependency(
		productSlug,
		release.ID,
		dependentRelease.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"release dependency removed successfully from %s/%s\n",
			productSlug,
			releaseVersion,
		)
	}

	return nil
}
