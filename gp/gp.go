package gp

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/extension"
	"github.com/pivotal-cf/go-pivnet/logger"
)

type Client struct {
	client pivnet.Client
}

func NewClient(config pivnet.ClientConfig, logger logger.Logger) *Client {
	return &Client{
		client: pivnet.NewClient(config, logger),
	}
}

type ExtendedClient struct {
	client extension.ExtendedClient
}

func NewExtendedClient(c Client, logger logger.Logger) *ExtendedClient {
	return &ExtendedClient{
		client: extension.NewExtendedClient(c, logger),
	}
}

type CompositeClient struct {
	Client
	ExtendedClient
}

func NewCompositeClient(config pivnet.ClientConfig, logger logger.Logger) *CompositeClient {
	c := NewClient(config, logger)
	e := NewExtendedClient(*c, logger)
	return &CompositeClient{
		*c,
		*e,
	}
}

func (c Client) ReleaseTypes() ([]pivnet.ReleaseType, error) {
	return c.client.ReleaseTypes.Get()
}

func (c Client) ReleasesForProductSlug(productSlug string) ([]pivnet.Release, error) {
	return c.client.Releases.List(productSlug)
}

func (c Client) Release(productSlug string, releaseID int) (pivnet.Release, error) {
	return c.client.Releases.Get(productSlug, releaseID)
}

func (c Client) ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error) {
	releases, err := c.ReleasesForProductSlug(productSlug)
	if err != nil {
		return pivnet.Release{}, err
	}

	release, err := c.releaseForReleaseVersion(releases, releaseVersion)
	if err != nil {
		return pivnet.Release{}, err
	}

	return c.client.Releases.Get(productSlug, release.ID)
}

func (c Client) releaseForReleaseVersion(releases []pivnet.Release, releaseVersion string) (pivnet.Release, error) {
	for _, r := range releases {
		if r.Version == releaseVersion {
			return r, nil
		}
	}

	return pivnet.Release{}, fmt.Errorf(
		"release not found for version: '%s'",
		releaseVersion,
	)
}

func (c Client) UpdateRelease(productSlug string, release pivnet.Release) (pivnet.Release, error) {
	return c.client.Releases.Update(productSlug, release)
}

func (c Client) CreateRelease(config pivnet.CreateReleaseConfig) (pivnet.Release, error) {
	return c.client.Releases.Create(config)
}

func (c Client) DeleteRelease(productSlug string, release pivnet.Release) error {
	return c.client.Releases.Delete(productSlug, release)
}

func (c Client) AddUserGroup(productSlug string, releaseID int, userGroupID int) error {
	return c.client.UserGroups.AddToRelease(productSlug, releaseID, userGroupID)
}

func (c Client) RemoveUserGroup(productSlug string, releaseID int, userGroupID int) error {
	return c.client.UserGroups.RemoveFromRelease(productSlug, releaseID, userGroupID)
}

func (c Client) UserGroups() ([]pivnet.UserGroup, error) {
	return c.client.UserGroups.List()
}

func (c Client) UserGroupsForRelease(productSlug string, releaseID int) ([]pivnet.UserGroup, error) {
	return c.client.UserGroups.ListForRelease(productSlug, releaseID)
}

func (c Client) UserGroup(userGroupID int) (pivnet.UserGroup, error) {
	return c.client.UserGroups.Get(userGroupID)
}

func (c Client) CreateUserGroup(name string, description string, members []string) (pivnet.UserGroup, error) {
	return c.client.UserGroups.Create(name, description, members)
}

func (c Client) UpdateUserGroup(userGroup pivnet.UserGroup) (pivnet.UserGroup, error) {
	return c.client.UserGroups.Update(userGroup)
}

func (c Client) DeleteUserGroup(userGroupID int) error {
	return c.client.UserGroups.Delete(userGroupID)
}

func (c Client) AddMemberToGroup(userGroupID int, emailAddress string, admin bool) (pivnet.UserGroup, error) {
	return c.client.UserGroups.AddMemberToGroup(userGroupID, emailAddress, admin)
}

func (c Client) RemoveMemberFromGroup(userGroupID int, emailAddress string) (pivnet.UserGroup, error) {
	return c.client.UserGroups.RemoveMemberFromGroup(userGroupID, emailAddress)
}

func (c Client) EULA(eulaSlug string) (pivnet.EULA, error) {
	return c.client.EULA.Get(eulaSlug)
}

func (c Client) AcceptEULA(productSlug string, releaseID int) error {
	return c.client.EULA.Accept(productSlug, releaseID)
}

func (c Client) EULAs() ([]pivnet.EULA, error) {
	return c.client.EULA.List()
}

func (c Client) GetProductFilesForRelease(productSlug string, releaseID int) ([]pivnet.ProductFile, error) {
	return c.client.ProductFiles.ListForRelease(productSlug, releaseID)
}

func (c Client) GetProductFiles(productSlug string) ([]pivnet.ProductFile, error) {
	return c.client.ProductFiles.List(productSlug)
}

func (c Client) GetProductFileForRelease(productSlug string, releaseID int, productFileID int) (pivnet.ProductFile, error) {
	return c.client.ProductFiles.GetForRelease(productSlug, releaseID, productFileID)
}

