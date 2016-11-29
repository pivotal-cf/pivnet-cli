package productfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logger"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
	ProductFiles(productSlug string) ([]pivnet.ProductFile, error)
	ProductFilesForRelease(productSlug string, releaseID int) ([]pivnet.ProductFile, error)
	ProductFile(productSlug string, productFileID int) (pivnet.ProductFile, error)
	ProductFileForRelease(productSlug string, releaseID int, productFileID int) (pivnet.ProductFile, error)
	CreateProductFile(config pivnet.CreateProductFileConfig) (pivnet.ProductFile, error)
	UpdateProductFile(productSlug string, productFile pivnet.ProductFile) (pivnet.ProductFile, error)
	AddProductFileToRelease(productSlug string, releaseID int, productFileID int) error
	RemoveProductFileFromRelease(productSlug string, releaseID int, productFileID int) error
	AddProductFileToFileGroup(productSlug string, fileGroupID int, productFileID int) error
	RemoveProductFileFromFileGroup(productSlug string, fileGroupID int, productFileID int) error
	DeleteProductFile(productSlug string, releaseID int) (pivnet.ProductFile, error)
	AcceptEULA(productSlug string, releaseID int) error
	DownloadProductFile(location *os.File, productSlug string, releaseID int, productFileID int) error
}

//go:generate counterfeiter . Filter
type Filter interface {
	ProductFileKeysByGlobs(productFiles []pivnet.ProductFile, glob []string) ([]pivnet.ProductFile, error)
}

type ProductFileClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	logWriter    io.Writer
	printer      printer.Printer
	l            logger.Logger
	filter       Filter
}

func NewProductFileClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	logWriter io.Writer,
	printer printer.Printer,
	l logger.Logger,
	filter Filter,
) *ProductFileClient {
	return &ProductFileClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		logWriter:    logWriter,
		printer:      printer,
		l:            l,
		filter:       filter,
	}
}

func (c *ProductFileClient) List(productSlug string, releaseVersion string) error {
	if releaseVersion == "" {
		productFiles, err := c.pivnetClient.ProductFiles(productSlug)
		if err != nil {
			return c.eh.HandleError(err)
		}

		return c.printProductFiles(productFiles)
	}

	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	productFiles, err := c.pivnetClient.ProductFilesForRelease(
		productSlug,
		release.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printProductFiles(productFiles)
}

func (c *ProductFileClient) printProductFiles(productFiles []pivnet.ProductFile) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Name",
			"File Version",
			"AWS Object Key",
		})

		for _, productFile := range productFiles {
			productFileAsString := []string{
				strconv.Itoa(productFile.ID),
				productFile.Name,
				productFile.FileVersion,
				productFile.AWSObjectKey,
			}
			table.Append(productFileAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(productFiles)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(productFiles)
	}

	return nil
}

func (c *ProductFileClient) printProductFile(productFile pivnet.ProductFile) error {
	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Name",
			"File Version",
			"File Type",
			"Description",
			"MD5",
			"AWS Object Key",
			"Size (Bytes)",
		})

		productFileAsString := []string{
			strconv.Itoa(productFile.ID),
			productFile.Name,
			productFile.FileVersion,
			productFile.FileType,
			productFile.Description,
			productFile.MD5,
			productFile.AWSObjectKey,
			fmt.Sprintf("%d", productFile.Size),
		}
		table.Append(productFileAsString)
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(productFile)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(productFile)
	}

	return nil
}

