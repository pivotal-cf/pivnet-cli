package eula

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
	AcceptEULA(productSlug string, releaseID int) error
	EULAs() ([]pivnet.EULA, error)
	EULA(eulaSlug string) (pivnet.EULA, error)
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
}

type EULAClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewEULAClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *EULAClient {
	return &EULAClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *EULAClient) List() error {
	eulas, err := c.pivnetClient.EULAs()
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printEULAs(eulas)
}

func (c *EULAClient) printEULA(eula pivnet.EULA) error {
	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{"ID", "Slug", "Name"})

		eulaAsString := []string{
			strconv.Itoa(eula.ID), eula.Slug, eula.Name,
		}
		table.Append(eulaAsString)
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(eula)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(eula)
	}

	return nil
}

func (c *EULAClient) Get(eulaSlug string) error {
	eula, err := c.pivnetClient.EULA(eulaSlug)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printEULA(eula)
}

func (c *EULAClient) printEULAs(eulas []pivnet.EULA) error {
	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{"ID", "Slug", "Name"})

		for _, e := range eulas {
			eulaAsString := []string{
				strconv.Itoa(e.ID), e.Slug, e.Name,
			}
			table.Append(eulaAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(eulas)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(eulas)
	}

	return nil
}

func (c *EULAClient) AcceptEULA(productSlug string, releaseVersion string) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.AcceptEULA(productSlug, release.ID)

	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"eula acccepted successfully for %s/%s\n",
			productSlug,
			releaseVersion,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
