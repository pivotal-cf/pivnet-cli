package pivnetversions

import (
	"fmt"
	"github.com/pivotal-cf/pivnet-cli/version"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	PivnetVersions() (pivnet.PivnetVersions, error)
}

type PivnetVersionsClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewPivnetVersionsClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *PivnetVersionsClient {
	return &PivnetVersionsClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *PivnetVersionsClient) List() error {
	pivnetVersions, err := c.pivnetClient.PivnetVersions()
	if err != nil {
		return c.eh.HandleError(err)
	}

	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{"Pivnet Component", "Latest Release Version"})
		table.Append([]string{"Pivnet CLI", pivnetVersions.PivnetCliVersion})
		table.Append([]string{"Pivnet Resource", pivnetVersions.PivnetResourceVersion})
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(pivnetVersions)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(pivnetVersions)
	}

	return nil
}

func (c *PivnetVersionsClient) Warn(currentVersion string) string {
	pivnetVersions, err := c.pivnetClient.PivnetVersions()
	if err == nil && currentVersion != pivnetVersions.PivnetCliVersion {
		return fmt.Sprintf("Warning: Your version of Pivnet CLI (%s) does not match the currently released version (%s).", version.Version, pivnetVersions.PivnetCliVersion)
	} else {
		return ""
	}
}
