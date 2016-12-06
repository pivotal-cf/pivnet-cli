package download

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

//go:generate counterfeiter -o ./fakes/ranger.go --fake-name Ranger . ranger
type ranger interface {
	BuildRange(contentLength int64) ([]Range, error)
}

//go:generate counterfeiter -o ./fakes/http_client.go --fake-name HTTPClient . httpClient
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

//go:generate counterfeiter -o ./fakes/bar.go --fake-name Bar . bar
type bar interface {
	SetTotal(contentLength int64)
	Add(totalWritten int) int
	Kickoff()
	Finish()
}

type Client struct {
	httpClient httpClient
	ranger     ranger
	bar        bar
}

func New(httpClient httpClient, ranger ranger, bar bar) Client {
	return Client{
		httpClient: httpClient,
		ranger:     ranger,
		bar:        bar,
	}
}

func (c Client) Get(location *os.File, contentURL string) error {
	req, err := http.NewRequest("HEAD", contentURL, nil)
	if err != nil {
		return fmt.Errorf("failed to construct HEAD request: %s", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HEAD request: %s", err)
	}

	contentURL = resp.Request.URL.String()

	ranges, err := c.ranger.BuildRange(resp.ContentLength)
	if err != nil {
		return fmt.Errorf("failed to construct range: %s", err)
	}

	c.bar.SetTotal(resp.ContentLength)
	c.bar.Kickoff()

	defer c.bar.Finish()

	var g errgroup.Group
	for _, r := range ranges {
		byteRange := r
		g.Go(func() error {
			respBytes, err := c.retryableRequest(contentURL, byteRange.HTTPHeader)
			if err != nil {
				return fmt.Errorf("failed during retryable request: %s", err)
			}

			bytesWritten, err := location.WriteAt(respBytes, byteRange.Lower)
			if err != nil {
				return fmt.Errorf("failed to write file: %s", err)
			}

			c.bar.Add(bytesWritten)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (c Client) retryableRequest(url string, rangeHeader http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header = rangeHeader

Retry:
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok {
			if netErr.Temporary() {
				goto Retry
			}
		}

		return []byte{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return []byte{}, fmt.Errorf("during GET unexpected status code was returned: %d", resp.StatusCode)
	}

	var respBytes []byte
	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			goto Retry
		}

		return []byte{}, err
	}

	return respBytes, err
}
