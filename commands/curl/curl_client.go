package curl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"io"
	"net/http"
	"strings"

	"github.com/pivotal-cf/pivnet-cli/v2/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/v2/printer"
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
	url string,
) error {
	expectedResponseCode := 0
	var body io.Reader
	if data != "" {
		body = strings.NewReader(data)
	}

	var output interface{}
	resp, err := c.pivnetClient.MakeRequest(
		method,
		url,
		expectedResponseCode,
		body,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	if strings.ToLower(contentType) == "text/csv" {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.eh.HandleError(err)
		}

		fmt.Println(string(bodyBytes))

		return nil
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
