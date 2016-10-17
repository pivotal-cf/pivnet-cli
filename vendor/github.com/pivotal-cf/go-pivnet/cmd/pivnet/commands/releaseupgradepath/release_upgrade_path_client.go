package releaseupgradepath

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
	ReleaseUpgradePaths(productSlug string, releaseID int) ([]pivnet.ReleaseUpgradePath, error)
	AddReleaseUpgradePath(productSlug string, releaseID int, previousReleaseID int) error
	RemoveReleaseUpgradePath(productSlug string, releaseID int, previousReleaseID int) error
}

type ReleaseUpgradePathClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewReleaseUpgradePathClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *ReleaseUpgradePathClient {
	return &ReleaseUpgradePathClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *ReleaseUpgradePathClient) List(productSlug string, releaseVersion string) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	releaseUpgradePaths, err := c.pivnetClient.ReleaseUpgradePaths(
		productSlug,
		release.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printReleaseUpgradePaths(releaseUpgradePaths)
}

func (c *ReleaseUpgradePathClient) printReleaseUpgradePaths(releaseUpgradePaths []pivnet.ReleaseUpgradePath) error {
	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Version",
		})

		for _, r := range releaseUpgradePaths {
			table.Append([]string{
				strconv.Itoa(r.Release.ID),
				r.Release.Version,
			})
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(releaseUpgradePaths)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(releaseUpgradePaths)
	}

	return nil
}

func (c *ReleaseUpgradePathClient) Add(
	productSlug string,
	releaseVersion string,
	previousReleaseVersion string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	previousRelease, err := c.pivnetClient.ReleaseForVersion(productSlug, previousReleaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.AddReleaseUpgradePath(
		productSlug,
		release.ID,
		previousRelease.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"release upgrade path added successfully to %s/%s\n",
			productSlug,
			releaseVersion,
		)
	}

	return nil
}

func (c *ReleaseUpgradePathClient) Remove(
	productSlug string,
	releaseVersion string,
	previousReleaseVersion string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	previousRelease, err := c.pivnetClient.ReleaseForVersion(productSlug, previousReleaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.RemoveReleaseUpgradePath(
		productSlug,
		release.ID,
		previousRelease.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"release upgrade path removed successfully from %s/%s\n",
			productSlug,
			releaseVersion,
		)
	}

	return nil
}
