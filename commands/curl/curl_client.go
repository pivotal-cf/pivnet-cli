package curl

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	MakeRequest(method string, url string, expectedResponseCode int, body io.Reader) (*http.Response, error)
}

type CurlClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewCurlClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *CurlClient {
	return &CurlClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *CurlClient) MakeRequest(
	method string,
	data string,
	args []string,
) error {
	expectedResponseCode := 0
	var body io.Reader
	if data != "" {
		body = strings.NewReader(data)
	}

	var output interface{}
	resp, err := c.pivnetClient.MakeRequest(
		method,
		args[0],
		expectedResponseCode,
		body,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return c.eh.HandleError(err)
	}

	c.printInterface(output)

	return nil
}

func (c *CurlClient) printInterface(object interface{}) error {
	if c.format == printer.PrintAsYAML {
		return c.printer.PrintYAML(object)
	}

	return c.printer.PrintJSON(object)
}
