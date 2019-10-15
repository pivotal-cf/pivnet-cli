package companygroup

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/pivnet-cli/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/printer"
	"io"
	"strconv"
	"strings"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	CompanyGroups() ([]pivnet.CompanyGroup, error)
	CompanyGroup(companyGroupID int) (pivnet.CompanyGroup, error)
	AddCompanyGroupMember(companyGroupID int, emailAddress string, isAdmin string) (pivnet.CompanyGroup, error)
	RemoveCompanyGroupMember(companyGroupID int, emailAddress string) (pivnet.CompanyGroup, error)
}

type CompanyGroupClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewCompanyGroupClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *CompanyGroupClient {
	return &CompanyGroupClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *CompanyGroupClient) List() error {
	companyGroups, err := c.pivnetClient.CompanyGroups()
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printCompanyGroups(companyGroups)
}

func (c *CompanyGroupClient) printCompanyGroups(companyGroups []pivnet.CompanyGroup) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{"ID", "Name"})

		for _, companyGroup := range companyGroups {
			companyGroupAsString := []string{
				strconv.Itoa(companyGroup.ID),
				companyGroup.Name,
			}
			table.Append(companyGroupAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(companyGroups)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(companyGroups)
	}

	return nil
}

func (c *CompanyGroupClient) Get(companyGroupID int) error {
	companyGroup, err := c.pivnetClient.CompanyGroup(companyGroupID)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printCompanyGroup(companyGroup)
}

func (c *CompanyGroupClient) AddMember(companyGroupID int, memberEmail string, isAdmin string) error {
	companyGroup, err := c.pivnetClient.AddCompanyGroupMember(companyGroupID, memberEmail, isAdmin)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printCompanyGroup(companyGroup)
}

func (c *CompanyGroupClient) RemoveMember(companyGroupID int, memberEmail string) error {
	companyGroup, err := c.pivnetClient.RemoveCompanyGroupMember(companyGroupID, memberEmail)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printCompanyGroup(companyGroup)
}

func (c *CompanyGroupClient) printCompanyGroup(companyGroup pivnet.CompanyGroup) error {
	switch c.format {
	case printer.PrintAsTable:
		b := color.New(color.Bold).SprintFunc()
		fmt.Println(b("ID: "), companyGroup.ID)
		fmt.Println(b("Company Name: "), companyGroup.Name)

		var entitlements []string
		for _, entitlement := range companyGroup.Entitlements {
			entitlements = append(entitlements, entitlement.Name)
		}
		fmt.Println(b("Entitlements: "), strings.Join(entitlements, ", "))
		fmt.Println("")
		fmt.Println(b("Members/Pending Invitations:"))
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{"Name", "Email", "Is Admin", "Pending"})

		for _, member := range companyGroup.Members {
			companyGroupAsString := []string{
				member.Name,
				member.Email,
				strconv.FormatBool(member.IsAdmin),
				"x",
			}

			table.Append(companyGroupAsString)
		}

		for _, pendingEmail := range companyGroup.PendingInvitations {
			companyGroupAsString := []string{
				"",
				pendingEmail,
				"",
				"\xE2\x9C\x94",
			}

			table.Append(companyGroupAsString)
		}

		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(companyGroup)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(companyGroup)
	}

	return nil
}
