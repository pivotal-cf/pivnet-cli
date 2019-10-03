package companygroup

import (
	"github.com/olekukonko/tablewriter"
	pivnet "github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
	"io"
	"strconv"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	CompanyGroups() ([]pivnet.CompanyGroup, error)
}

type CompanyGroupClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewCompanyGroupClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *CompanyGroupClient {
	return &CompanyGroupClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *CompanyGroupClient) List() error {
	companyGroups, err := c.pivnetClient.CompanyGroups()
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printCompanyGroups(companyGroups)
}

func (c *CompanyGroupClient) printCompanyGroups(companyGroups []pivnet.CompanyGroup) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{"ID", "Name"})

		for _, companyGroup := range companyGroups {
			companyGroupAsString := []string{
				strconv.Itoa(companyGroup.ID),
				companyGroup.Name,
			}
			table.Append(companyGroupAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(companyGroups)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(companyGroups)
	}

	return nil
}
