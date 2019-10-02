package imagereference

import (
	"fmt"
	"github.com/pivotal-cf/pivnet-cli/ui"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/go-pivnet/v2/logger"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
	CreateImageReference(config pivnet.CreateImageReferenceConfig) (pivnet.ImageReference, error)
	DeleteImageReference(productSlug string, releaseID int) (pivnet.ImageReference, error)
	AddImageReferenceToRelease(productSlug string, imageReferenceID int, releaseID int) error
	RemoveImageReferenceFromRelease(productSlug string, imageReferenceID int, releaseID int) error
}

type ImageReferenceClient struct {
	pivnetClient     PivnetClient
	eh               errorhandler.ErrorHandler
	format           string
	outputWriter     io.Writer
	logWriter        io.Writer
	printer          printer.Printer
	l                logger.Logger
}

func NewImageReferenceClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	logWriter io.Writer,
	printer printer.Printer,
	l logger.Logger,
) *ImageReferenceClient {
	return &ImageReferenceClient{
		pivnetClient:     pivnetClient,
		eh:               eh,
		format:           format,
		outputWriter:     outputWriter,
		logWriter:        logWriter,
		printer:          printer,
		l:                l,
	}
}

func (c *ImageReferenceClient) printImageReference(imageReference pivnet.ImageReference) error {
	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Image Path",
			"Digest",
		})

		imageReferenceAsString := []string{
			strconv.Itoa(imageReference.ID),
			imageReference.Name,
			imageReference.ImagePath,
			imageReference.Digest,
		}
		table.Append(imageReferenceAsString)
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(imageReference)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(imageReference)
	}

	return nil
}

func (c *ImageReferenceClient) Create(config pivnet.CreateImageReferenceConfig) error {
	imageReference, err := c.pivnetClient.CreateImageReference(config)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printImageReference(imageReference)
}

func (c *ImageReferenceClient) Delete(productSlug string, imageReferenceID int) error {
	imageReference, err := c.pivnetClient.DeleteImageReference(
		productSlug,
		imageReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		message := fmt.Sprintf(
			"Image reference %d deleted for %s",
			imageReferenceID,
			productSlug,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err := fmt.Fprintln(c.outputWriter, coloredMessage)

		return err
	}

	return c.printImageReference(imageReference)
}

func (c *ImageReferenceClient) AddToRelease(
	productSlug string,
	imageReferenceID int,
	releaseVersion string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.AddImageReferenceToRelease(
		productSlug,
		release.ID,
		imageReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		message := fmt.Sprintf(
			"Image reference %d added to %s/%s",
			imageReferenceID,
			productSlug,
			releaseVersion,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err = fmt.Fprintln(c.outputWriter, coloredMessage)
	}

	return nil
}

func (c *ImageReferenceClient) RemoveFromRelease(
	productSlug string,
	imageReferenceID int,
	releaseVersion string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.RemoveImageReferenceFromRelease(
		productSlug,
		release.ID,
		imageReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		message := fmt.Sprintf(
			"Image reference %d removed from %s/%s",
			imageReferenceID,
			productSlug,
			releaseVersion,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err = fmt.Fprintln(c.outputWriter, coloredMessage)
	}

	return nil
}