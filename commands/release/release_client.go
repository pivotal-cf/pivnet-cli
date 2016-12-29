package release

import (
	"fmt"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
	"github.com/pivotal-cf/pivnet-cli/ui"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	ReleasesForProductSlug(productSlug string) ([]pivnet.Release, error)
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
	CreateRelease(config pivnet.CreateReleaseConfig) (pivnet.Release, error)
	DeleteRelease(productSlug string, release pivnet.Release) error
	EULAs() ([]pivnet.EULA, error)
	ReleaseTypes() ([]pivnet.ReleaseType, error)
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
			"Updated At",
		})

		for _, release := range releases {
			releaseAsString := []string{
				strconv.Itoa(release.ID),
				release.Version,
				release.Description,
				release.UpdatedAt,
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

	return c.printRelease(release)
}

func (c *ReleaseClient) printRelease(release pivnet.Release) error {
	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{
			"ID",
			"Version",
			"Description",
			"Updated At",
		})

		releaseAsString := []string{
			strconv.Itoa(release.ID),
			release.Version,
			release.Description,
			release.UpdatedAt,
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
	err := c.validateEULA(eulaSlug)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.validateReleaseType(releaseType)
	if err != nil {
		return c.eh.HandleError(err)
	}

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

	return c.printRelease(release)
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
		message := fmt.Sprintf(
			"Release %d deleted for %s",
			release.ID,
			productSlug,
		)
		coloredMessage := ui.SuccessColor.SprintFunc()(message)

		_, err := fmt.Fprintln(c.outputWriter, coloredMessage)

		return err
	}

	return nil
}

func (c *ReleaseClient) validateEULA(eulaSlug string) error {
	eulas, err := c.pivnetClient.EULAs()
	if err != nil {
		return err
	}

	eulaSlugs := eulaSlugsFromEULAs(eulas)

	if !stringsContains(eulaSlugs, eulaSlug) {
		err := fmt.Errorf(
			"provided EULA Slug: '%s' must be one of: %v",
			eulaSlug,
			eulaSlugs,
		)
		return err
	}

	return nil
}

func (c *ReleaseClient) validateReleaseType(releaseType string) error {
	releaseTypes, err := c.pivnetClient.ReleaseTypes()
	if err != nil {
		return err
	}

	releaseTypesStrings := releaseTypesToStrings(releaseTypes)

	if !stringsContains(releaseTypesStrings, releaseType) {
		err := fmt.Errorf(
			"provided release type: '%s' must be one of: %v",
			releaseType,
			releaseTypes,
		)
		return err
	}

	return nil
}

func releaseTypesToStrings(releaseTypes []pivnet.ReleaseType) []string {
	var releaseTypeStrings []string
	for _, r := range releaseTypes {
		releaseTypeStrings = append(releaseTypeStrings, string(r))
	}

	return releaseTypeStrings
}

func eulaSlugsFromEULAs(eulas []pivnet.EULA) []string {
	var eulaSlugs []string
	for _, e := range eulas {
		eulaSlugs = append(eulaSlugs, e.Slug)
	}

	return eulaSlugs
}

func stringsContains(strings []string, val string) bool {
	for _, s := range strings {
		if s == val {
			return true
		}
	}
	return false
}
