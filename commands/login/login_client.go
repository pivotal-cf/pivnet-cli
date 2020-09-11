package login

import (
	"fmt"
	"io"

	"github.com/pivotal-cf/pivnet-cli/v2/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/v2/printer"
	"github.com/pivotal-cf/pivnet-cli/v2/ui"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	Auth() (bool, error)
}

//go:generate counterfeiter . RCHandler
type RCHandler interface {
	SaveProfile(profileName string, apiToken string, host string, accessToken string, accessTokenExpiry int64) error
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
	const legacyAPITokenLength = 20
	if len(apiToken) <= legacyAPITokenLength {
		message := "Warning: The use of static Pivnet API tokens is deprecated and will be removed. Please see https://network.pivotal.io/docs/api#how-to-authenticate for details on the new UAA API Token mechanism."
		coloredMessage := ui.ErrorColor.SprintFunc()(message)
		fmt.Fprintln(c.outputWriter, coloredMessage)
	}

	ok, err := c.pivnetClient.Auth()
	if err != nil {
		return c.eh.HandleError(err)
	}

	if !ok {
		err := fmt.Errorf("Failed to login")
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