func (c Client) GetProductFile(productSlug string, productFileID int) (pivnet.ProductFile, error) {
	return c.client.ProductFiles.Get(productSlug, productFileID)
}

func (c Client) DeleteProductFile(productSlug string, releaseID int) (pivnet.ProductFile, error) {
	return c.client.ProductFiles.Delete(productSlug, releaseID)
}

func (c Client) Products() ([]pivnet.Product, error) {
	return c.client.Products.List()
}

func (c Client) FindProductForSlug(slug string) (pivnet.Product, error) {
	return c.client.Products.Get(slug)
}

func (c Client) CreateProductFile(config pivnet.CreateProductFileConfig) (pivnet.ProductFile, error) {
	return c.client.ProductFiles.Create(config)
}

func (c Client) UpdateProductFile(productSlug string, productFile pivnet.ProductFile) (pivnet.ProductFile, error) {
	return c.client.ProductFiles.Update(productSlug, productFile)
}

func (c Client) AddProductFileToRelease(productSlug string, releaseID int, productFileID int) error {
	return c.client.ProductFiles.AddToRelease(productSlug, releaseID, productFileID)
}

func (c Client) RemoveProductFileFromRelease(productSlug string, releaseID int, productFileID int) error {
	return c.client.ProductFiles.RemoveFromRelease(productSlug, releaseID, productFileID)
}

func (c Client) AddProductFileToFileGroup(productSlug string, fileGroupID int, productFileID int) error {
	return c.client.ProductFiles.AddToFileGroup(productSlug, fileGroupID, productFileID)
}

func (c Client) RemoveProductFileFromFileGroup(productSlug string, fileGroupID int, productFileID int) error {
	return c.client.ProductFiles.RemoveFromFileGroup(productSlug, fileGroupID, productFileID)
}

func (c Client) ReleaseDependencies(productSlug string, releaseID int) ([]pivnet.ReleaseDependency, error) {
	return c.client.ReleaseDependencies.List(productSlug, releaseID)
}

func (c Client) AddReleaseDependency(productSlug string, releaseID int, dependentReleaseID int) error {
	return c.client.ReleaseDependencies.Add(productSlug, releaseID, dependentReleaseID)
}

func (c Client) RemoveReleaseDependency(productSlug string, releaseID int, dependentReleaseID int) error {
	return c.client.ReleaseDependencies.Remove(productSlug, releaseID, dependentReleaseID)
}

func (c Client) ReleaseUpgradePaths(productSlug string, releaseID int) ([]pivnet.ReleaseUpgradePath, error) {
	return c.client.ReleaseUpgradePaths.Get(productSlug, releaseID)
}

func (c Client) AddReleaseUpgradePath(productSlug string, releaseID int, previousReleaseID int) error {
	return c.client.ReleaseUpgradePaths.Add(productSlug, releaseID, previousReleaseID)
}

func (c Client) RemoveReleaseUpgradePath(productSlug string, releaseID int, previousReleaseID int) error {
	return c.client.ReleaseUpgradePaths.Remove(productSlug, releaseID, previousReleaseID)
}

func (c Client) FileGroups(productSlug string) ([]pivnet.FileGroup, error) {
	return c.client.FileGroups.List(productSlug)
}

func (c Client) FileGroupsForRelease(productSlug string, releaseID int) ([]pivnet.FileGroup, error) {
	return c.client.FileGroups.ListForRelease(productSlug, releaseID)
}

func (c Client) FileGroup(productSlug string, fileGroupID int) (pivnet.FileGroup, error) {
	return c.client.FileGroups.Get(productSlug, fileGroupID)
}

func (c Client) CreateFileGroup(productSlug string, name string) (pivnet.FileGroup, error) {
	return c.client.FileGroups.Create(productSlug, name)
}

func (c Client) UpdateFileGroup(productSlug string, fileGroup pivnet.FileGroup) (pivnet.FileGroup, error) {
	return c.client.FileGroups.Update(productSlug, fileGroup)
}

func (c Client) DeleteFileGroup(productSlug string, fileGroupID int) (pivnet.FileGroup, error) {
	return c.client.FileGroups.Delete(productSlug, fileGroupID)
}

func (c Client) AddFileGroupToRelease(productSlug string, fileGroupID int, releaseID int) error {
	return c.client.FileGroups.AddToRelease(productSlug, fileGroupID, releaseID)
}

func (c Client) RemoveFileGroupFromRelease(productSlug string, fileGroupID int, releaseID int) error {
	return c.client.FileGroups.RemoveFromRelease(productSlug, fileGroupID, releaseID)
}

func (c Client) MakeRequest(method string, url string, expectedResponseCode int, body io.Reader, data interface{}) (*http.Response, []byte, error) {
	return c.client.MakeRequest(method, url, expectedResponseCode, body, data)
}

func (c Client) CreateRequest(method string, url string, body io.Reader) (*http.Request, error) {
	return c.client.CreateRequest(method, url, body)
}

func (c ExtendedClient) ReleaseFingerprint(productSlug string, releaseID int) (string, error) {
	return c.client.ReleaseFingerprint(productSlug, releaseID)
}

func (c ExtendedClient) DownloadFile(writer io.Writer, downloadLink string) error {
	return c.client.DownloadFile(writer, downloadLink)
}
