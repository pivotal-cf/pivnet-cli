package pivnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ReleaseDependenciesService struct {
	client Client
}

type ReleaseDependenciesResponse struct {
	ReleaseDependencies []ReleaseDependency `json:"dependencies,omitempty"`
}

type ReleaseDependency struct {
	Release DependentRelease `json:"release,omitempty" yaml:"release,omitempty"`
}

type DependentRelease struct {
	ID      int     `json:"id,omitempty" yaml:"id,omitempty"`
	Version string  `json:"version,omitempty" yaml:"version,omitempty"`
	Product Product `json:"product,omitempty" yaml:"product,omitempty"`
}

func (r ReleaseDependenciesService) List(productSlug string, releaseID int) ([]ReleaseDependency, error) {
	url := fmt.Sprintf(
		"/products/%s/releases/%d/dependencies",
		productSlug,
		releaseID,
	)

	var response ReleaseDependenciesResponse
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

	return response.ReleaseDependencies, nil
}

func (r ReleaseDependenciesService) Add(
	productSlug string,
	releaseID int,
	dependentReleaseID int,
) error {
	url := fmt.Sprintf(
		"/products/%s/releases/%d/add_dependency",
		productSlug,
		releaseID,
	)

	body := addRemoveDependencyBody{
		Dependency: addRemoveDependencyBodyDependency{
			ReleaseID: dependentReleaseID,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		// Untested as we cannot force an error because we are marshalling
		// a known-good body
		return err
	}

	_, _, err = r.client.MakeRequest(
		"PATCH",
		url,
		http.StatusNoContent,
		bytes.NewReader(b),
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r ReleaseDependenciesService) Remove(
	productSlug string,
	releaseID int,
	dependentReleaseID int,
) error {
	url := fmt.Sprintf(
		"/products/%s/releases/%d/remove_dependency",
		productSlug,
		releaseID,
	)

	body := addRemoveDependencyBody{
		Dependency: addRemoveDependencyBodyDependency{
			ReleaseID: dependentReleaseID,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		// Untested as we cannot force an error because we are marshalling
		// a known-good body
		return err
	}

	_, _, err = r.client.MakeRequest(
		"PATCH",
		url,
		http.StatusNoContent,
		bytes.NewReader(b),
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

type addRemoveDependencyBody struct {
	Dependency addRemoveDependencyBodyDependency `json:"dependency"`
}

type addRemoveDependencyBodyDependency struct {
	ReleaseID int `json:"release_id"`
}
