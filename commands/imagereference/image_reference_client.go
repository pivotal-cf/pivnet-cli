package imagereference

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/pivotal-cf/pivnet-cli/ui"

	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf/go-pivnet/v4"
	"github.com/pivotal-cf/go-pivnet/v4/logger"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	ImageReferences(productSlug string) ([]pivnet.ImageReference, error)
	ImageReferencesForRelease(productSlug string, releaseID int) ([]pivnet.ImageReference, error)
	ImageReferencesForDigest(productSlug string, imageReferenceDigest string) ([]pivnet.ImageReference, error)
	ImageReference(productSlug string, imageReferenceID int) (pivnet.ImageReference, error)
	ImageReferenceForRelease(productSlug string, releaseID int, imageReferenceID int) (pivnet.ImageReference, error)
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
	CreateImageReference(config pivnet.CreateImageReferenceConfig) (pivnet.ImageReference, error)
	DeleteImageReference(productSlug string, releaseID int) (pivnet.ImageReference, error)
	AddImageReferenceToRelease(productSlug string, imageReferenceID int, releaseID int) error
	RemoveImageReferenceFromRelease(productSlug string, imageReferenceID int, releaseID int) error
	UpdateImageReference(productSlug string, imageReference pivnet.ImageReference) (pivnet.ImageReference, error)
}

type ImageReferenceClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	logWriter    io.Writer
	printer      printer.Printer
	l            logger.Logger
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
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		logWriter:    logWriter,
		printer:      printer,
		l:            l,
	}
}

func (c *ImageReferenceClient) Update(productSlug string, imageReferenceID int, name *string, description *string, docsURL *string, systemRequirements *[]string) error {
	imageReference, err := c.pivnetClient.ImageReference(
		productSlug,
		imageReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if name != nil {
		imageReference.Name = *name
	}

	if description != nil {
		imageReference.Description = *description
	}

	if docsURL != nil {
		imageReference.DocsURL = *docsURL
	}

	if systemRequirements != nil {
		imageReference.SystemRequirements = *systemRequirements
	}

	updatedImageReference, err := c.pivnetClient.UpdateImageReference(productSlug, imageReference)
	if err != nil {
		return c.eh.HandleError(err)
	}
	return c.printImageReference(updatedImageReference)
}

func (c *ImageReferenceClient) List(productSlug string, releaseVersion string, imageDigest string) error {
	if imageDigest != "" {
		imageReferences, err := c.pivnetClient.ImageReferencesForDigest(productSlug, imageDigest)
		if err != nil {
			return c.eh.HandleError(err)
		}

		return c.printImageReferencesForDigest(imageReferences)
	}

	if releaseVersion == "" {
		imageReferences, err := c.pivnetClient.ImageReferences(productSlug)
		if err != nil {
			return c.eh.HandleError(err)
		}

		return c.printImageReferences(imageReferences)
	}

	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	imageReferences, err := c.pivnetClient.ImageReferencesForRelease(
		productSlug,
		release.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printImageReferences(imageReferences)
}

func (c *ImageReferenceClient) printImageReferences(imageReferences []pivnet.ImageReference) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Image Path",
			"Digest",
		})

		for _, imageReference := range imageReferences {
			imageReferenceAsString := []string{
				strconv.Itoa(imageReference.ID),
				imageReference.Name,
				imageReference.ImagePath,
				imageReference.Digest,
			}
			table.Append(imageReferenceAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(imageReferences)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(imageReferences)
	}

	return nil
}

func (c *ImageReferenceClient) printImageReferencesForDigest(imageReferences []pivnet.ImageReference) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Image Path",
			"releases",
		})

		for _, imageReference := range imageReferences {
			imageReferenceAsString := []string{
				strconv.Itoa(imageReference.ID),
				imageReference.Name,
				imageReference.ImagePath,
				strings.Join(imageReference.ReleaseVersions, ","),
			}
			table.Append(imageReferenceAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(imageReferences)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(imageReferences)
	}

	return nil
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

func (c *ImageReferenceClient) Get(
	productSlug string,
	releaseVersion string,
	imageReferenceID int,
) error {
	if releaseVersion == "" {
		imageReference, err := c.pivnetClient.ImageReference(
			productSlug,
			imageReferenceID,
		)
		if err != nil {
			return c.eh.HandleError(err)
		}
		return c.printImageReference(imageReference)
	}

	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	imageReference, err := c.pivnetClient.ImageReferenceForRelease(
		productSlug,
		release.ID,
		imageReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printImageReference(imageReference)
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
