package subscriptiongroup

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf/go-pivnet/v7"
	"github.com/pivotal-cf/pivnet-cli/v3/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/v3/printer"
	"io"
	"strconv"
	"strings"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	SubscriptionGroups() ([]pivnet.SubscriptionGroup, error)
	SubscriptionGroup(subscriptionGroupID int) (pivnet.SubscriptionGroup, error)
	AddSubscriptionGroupMember(subscriptionGroupID int, emailAddress string, isAdmin string) (pivnet.SubscriptionGroup, error)
	RemoveSubscriptionGroupMember(subscriptionGroupID int, emailAddress string) (pivnet.SubscriptionGroup, error)
}

type SubscriptionGroupClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewSubscriptionGroupClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *SubscriptionGroupClient {
	return &SubscriptionGroupClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *SubscriptionGroupClient) List() error {
	subscriptionGroups, err := c.pivnetClient.SubscriptionGroups()
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printSubscriptionGroups(subscriptionGroups)
}

func (c *SubscriptionGroupClient) printSubscriptionGroups(subscriptionGroups []pivnet.SubscriptionGroup) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{"ID", "Name"})

		for _, subscriptionGroup := range subscriptionGroups {
			subscriptionGroupAsString := []string{
				strconv.Itoa(subscriptionGroup.ID),
				subscriptionGroup.Name,
			}
			table.Append(subscriptionGroupAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(subscriptionGroups)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(subscriptionGroups)
	}

	return nil
}

func (c *SubscriptionGroupClient) Get(subscriptionGroupID int) error {
	subscriptionGroup, err := c.pivnetClient.SubscriptionGroup(subscriptionGroupID)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printSubscriptionGroup(subscriptionGroup)
}

func (c *SubscriptionGroupClient) AddMember(subscriptionGroupID int, memberEmail string, isAdmin string) error {
	subscriptionGroup, err := c.pivnetClient.AddSubscriptionGroupMember(subscriptionGroupID, memberEmail, isAdmin)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printSubscriptionGroup(subscriptionGroup)
}

func (c *SubscriptionGroupClient) RemoveMember(subscriptionGroupID int, memberEmail string) error {
	subscriptionGroup, err := c.pivnetClient.RemoveSubscriptionGroupMember(subscriptionGroupID, memberEmail)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printSubscriptionGroup(subscriptionGroup)
}

func (c *SubscriptionGroupClient) printSubscriptionGroup(subscriptionGroup pivnet.SubscriptionGroup) error {
	switch c.format {
	case printer.PrintAsTable:
		b := color.New(color.Bold).SprintFunc()
		fmt.Println(b("ID: "), subscriptionGroup.ID)
		fmt.Println(b("Subscription Name: "), subscriptionGroup.Name)

		var subscriptions []string
		for _, subscription := range subscriptionGroup.Subscriptions {
			subscriptions = append(subscriptions, subscription.Name)
		}
		fmt.Println(b("Subscriptions: "), strings.Join(subscriptions, ", "))
		fmt.Println("")
		fmt.Println(b("Members/Pending Invitations:"))
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{"Name", "Email", "Is Admin", "Pending"})

		for _, member := range subscriptionGroup.Members {
			subscriptionGroupAsString := []string{
				member.Name,
				member.Email,
				strconv.FormatBool(member.IsAdmin),
				"x",
			}

			table.Append(subscriptionGroupAsString)
		}

		for _, pendingEmail := range subscriptionGroup.PendingInvitations {
			subscriptionGroupAsString := []string{
				"",
				pendingEmail,
				"",
				"\xE2\x9C\x94",
			}

			table.Append(subscriptionGroupAsString)
		}

		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(subscriptionGroup)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(subscriptionGroup)
	}

	return nil
}
