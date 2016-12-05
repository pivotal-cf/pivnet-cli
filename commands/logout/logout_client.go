package logout

import (
	"fmt"
	"io"

	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

//go:generate counterfeiter . RCHandler
type RCHandler interface {
	SaveProfile(profileName string, apiToken string, host string) error
}

type LogoutClient struct {
	rcHandler    RCHandler
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewLogoutClient(
	rcHandler RCHandler,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *LogoutClient {
	return &LogoutClient{
		rcHandler:    rcHandler,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *LogoutClient) Logout(profileName string) error {
	apiToken := ""
	host := ""

	err := c.rcHandler.SaveProfile(profileName, apiToken, host)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printLogout()
}

func (c *LogoutClient) printLogout() error {
	switch c.format {

	case printer.PrintAsTable:
		fmt.Fprintln(c.outputWriter, "logged-out successfully")
		return nil
	}

	return nil
}
