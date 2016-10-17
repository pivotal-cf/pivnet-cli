package pivnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProductFilesService struct {
	client Client
}

type CreateProductFileConfig struct {
	ProductSlug  string
	FileVersion  string
	AWSObjectKey string
	Name         string
	MD5          string
	Description  string
	FileType     string
}

type ProductFilesResponse struct {
	ProductFiles []ProductFile `json:"product_files,omitempty"`
}

type ProductFileResponse struct {
	ProductFile ProductFile `json:"product_file,omitempty"`
}

type ProductFile struct {
	ID           int    `json:"id,omitempty" yaml:"id,omitempty"`
	AWSObjectKey string `json:"aws_object_key,omitempty" yaml:"aws_object_key,omitempty"`
	Links        *Links `json:"_links,omitempty" yaml:"_links,omitempty"`
	FileType     string `json:"file_type,omitempty" yaml:"file_type,omitempty"`
	FileVersion  string `json:"file_version,omitempty" yaml:"file_version,omitempty"`
	Name         string `json:"name,omitempty" yaml:"name,omitempty"`
	MD5          string `json:"md5,omitempty" yaml:"md5,omitempty"`
	Description  string `json:"description,omitempty" yaml:"description,omitempty"`
	Size         int    `json:"size,omitempty" yaml:"size,omitempty"`
}

const (
	FileTypeSoftware          = "Software"
	FileTypeDocumentation     = "Documentation"
	FileTypeOpenSourceLicense = "Open Source License"
)

func (p ProductFilesService) List(productSlug string) ([]ProductFile, error) {
	url := fmt.Sprintf("/products/%s/product_files", productSlug)

	var response ProductFilesResponse
	_, _, err := p.client.MakeRequest(
		"GET",
		url,
		http.StatusOK,
		nil,
		&response,
	)
	if err != nil {
		return []ProductFile{}, err
	}

	return response.ProductFiles, nil
}

func (p ProductFilesService) ListForRelease(productSlug string, releaseID int) ([]ProductFile, error) {
	url := fmt.Sprintf(
		"/products/%s/releases/%d/product_files",
		productSlug,
		releaseID,
	)

	var response ProductFilesResponse
	_, _, err := p.client.MakeRequest(
		"GET",
		url,
		http.StatusOK,
		nil,
		&response,
	)
	if err != nil {
		return []ProductFile{}, err
	}

	return response.ProductFiles, nil
}

func (p ProductFilesService) Get(productSlug string, productFileID int) (ProductFile, error) {
	url := fmt.Sprintf(
		"/products/%s/product_files/%d",
		productSlug,
		productFileID,
	)

	var response ProductFileResponse
	_, _, err := p.client.MakeRequest(
		"GET",
		url,
		http.StatusOK,
		nil,
		&response,
	)
	if err != nil {
		return ProductFile{}, err
	}

	return response.ProductFile, nil
}

func (p ProductFilesService) GetForRelease(productSlug string, releaseID int, productFileID int) (ProductFile, error) {
	url := fmt.Sprintf(
		"/products/%s/releases/%d/product_files/%d",
		productSlug,
		releaseID,
		productFileID,
	)

	var response ProductFileResponse
	_, _, err := p.client.MakeRequest(
		"GET",
		url,
		http.StatusOK,
		nil,
		&response,
	)
	if err != nil {
		return ProductFile{}, err
	}

	return response.ProductFile, nil
}