func (c *ProductFileClient) Get(
	productSlug string,
	releaseVersion string,
	productFileID int,
) error {
	if releaseVersion == "" {
		productFile, err := c.pivnetClient.ProductFile(
			productSlug,
			productFileID,
		)
		if err != nil {
			return c.eh.HandleError(err)
		}
		return c.printProductFile(productFile)
	}

	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	productFile, err := c.pivnetClient.ProductFileForRelease(
		productSlug,
		release.ID,
		productFileID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printProductFile(productFile)
}

func (c *ProductFileClient) Create(config pivnet.CreateProductFileConfig) error {
	productFile, err := c.pivnetClient.CreateProductFile(config)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printProductFile(productFile)
}

func (c *ProductFileClient) Update(
	productFileID int,
	productSlug string,
	name *string,
	fileType *string,
	fileVersion *string,
	md5 *string,
	description *string,
) error {
	productFile, err := c.pivnetClient.ProductFile(
		productSlug,
		productFileID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if name != nil {
		productFile.Name = *name
	}

	if fileType != nil {
		productFile.FileType = *fileType
	}

	if fileVersion != nil {
		productFile.FileVersion = *fileVersion
	}

	if md5 != nil {
		productFile.MD5 = *md5
	}

	if description != nil {
		productFile.Description = *description
	}

	updatedProductFile, err := c.pivnetClient.UpdateProductFile(productSlug, productFile)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printProductFile(updatedProductFile)
}

func (c *ProductFileClient) AddToRelease(
	productSlug string,
	releaseVersion string,
	productFileID int,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.AddProductFileToRelease(
		productSlug,
		release.ID,
		productFileID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"product file %d added successfully to %s/%s\n",
			productFileID,
			productSlug,
			releaseVersion,
		)
	}

	return nil
}

func (c *ProductFileClient) RemoveFromRelease(
	productSlug string,
	releaseVersion string,
	productFileID int,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.RemoveProductFileFromRelease(
		productSlug,
		release.ID,
		productFileID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"product file %d removed successfully from %s/%s\n",
			productFileID,
			productSlug,
			releaseVersion,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ProductFileClient) AddToFileGroup(
	productSlug string,
	fileGroupID int,
	productFileID int,
) error {
	err := c.pivnetClient.AddProductFileToFileGroup(
		productSlug,
		fileGroupID,
		productFileID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"product file %d successfully added to file group %d\n",
			productFileID,
			fileGroupID,
		)
	}

	return nil
}

func (c *ProductFileClient) RemoveFromFileGroup(
	productSlug string,
	fileGroupID int,
	productFileID int,
) error {
	err := c.pivnetClient.RemoveProductFileFromFileGroup(
		productSlug,
		fileGroupID,
		productFileID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"product file %d successfully removed from file group %d\n",
			productFileID,
			fileGroupID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ProductFileClient) Delete(productSlug string, productFileID int) error {
	productFile, err := c.pivnetClient.DeleteProductFile(
		productSlug,
		productFileID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"product file %d deleted successfully for %s\n",
			productFileID,
			productSlug,
		)
	}

	return c.printProductFile(productFile)
}

func (c *ProductFileClient) Download(
	productSlug string,
	releaseVersion string,
	globs []string,
	productFileIDs []int,
	downloadDir string,
	acceptEULA bool,
) error {
	if len(globs) > 0 && len(productFileIDs) > 0 {
		err := fmt.Errorf("Cannot provide both globs and product file IDs")
		return c.eh.HandleError(err)
	}

	if len(globs) == 0 && len(productFileIDs) == 0 {
		err := fmt.Errorf("Must provide either globs or product file IDs")
		return c.eh.HandleError(err)
	}

	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	productFiles, err := c.pivnetClient.ProductFilesForRelease(
		productSlug,
		release.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	var filteredProductFiles []pivnet.ProductFile

	if len(productFileIDs) > 0 {
		filteredProductFiles = filterProductFilesByIDs(productFiles, productFileIDs)
	}

	if len(globs) > 0 {
		var err error
		filteredProductFiles, err = c.filter.ProductFileKeysByGlobs(productFiles, globs)
		if err != nil {
			return c.eh.HandleError(err)
		}
	}

	if acceptEULA {
		c.l.Debug("Accepting EULA")
		err = c.pivnetClient.AcceptEULA(productSlug, release.ID)
		if err != nil {
			return c.eh.HandleError(err)
		}
	}

	for _, pf := range filteredProductFiles {
		parts := strings.Split(pf.AWSObjectKey, "/")
		fileName := parts[len(parts)-1]

		localFilepath := filepath.Join(downloadDir, fileName)

		c.l.Debug(
			"Creating local file",
			logger.Data{"name": pf.Name, "localFilepath": localFilepath},
		)
		file, err := os.Create(localFilepath)
		if err != nil {
			return c.eh.HandleError(err)
		}

		c.l.Info(fmt.Sprintf(
			"Downloading '%s' to '%s'",
			fileName,
			localFilepath,
		))

		err = c.pivnetClient.DownloadProductFile(file, productSlug, release.ID, pf.ID)
		if err != nil {
			return c.eh.HandleError(err)
		}
	}

	return nil
}

func filterProductFilesByIDs(productFiles []pivnet.ProductFile, ids []int) []pivnet.ProductFile {
	var foundProductFiles []pivnet.ProductFile
	for _, pf := range productFiles {
		if intContains(pf.ID, ids) {
			foundProductFiles = append(foundProductFiles, pf)
		}
	}
	return foundProductFiles
}

func intContains(i int, ints []int) bool {
	for _, val := range ints {
		if val == i {
			return true
		}
	}
	return false
}
