package productfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/go-pivnet/v2/download"
	"github.com/pivotal-cf/go-pivnet/v2/logger"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
	"github.com/pivotal-cf/pivnet-cli/ui"
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
	DeleteProductFile(productSlug string, productFileID int) (pivnet.ProductFile, error)
	AcceptEULA(productSlug string, releaseID int) error
	DownloadProductFile(location *download.FileInfo, productSlug string, releaseID int, productFileID int, progressWriter io.Writer) error
}

//go:generate counterfeiter . Filter
type Filter interface {
	ProductFileKeysByGlobs(productFiles []pivnet.ProductFile, glob []string) ([]pivnet.ProductFile, error)
}

//go:generate counterfeiter --fake-name FakeFileSummer . FileSummer
type FileSummer interface {
	SumFile(filepath string) (string, error)
}

type ProductFileClient struct {
	pivnetClient     PivnetClient
	sha256FileSummer FileSummer
	md5FileSummer    FileSummer
	eh               errorhandler.ErrorHandler
	format           string
	outputWriter     io.Writer
	logWriter        io.Writer
	printer          printer.Printer
	l                logger.Logger
	filter           Filter
}

func NewProductFileClient(
	pivnetClient PivnetClient,
	sha256FileSummer FileSummer,
	md5FileSummer FileSummer,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	logWriter io.Writer,
	printer printer.Printer,
	l logger.Logger,
	filter Filter,
) *ProductFileClient {
	return &ProductFileClient{
		pivnetClient:     pivnetClient,
		sha256FileSummer: sha256FileSummer,
		md5FileSummer:    md5FileSummer,
		eh:               eh,
		format:           format,
		outputWriter:     outputWriter,
		logWriter:        logWriter,
		printer:          printer,
		l:                l,
		filter:           filter,
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
			"File Type",
			"SHA256",
			"AWS Object Key",
		})

		for _, productFile := range productFiles {
			productFileAsString := []string{
				strconv.Itoa(productFile.ID),
				productFile.Name,
				productFile.FileVersion,
				productFile.FileType,
				productFile.SHA256,
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
			"SHA256",
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
			productFile.SHA256,
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
	fileVersion *string,
	sha256 *string,
	md5 *string,
	description *string,
	docsURL *string,
	systemRequirements *[]string,
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

	if fileVersion != nil {
		productFile.FileVersion = *fileVersion
	}

	if sha256 != nil {
		productFile.SHA256 = *sha256
	}

	if md5 != nil {
		productFile.MD5 = *md5
	}

	if description != nil {
		productFile.Description = *description
	}

	if docsURL != nil {
		productFile.DocsURL = *docsURL
	}
	
	if systemRequirements != nil {
		productFile.SystemRequirements = *systemRequirements
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
		message := fmt.Sprintf(
			"Product file %d added to %s/%s",
			productFileID,
			productSlug,
			releaseVersion,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err := fmt.Fprintln(c.outputWriter, coloredMessage)

		return err
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
		message := fmt.Sprintf(
			"Product file %d removed from %s/%s",
			productFileID,
			productSlug,
			releaseVersion,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err := fmt.Fprintln(c.outputWriter, coloredMessage)

		return err
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
		message := fmt.Sprintf(
			"Product file %d added to file group %d",
			productFileID,
			fileGroupID,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err := fmt.Fprintln(c.outputWriter, coloredMessage)

		return err
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
		message := fmt.Sprintf(
			"Product file %d removed from file group %d",
			productFileID,
			fileGroupID,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err := fmt.Fprintln(c.outputWriter, coloredMessage)

		return err
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
		message := fmt.Sprintf(
			"Product file %d deleted for %s",
			productFileID,
			productSlug,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err := fmt.Fprintln(c.outputWriter, coloredMessage)

		return err
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
	progressWriter io.Writer,
) error {
	if len(globs) > 0 && len(productFileIDs) > 0 {
		err := fmt.Errorf("Cannot provide both globs and product file IDs")
		return c.eh.HandleError(err)
	}

	if len(globs) == 0 && len(productFileIDs) == 0 {
		err := fmt.Errorf("Must provide either globs (-g) or product file IDs (-i)")
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

	if len(filteredProductFiles) == 0 {
		err := fmt.Errorf(
			"No product files found for ids: %v or globs: %v",
			productFileIDs,
			globs,
		)
		return c.eh.HandleError(err)
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

		fileInfo, err := download.NewFileInfo(file)
		if err != nil {
			return c.eh.HandleError(err)
		}

		err = file.Close()
		if err != nil {
			return c.eh.HandleError(err)
		}

		c.l.Info(fmt.Sprintf(
			"Downloading '%s' to '%s'",
			fileName,
			localFilepath,
		))

		err = c.pivnetClient.DownloadProductFile(fileInfo, productSlug, release.ID, pf.ID, progressWriter)
		if err != nil {
			return c.eh.HandleError(err)
		}

		if pf.FileType == pivnet.FileTypeSoftware && pf.SHA256 == "" && pf.MD5 == "" {
			return fmt.Errorf("cannot check file integrity of file %s: missing sha256 and md5 fields", file.Name())
		}

		if pf.SHA256 != "" {
			c.l.Info("Verifying SHA256")

			actualSHA256, err := c.sha256FileSummer.SumFile(file.Name())
			if err != nil {
				return err
			}

			if actualSHA256 != pf.SHA256 {
				return fmt.Errorf(
					"SHA256 comparison failed for downloaded file: '%s'. Expected (from pivnet): '%s' - actual (from file): '%s'",
					file.Name(),
					pf.SHA256,
					actualSHA256,
				)
			}

			c.l.Info("Successfully verified SHA256")
		}

		if pf.MD5 != "" {
			c.l.Info("Verifying MD5")

			actualMD5, err := c.md5FileSummer.SumFile(file.Name())
			if err != nil {
				return err
			}

			if actualMD5 != pf.MD5 {
				return fmt.Errorf(
					"MD5 comparison failed for downloaded file: '%s'. Expected (from pivnet): '%s' - actual (from file): '%s'",
					file.Name(),
					pf.MD5,
					actualMD5,
				)
			}

			c.l.Info("Successfully verified MD5")
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
