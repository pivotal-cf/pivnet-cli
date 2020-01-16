package helmchartreference

import (
	"fmt"
	"github.com/pivotal-cf/pivnet-cli/ui"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf/go-pivnet/v3"
	"github.com/pivotal-cf/go-pivnet/v3/logger"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	HelmChartReferences(productSlug string) ([]pivnet.HelmChartReference, error)
	HelmChartReferencesForRelease(productSlug string, releaseID int) ([]pivnet.HelmChartReference, error)
	HelmChartReference(productSlug string, helmChartReferenceID int) (pivnet.HelmChartReference, error)
	HelmChartReferenceForRelease(productSlug string, releaseID int, helmChartReferenceID int) (pivnet.HelmChartReference, error)
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
	CreateHelmChartReference(config pivnet.CreateHelmChartReferenceConfig) (pivnet.HelmChartReference, error)
	DeleteHelmChartReference(productSlug string, releaseID int) (pivnet.HelmChartReference, error)
	AddHelmChartReferenceToRelease(productSlug string, helmChartReferenceID int, releaseID int) error
	RemoveHelmChartReferenceFromRelease(productSlug string, helmChartReferenceID int, releaseID int) error
	UpdateHelmChartReference(productSlug string, helmChartReference pivnet.HelmChartReference) (pivnet.HelmChartReference, error)
}

type HelmChartReferenceClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	logWriter    io.Writer
	printer      printer.Printer
	l            logger.Logger
}

func NewHelmChartReferenceClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	logWriter io.Writer,
	printer printer.Printer,
	l logger.Logger,
) *HelmChartReferenceClient {
	return &HelmChartReferenceClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		logWriter:    logWriter,
		printer:      printer,
		l:            l,
	}
}

func (c *HelmChartReferenceClient) Update(productSlug string, helmChartReferenceID int, description *string, docsURL *string, systemRequirements *[]string) error {
	helmChartReference, err := c.pivnetClient.HelmChartReference(
		productSlug,
		helmChartReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if description != nil {
		helmChartReference.Description = *description
	}

	if docsURL != nil {
		helmChartReference.DocsURL = *docsURL
	}

	if systemRequirements != nil {
		helmChartReference.SystemRequirements = *systemRequirements
	}

	updatedHelmChartReference, err := c.pivnetClient.UpdateHelmChartReference(productSlug, helmChartReference)
	if err != nil {
		return c.eh.HandleError(err)
	}
	return c.printHelmChartReference(updatedHelmChartReference)
}

func (c *HelmChartReferenceClient) List(productSlug string, releaseVersion string) error {
	if releaseVersion == "" {
		helmChartReferences, err := c.pivnetClient.HelmChartReferences(productSlug)
		if err != nil {
			return c.eh.HandleError(err)
		}

		return c.printHelmChartReferences(helmChartReferences)
	}

	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	helmChartReferences, err := c.pivnetClient.HelmChartReferencesForRelease(
		productSlug,
		release.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printHelmChartReferences(helmChartReferences)
}

func (c *HelmChartReferenceClient) printHelmChartReferences(helmChartReferences []pivnet.HelmChartReference) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Version",
		})

		for _, helmChartReference := range helmChartReferences {
			helmChartReferenceAsString := []string{
				strconv.Itoa(helmChartReference.ID),
				helmChartReference.Name,
				helmChartReference.Version,
			}
			table.Append(helmChartReferenceAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(helmChartReferences)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(helmChartReferences)
	}

	return nil
}

func (c *HelmChartReferenceClient) printHelmChartReference(helmChartReference pivnet.HelmChartReference) error {
	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Version",
		})

		helmChartReferenceAsString := []string{
			strconv.Itoa(helmChartReference.ID),
			helmChartReference.Name,
			helmChartReference.Version,
		}
		table.Append(helmChartReferenceAsString)
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(helmChartReference)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(helmChartReference)
	}

	return nil
}

func (c *HelmChartReferenceClient) Get(
	productSlug string,
	releaseVersion string,
	helmChartReferenceID int,
) error {
	if releaseVersion == "" {
		helmChartReference, err := c.pivnetClient.HelmChartReference(
			productSlug,
			helmChartReferenceID,
		)
		if err != nil {
			return c.eh.HandleError(err)
		}
		return c.printHelmChartReference(helmChartReference)
	}

	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	helmChartReference, err := c.pivnetClient.HelmChartReferenceForRelease(
		productSlug,
		release.ID,
		helmChartReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printHelmChartReference(helmChartReference)
}

func (c *HelmChartReferenceClient) Create(config pivnet.CreateHelmChartReferenceConfig) error {
	helmChartReference, err := c.pivnetClient.CreateHelmChartReference(config)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printHelmChartReference(helmChartReference)
}

func (c *HelmChartReferenceClient) Delete(productSlug string, helmChartReferenceID int) error {
	helmChartReference, err := c.pivnetClient.DeleteHelmChartReference(
		productSlug,
		helmChartReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		message := fmt.Sprintf(
			"Helm chart reference %d deleted for %s",
			helmChartReferenceID,
			productSlug,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err := fmt.Fprintln(c.outputWriter, coloredMessage)

		return err
	}

	return c.printHelmChartReference(helmChartReference)
}

func (c *HelmChartReferenceClient) AddToRelease(
	productSlug string,
	helmChartReferenceID int,
	releaseVersion string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.AddHelmChartReferenceToRelease(
		productSlug,
		release.ID,
		helmChartReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		message := fmt.Sprintf(
			"Helm chart reference %d added to %s/%s",
			helmChartReferenceID,
			productSlug,
			releaseVersion,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err = fmt.Fprintln(c.outputWriter, coloredMessage)
	}

	return nil
}

func (c *HelmChartReferenceClient) RemoveFromRelease(
	productSlug string,
	helmChartReferenceID int,
	releaseVersion string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.RemoveHelmChartReferenceFromRelease(
		productSlug,
		release.ID,
		helmChartReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		message := fmt.Sprintf(
			"Helm chart reference %d removed from %s/%s",
			helmChartReferenceID,
			productSlug,
			releaseVersion,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err = fmt.Fprintln(c.outputWriter, coloredMessage)
	}

	return nil
}
