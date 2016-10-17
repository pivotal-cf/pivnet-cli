package pivnet

import (
	"fmt"
	"net/http"
	"strings"
)

type EULAsService struct {
	client Client
}

type EULA struct {
	Slug    string `json:"slug,omitempty" yaml:"slug,omitempty"`
	ID      int    `json:"id,omitempty" yaml:"id,omitempty"`
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`
	Content string `json:"content,omitempty" yaml:"content,omitempty"`
	Links   *Links `json:"_links,omitempty" yaml:"_links,omitempty"`
}

type EULAsResponse struct {
	EULAs []EULA `json:"eulas,omitempty"`
	Links *Links `json:"_links,omitempty"`
}

type EULAAcceptanceResponse struct {
	AcceptedAt string `json:"accepted_at,omitempty"`
	Links      *Links `json:"_links,omitempty"`
}

func (e EULAsService) List() ([]EULA, error) {
	url := "/eulas"

	var response EULAsResponse
	_, _, err := e.client.MakeRequest(
		"GET",
		url,
		http.StatusOK,
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return response.EULAs, nil
}

func (e EULAsService) Get(eulaSlug string) (EULA, error) {
	url := fmt.Sprintf("/eulas/%s", eulaSlug)

	var response EULA
	_, _, err := e.client.MakeRequest(
		"GET",
		url,
		http.StatusOK,
		nil,
		&response,
	)
	if err != nil {
		return EULA{}, err
	}

	return response, nil
}

func (e EULAsService) Accept(productSlug string, releaseID int) error {
	url := fmt.Sprintf(
		"/products/%s/releases/%d/eula_acceptance",
		productSlug,
		releaseID,
	)

	var response EULAAcceptanceResponse
	_, _, err := e.client.MakeRequest(
		"POST",
		url,
		http.StatusOK,
		strings.NewReader(`{}`),
		&response,
	)
	if err != nil {
		return err
	}

	return nil
}
