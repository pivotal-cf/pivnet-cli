package productfile

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/errorhandler"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/printer"
	"github.com/pivotal-cf/go-pivnet/logger"
	"gopkg.in/cheggaaa/pb.v1"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
	GetProductFiles(productSlug string) ([]pivnet.ProductFile, error)
	GetProductFilesForRelease(productSlug string, releaseID int) ([]pivnet.ProductFile, error)
	GetProductFile(productSlug string, productFileID int) (pivnet.ProductFile, error)
	GetProductFileForRelease(productSlug string, releaseID int, productFileID int) (pivnet.ProductFile, error)
	CreateProductFile(config pivnet.CreateProductFileConfig) (pivnet.ProductFile, error)
	UpdateProductFile(productSlug string, productFile pivnet.ProductFile) (pivnet.ProductFile, error)
	AddProductFileToRelease(productSlug string, releaseID int, productFileID int) error
	RemoveProductFileFromRelease(productSlug string, releaseID int, productFileID int) error
	AddProductFileToFileGroup(productSlug string, fileGroupID int, productFileID int) error
	RemoveProductFileFromFileGroup(productSlug string, fileGroupID int, productFileID int) error
	DeleteProductFile(productSlug string, releaseID int) (pivnet.ProductFile, error)
	AcceptEULA(productSlug string, releaseID int) error
	DownloadFile(writer io.Writer, downloadLink string) error
}

type ProductFileClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	logWriter    io.Writer
	printer      printer.Printer
	l            logger.Logger
}

func NewProductFileClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	logWriter io.Writer,
	printer printer.Printer,
	l logger.Logger,
) *ProductFileClient {
	return &ProductFileClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		logWriter:    logWriter,
		printer:      printer,
		l:            l,
	}
}

func (c *ProductFileClient) List(productSlug string, releaseVersion string) error {
	if releaseVersion == "" {
		productFiles, err := c.pivnetClient.GetProductFiles(productSlug)
		if err != nil {
			return c.eh.HandleError(err)
		}

		return c.printProductFiles(productFiles)
	}

	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	productFiles, err := c.pivnetClient.GetProductFilesForRelease(
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
		productFile, err := c.pivnetClient.GetProductFile(
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

	productFile, err := c.pivnetClient.GetProductFileForRelease(
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
	productFile, err := c.pivnetClient.GetProductFile(
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
	productFileID int,
	filepath string,
	acceptEULA bool,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	downloadLink := fmt.Sprintf(
		"/products/%s/releases/%d/product_files/%d/download",
		productSlug,
		release.ID,
		productFileID,
	)

	c.l.Debug(
		"Creating local file",
		logger.Data{"downloadLink": downloadLink, "localFilepath": filepath},
	)
	file, err := os.Create(filepath)
	if err != nil {
		return c.eh.HandleError(err)
	}

	c.l.Debug("Determining file size", logger.Data{"downloadLink": downloadLink})
	productFile, err := c.pivnetClient.GetProductFileForRelease(
		productSlug,
		release.ID,
		productFileID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if acceptEULA {
		c.l.Debug("Accepting EULA")
		err = c.pivnetClient.AcceptEULA(productSlug, release.ID)
		if err != nil {
			return c.eh.HandleError(err)
		}
	}

	progress := newProgressBar(productFile.Size, c.logWriter)
	onDemandProgress := &startOnDemandProgressBar{progress, false}

	multiWriter := io.MultiWriter(file, onDemandProgress)

	c.l.Debug(
		"Downloading link to local file",
		logger.Data{
			"downloadLink":  downloadLink,
			"localFilepath": filepath,
		},
	)
	err = c.pivnetClient.DownloadFile(multiWriter, downloadLink)
	if err != nil {
		return c.eh.HandleError(err)
	}

	progress.Finish()
	return nil
}

type startOnDemandProgressBar struct {
	progressbar *pb.ProgressBar
	started     bool
}

func (w *startOnDemandProgressBar) Write(b []byte) (int, error) {
	if !w.started {
		w.progressbar.Start()
		w.started = true
	}
	return w.progressbar.Write(b)
}

func newProgressBar(total int, output io.Writer) *pb.ProgressBar {
	progress := pb.New(total)

	progress.Output = output
	progress.ShowSpeed = true
	progress.Units = pb.U_BYTES
	progress.NotPrint = true

	return progress.SetWidth(80)
}