func (p ProductFilesService) Create(config CreateProductFileConfig) (ProductFile, error) {
	if config.AWSObjectKey == "" {
		return ProductFile{}, fmt.Errorf("AWS object key must not be empty")
	}

	url := fmt.Sprintf("/products/%s/product_files", config.ProductSlug)

	body := createUpdateProductFileBody{
		ProductFile: ProductFile{
			MD5:          config.MD5,
			FileType:     config.FileType,
			FileVersion:  config.FileVersion,
			AWSObjectKey: config.AWSObjectKey,
			Name:         config.Name,
			Description:  config.Description,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		// Untested as we cannot force an error because we are marshalling
		// a known-good body
		return ProductFile{}, err
	}

	var response ProductFileResponse
	_, _, err = p.client.MakeRequest(
		"POST",
		url,
		http.StatusCreated,
		bytes.NewReader(b),
		&response,
	)
	if err != nil {
		return ProductFile{}, err
	}

	return response.ProductFile, nil
}

func (p ProductFilesService) Update(productSlug string, productFile ProductFile) (ProductFile, error) {
	url := fmt.Sprintf("/products/%s/product_files/%d", productSlug, productFile.ID)

	body := createUpdateProductFileBody{
		ProductFile: ProductFile{
			Description: productFile.Description,
			FileType:    productFile.FileType,
			FileVersion: productFile.FileVersion,
			MD5:         productFile.MD5,
			Name:        productFile.Name,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		// Untested as we cannot force an error because we are marshalling
		// a known-good body
		return ProductFile{}, err
	}

	var response ProductFileResponse
	_, _, err = p.client.MakeRequest(
		"PATCH",
		url,
		http.StatusOK,
		bytes.NewReader(b),
		&response,
	)
	if err != nil {
		return ProductFile{}, err
	}

	return response.ProductFile, nil
}

type createUpdateProductFileBody struct {
	ProductFile ProductFile `json:"product_file"`
}

func (p ProductFilesService) Delete(productSlug string, id int) (ProductFile, error) {
	url := fmt.Sprintf(
		"/products/%s/product_files/%d",
		productSlug,
		id,
	)

	var response ProductFileResponse
	_, _, err := p.client.MakeRequest(
		"DELETE",
		url,
		http.StatusOK,
		nil,
		&response,
	)
	if err != nil {
		return ProductFile{}, err
	}

	return response.ProductFile, nil
}

func (p ProductFilesService) AddToRelease(
	productSlug string,
	releaseID int,
	productFileID int,
) error {
	url := fmt.Sprintf(
		"/products/%s/releases/%d/add_product_file",
		productSlug,
		releaseID,
	)

	body := createUpdateProductFileBody{
		ProductFile: ProductFile{
			ID: productFileID,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		// Untested as we cannot force an error because we are marshalling
		// a known-good body
		return err
	}

	_, _, err = p.client.MakeRequest(
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

func (p ProductFilesService) RemoveFromRelease(
	productSlug string,
	releaseID int,
	productFileID int,
) error {
	url := fmt.Sprintf(
		"/products/%s/releases/%d/remove_product_file",
		productSlug,
		releaseID,
	)

	body := createUpdateProductFileBody{
		ProductFile: ProductFile{
			ID: productFileID,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		// Untested as we cannot force an error because we are marshalling
		// a known-good body
		return err
	}

	_, _, err = p.client.MakeRequest(
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

func (p ProductFilesService) AddToFileGroup(
	productSlug string,
	fileGroupID int,
	productFileID int,
) error {
	url := fmt.Sprintf(
		"/products/%s/file_groups/%d/add_product_file",
		productSlug,
		fileGroupID,
	)

	body := createUpdateProductFileBody{
		ProductFile: ProductFile{
			ID: productFileID,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		// Untested as we cannot force an error because we are marshalling
		// a known-good body
		return err
	}

	_, _, err = p.client.MakeRequest(
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

func (p ProductFilesService) RemoveFromFileGroup(
	productSlug string,
	fileGroupID int,
	productFileID int,
) error {
	url := fmt.Sprintf(
		"/products/%s/file_groups/%d/remove_product_file",
		productSlug,
		fileGroupID,
	)

	body := createUpdateProductFileBody{
		ProductFile: ProductFile{
			ID: productFileID,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		// Untested as we cannot force an error because we are marshalling
		// a known-good body
		return err
	}

	_, _, err = p.client.MakeRequest(
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
