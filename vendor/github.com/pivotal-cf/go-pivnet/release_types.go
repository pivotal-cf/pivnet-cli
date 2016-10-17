package pivnet

import (
	"fmt"
	"net/http"
)

type ReleaseTypesService struct {
	client Client
}

type ReleaseType string

type ReleaseTypesResponse struct {
	ReleaseTypes []ReleaseType `json:"release_types" yaml:"release_types"`
}

func (r ReleaseTypesService) Get() ([]ReleaseType, error) {
	url := fmt.Sprintf("/releases/release_types")

	var response ReleaseTypesResponse
	_, _, err := r.client.MakeRequest(
		"GET",
		url,
		http.StatusOK,
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return response.ReleaseTypes, nil
}
