package login

import (
	"fmt"
	"io"

	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
	"github.com/pivotal-cf/pivnet-cli/ui"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	Auth() (bool, error)
}

//go:generate counterfeiter . RCHandler
type RCHandler interface {
	SaveProfile(profileName string, apiToken string, host string) error
}

type LoginClient struct {
	pivnetClient PivnetClient
	rcHandler    RCHandler
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewLoginClient(
	pivnetClient PivnetClient,
	rcHandler RCHandler,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *LoginClient {
	return &LoginClient{
		pivnetClient: pivnetClient,
		rcHandler:    rcHandler,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *LoginClient) Login(
	profileName string,
	apiToken string,
	host string,
) error {
	ok, err := c.pivnetClient.Auth()
	if err != nil {
		return c.eh.HandleError(err)
	}

	if !ok {
		err := fmt.Errorf("Failed to login")
		return c.eh.HandleError(err)
	}

	err = c.rcHandler.SaveProfile(profileName, apiToken, host)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printLogin()
}

func (c *LoginClient) printLogin() error {
	switch c.format {

	case printer.PrintAsTable:
		message := "Logged-in successfully"
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err := fmt.Fprintln(c.outputWriter, coloredMessage)

		return err
	}

	return nil
}
