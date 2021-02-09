package artifactreference

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/pivotal-cf/pivnet-cli/v2/ui"

	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf/go-pivnet/v7"
	"github.com/pivotal-cf/go-pivnet/v7/logger"
	"github.com/pivotal-cf/pivnet-cli/v2/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/v2/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	ArtifactReferences(productSlug string) ([]pivnet.ArtifactReference, error)
	ArtifactReferencesForRelease(productSlug string, releaseID int) ([]pivnet.ArtifactReference, error)
	ArtifactReferencesForDigest(productSlug string, artifactReferenceDigest string) ([]pivnet.ArtifactReference, error)
	ArtifactReference(productSlug string, artifactReferenceID int) (pivnet.ArtifactReference, error)
	ArtifactReferenceForRelease(productSlug string, releaseID int, artifactReferenceID int) (pivnet.ArtifactReference, error)
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
	CreateArtifactReference(config pivnet.CreateArtifactReferenceConfig) (pivnet.ArtifactReference, error)
	DeleteArtifactReference(productSlug string, releaseID int) (pivnet.ArtifactReference, error)
	AddArtifactReferenceToRelease(productSlug string, artifactReferenceID int, releaseID int) error
	RemoveArtifactReferenceFromRelease(productSlug string, artifactReferenceID int, releaseID int) error
	UpdateArtifactReference(productSlug string, artifactReference pivnet.ArtifactReference) (pivnet.ArtifactReference, error)
}

type ArtifactReferenceClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	logWriter    io.Writer
	printer      printer.Printer
	l            logger.Logger
}

func NewArtifactReferenceClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	logWriter io.Writer,
	printer printer.Printer,
	l logger.Logger,
) *ArtifactReferenceClient {
	return &ArtifactReferenceClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		logWriter:    logWriter,
		printer:      printer,
		l:            l,
	}
}

func (c *ArtifactReferenceClient) Update(productSlug string, artifactReferenceID int, name *string, description *string, docsURL *string, systemRequirements *[]string) error {
	artifactReference, err := c.pivnetClient.ArtifactReference(
		productSlug,
		artifactReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if name != nil {
		artifactReference.Name = *name
	}

	if description != nil {
		artifactReference.Description = *description
	}

	if docsURL != nil {
		artifactReference.DocsURL = *docsURL
	}

	if systemRequirements != nil {
		artifactReference.SystemRequirements = *systemRequirements
	}

	updatedArtifactReference, err := c.pivnetClient.UpdateArtifactReference(productSlug, artifactReference)
	if err != nil {
		return c.eh.HandleError(err)
	}
	return c.printArtifactReference(updatedArtifactReference)
}

func (c *ArtifactReferenceClient) List(productSlug string, releaseVersion string, artifactDigest string) error {
	if artifactDigest != "" {
		artifactReferences, err := c.pivnetClient.ArtifactReferencesForDigest(productSlug, artifactDigest)
		if err != nil {
			return c.eh.HandleError(err)
		}

		return c.printArtifactReferencesForDigest(artifactReferences)
	}

	if releaseVersion == "" {
		artifactReferences, err := c.pivnetClient.ArtifactReferences(productSlug)
		if err != nil {
			return c.eh.HandleError(err)
		}

		return c.printArtifactReferences(artifactReferences)
	}

	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	artifactReferences, err := c.pivnetClient.ArtifactReferencesForRelease(
		productSlug,
		release.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printArtifactReferences(artifactReferences)
}

func (c *ArtifactReferenceClient) printArtifactReferences(artifactReferences []pivnet.ArtifactReference) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Artifact Path",
			"Digest",
		})

		for _, artifactReference := range artifactReferences {
			artifactReferenceAsString := []string{
				strconv.Itoa(artifactReference.ID),
				artifactReference.Name,
				artifactReference.ArtifactPath,
				artifactReference.Digest,
			}
			table.Append(artifactReferenceAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(artifactReferences)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(artifactReferences)
	}

	return nil
}

func (c *ArtifactReferenceClient) printArtifactReferencesForDigest(artifactReferences []pivnet.ArtifactReference) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Artifact Path",
			"releases",
		})

		for _, artifactReference := range artifactReferences {
			artifactReferenceAsString := []string{
				strconv.Itoa(artifactReference.ID),
				artifactReference.Name,
				artifactReference.ArtifactPath,
				strings.Join(artifactReference.ReleaseVersions, ","),
			}
			table.Append(artifactReferenceAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(artifactReferences)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(artifactReferences)
	}

	return nil
}

func (c *ArtifactReferenceClient) printArtifactReference(artifactReference pivnet.ArtifactReference) error {
	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Artifact Path",
			"Digest",
		})

		artifactReferenceAsString := []string{
			strconv.Itoa(artifactReference.ID),
			artifactReference.Name,
			artifactReference.ArtifactPath,
			artifactReference.Digest,
		}
		table.Append(artifactReferenceAsString)
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(artifactReference)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(artifactReference)
	}

	return nil
}

func (c *ArtifactReferenceClient) Get(
	productSlug string,
	releaseVersion string,
	artifactReferenceID int,
) error {
	if releaseVersion == "" {
		artifactReference, err := c.pivnetClient.ArtifactReference(
			productSlug,
			artifactReferenceID,
		)
		if err != nil {
			return c.eh.HandleError(err)
		}
		return c.printArtifactReference(artifactReference)
	}

	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	artifactReference, err := c.pivnetClient.ArtifactReferenceForRelease(
		productSlug,
		release.ID,
		artifactReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printArtifactReference(artifactReference)
}

func (c *ArtifactReferenceClient) Create(config pivnet.CreateArtifactReferenceConfig) error {
	artifactReference, err := c.pivnetClient.CreateArtifactReference(config)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printArtifactReference(artifactReference)
}

func (c *ArtifactReferenceClient) Delete(productSlug string, artifactReferenceID int) error {
	artifactReference, err := c.pivnetClient.DeleteArtifactReference(
		productSlug,
		artifactReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		message := fmt.Sprintf(
			"Artifact reference %d deleted for %s",
			artifactReferenceID,
			productSlug,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err := fmt.Fprintln(c.outputWriter, coloredMessage)

		return err
	}

	return c.printArtifactReference(artifactReference)
}

func (c *ArtifactReferenceClient) AddToRelease(
	productSlug string,
	artifactReferenceID int,
	releaseVersion string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.AddArtifactReferenceToRelease(
		productSlug,
		release.ID,
		artifactReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		message := fmt.Sprintf(
			"Artifact reference %d added to %s/%s",
			artifactReferenceID,
			productSlug,
			releaseVersion,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err = fmt.Fprintln(c.outputWriter, coloredMessage)
	}

	return nil
}

func (c *ArtifactReferenceClient) RemoveFromRelease(
	productSlug string,
	artifactReferenceID int,
	releaseVersion string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.RemoveArtifactReferenceFromRelease(
		productSlug,
		release.ID,
		artifactReferenceID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		message := fmt.Sprintf(
			"Artifact reference %d removed from %s/%s",
			artifactReferenceID,
			productSlug,
			releaseVersion,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err = fmt.Fprintln(c.outputWriter, coloredMessage)
	}

	return nil
}
