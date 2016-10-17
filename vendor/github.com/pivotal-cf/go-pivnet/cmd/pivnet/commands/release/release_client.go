package release

import (
	"fmt"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/errorhandler"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	ReleasesForProductSlug(productSlug string) ([]pivnet.Release, error)
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
	CreateRelease(config pivnet.CreateReleaseConfig) (pivnet.Release, error)
	DeleteRelease(productSlug string, release pivnet.Release) error
	ReleaseFingerprint(productSlug string, releaseID int) (string, error)
}

type ReleaseClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewReleaseClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *ReleaseClient {
	return &ReleaseClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *ReleaseClient) List(productSlug string) error {
	releases, err := c.pivnetClient.ReleasesForProductSlug(productSlug)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printReleases(releases)
}

func (c *ReleaseClient) printReleases(releases []pivnet.Release) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Version",
			"Description",
		})

		for _, release := range releases {
			releaseAsString := []string{
				strconv.Itoa(release.ID),
				release.Version,
				release.Description,
			}
			table.Append(releaseAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(releases)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(releases)
	}

	return nil
}

func (c *ReleaseClient) Get(
	productSlug string,
	releaseVersion string,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(
		productSlug,
		releaseVersion,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	fingerprint, err := c.pivnetClient.ReleaseFingerprint(productSlug, release.ID)
	if err != nil {
		return c.eh.HandleError(err)
	}

	r := CLIRelease{
		release,
		fingerprint,
	}

	return c.printRelease(r)
}

func (c *ReleaseClient) printRelease(release CLIRelease) error {
	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Version",
			"Description",
			"Fingerprint",
		})

		releaseAsString := []string{
			strconv.Itoa(release.ID),
			release.Version,
			release.Description,
			release.Fingerprint,
		}
		table.Append(releaseAsString)
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(release)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(release)
	}

	return nil
}

func (c *ReleaseClient) Create(
	productSlug string,
	releaseVersion string,
	releaseType string,
	eulaSlug string,
) error {
	newReleaseConfig := pivnet.CreateReleaseConfig{
		ProductSlug: productSlug,
		Version:     releaseVersion,
		ReleaseType: releaseType,
		EULASlug:    eulaSlug,
	}

	release, err := c.pivnetClient.CreateRelease(newReleaseConfig)
	if err != nil {
		return c.eh.HandleError(err)
	}

	fingerprint, err := c.pivnetClient.ReleaseFingerprint(productSlug, release.ID)
	if err != nil {
		return c.eh.HandleError(err)
	}

	r := CLIRelease{
		release,
		fingerprint,
	}

	return c.printRelease(r)
}

func (c *ReleaseClient) Delete(productSlug string, releaseVersion string) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.DeleteRelease(
		productSlug,
		release,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"release %d deleted successfully for %s\n",
			release.ID,
			productSlug,
		)
	}

	return nil
}

type CLIRelease struct {
	pivnet.Release `yaml:",inline"`
	Fingerprint    string `json:"fingerprint,omitempty"`
}
