package pivnet

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/pivotal-cf/go-pivnet/logger"
)

const (
	DefaultHost = "https://network.pivotal.io"
	apiVersion  = "/api/v2"
)

type pivnetErr struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

type pivnetInternalServerErr struct {
	Error string `json:"error"`
}

type ErrPivnetOther struct {
	ResponseCode int      `json:"response_code" yaml:"response_code"`
	Message      string   `json:"message" yaml:"message"`
	Errors       []string `json:"errors" yaml:"errors"`
}

func (e ErrPivnetOther) Error() string {
	return fmt.Sprintf(
		"%d - %s. Errors: %v",
		e.ResponseCode,
		e.Message,
		strings.Join(e.Errors, ","),
	)
}

type ErrUnauthorized struct {
	ResponseCode int    `json:"response_code" yaml:"response_code"`
	Message      string `json:"message" yaml:"message"`
}

func (e ErrUnauthorized) Error() string {
	return e.Message
}

func newErrUnauthorized(message string) ErrUnauthorized {
	return ErrUnauthorized{
		ResponseCode: http.StatusUnauthorized,
		Message:      message,
	}
}

type ErrNotFound struct {
	ResponseCode int    `json:"response_code" yaml:"response_code"`
	Message      string `json:"message" yaml:"message"`
}

func (e ErrNotFound) Error() string {
	return e.Message
}

func newErrNotFound(message string) ErrNotFound {
	return ErrNotFound{
		ResponseCode: http.StatusNotFound,
		Message:      message,
	}
}

type ErrUnavailableForLegalReasons struct {
	ResponseCode int    `json:"response_code" yaml:"response_code"`
	Message      string `json:"message" yaml:"message"`
}

func (e ErrUnavailableForLegalReasons) Error() string {
	return e.Message
}

func newErrUnavailableForLegalReasons() ErrUnavailableForLegalReasons {
	return ErrUnavailableForLegalReasons{
		ResponseCode: http.StatusUnavailableForLegalReasons,
		Message:      "The EULA has not been accepted.",
	}
}

type Client struct {
	baseURL           string
	token             string
	userAgent         string
	logger            logger.Logger
	skipSSLValidation bool

	Auth                *AuthService
	EULA                *EULAsService
	ProductFiles        *ProductFilesService
	FileGroups          *FileGroupsService
	Releases            *ReleasesService
	Products            *ProductsService
	UserGroups          *UserGroupsService
	ReleaseDependencies *ReleaseDependenciesService
	ReleaseTypes        *ReleaseTypesService
	ReleaseUpgradePaths *ReleaseUpgradePathsService
}

type ClientConfig struct {
	Host              string
	Token             string
	UserAgent         string
	SkipSSLValidation bool
}

func NewClient(config ClientConfig, logger logger.Logger) Client {
	baseURL := fmt.Sprintf("%s%s", config.Host, apiVersion)

	client := Client{
		baseURL:           baseURL,
		token:             config.Token,
		userAgent:         config.UserAgent,
		logger:            logger,
		skipSSLValidation: config.SkipSSLValidation,
	}

	client.Auth = &AuthService{client: client}
	client.EULA = &EULAsService{client: client}
	client.ProductFiles = &ProductFilesService{client: client}
	client.FileGroups = &FileGroupsService{client: client}
	client.Releases = &ReleasesService{client: client, l: logger}
	client.Products = &ProductsService{client: client, l: logger}
	client.UserGroups = &UserGroupsService{client: client}
	client.ReleaseDependencies = &ReleaseDependenciesService{client: client}
	client.ReleaseTypes = &ReleaseTypesService{client: client}
	client.ReleaseUpgradePaths = &ReleaseUpgradePathsService{client: client}

	return client
}

func (c Client) CreateRequest(
	requestType string,
	endpoint string,
	body io.Reader,
) (*http.Request, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}

	endpoint = c.stripHostPrefix(endpoint)

	u.Path = u.Path + endpoint

	req, err := http.NewRequest(requestType, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", c.token))
	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

func (c Client) MakeRequest(
	requestType string,
	endpoint string,
	expectedStatusCode int,
	body io.Reader,
) (*http.Response, error) {
	req, err := c.CreateRequest(requestType, endpoint, body)
	if err != nil {
		return nil, err
	}

	reqBytes, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("Making request", logger.Data{"request": string(reqBytes)})
	var httpClient *http.Client

	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.skipSSLValidation},
		},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("Response status code", logger.Data{"status code": resp.StatusCode})
	c.logger.Debug("Response headers", logger.Data{"headers": resp.Header})

	if expectedStatusCode > 0 && resp.StatusCode != expectedStatusCode {
		var pErr pivnetErr

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// We have to handle 500 differently because it has a different structure
		if resp.StatusCode == http.StatusInternalServerError {
			var internalServerError pivnetInternalServerErr
			err = json.Unmarshal(b, &internalServerError)
			if err != nil {
				return nil, err
			}

			pErr = pivnetErr{
				Message: internalServerError.Error,
			}
		} else {
			err = json.Unmarshal(b, &pErr)
			if err != nil {
				return nil, err
			}
		}

		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, newErrUnauthorized(pErr.Message)
		case http.StatusNotFound:
			return nil, newErrNotFound(pErr.Message)
		case http.StatusUnavailableForLegalReasons:
			return nil, newErrUnavailableForLegalReasons()
		default:
			return nil, ErrPivnetOther{
				ResponseCode: resp.StatusCode,
				Message:      pErr.Message,
				Errors:       pErr.Errors,
			}
		}
	}

	return resp, nil
}

func (c Client) stripHostPrefix(downloadLink string) string {
	if strings.HasPrefix(downloadLink, apiVersion) {
		return downloadLink
	}
	sp := strings.Split(downloadLink, apiVersion)
	return sp[len(sp)-1]
}
