package extension

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/pivotal-cf/go-pivnet/logger"
)

//go:generate counterfeiter . Client
type Client interface {
	MakeRequest(method string, url string, expectedResponseCode int, body io.Reader, data interface{}) (*http.Response, []byte, error)
	CreateRequest(method string, url string, body io.Reader) (*http.Request, error)
}

type ExtendedClient struct {
	c      Client
	logger logger.Logger
}

func NewExtendedClient(client Client, logger logger.Logger) ExtendedClient {
	return ExtendedClient{
		c:      client,
		logger: logger,
	}
}

func (c ExtendedClient) DownloadFile(writer io.Writer, downloadLink string) (err error, retryable bool) {
	c.logger.Debug("Downloading file", logger.Data{"downloadLink": downloadLink})

	req, err := c.c.CreateRequest(
		"POST",
		downloadLink,
		nil,
	)
	if err != nil {
		return err, false
	}

	reqBytes, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return err, false
	}

	c.logger.Debug("Making request", logger.Data{"request": string(reqBytes)})
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err, false
	}

	if resp.StatusCode == http.StatusUnavailableForLegalReasons {
		return errors.New(fmt.Sprintf("the EULA has not been accepted for the file: %s", downloadLink)), false
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("pivnet returned an error code of %d for the file: %s", resp.StatusCode, downloadLink)), false
	}

	c.logger.Debug("Copying body", logger.Data{"downloadLink": downloadLink})

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err, true
	}

	return nil, true
}
