package usergroup

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/errorhandler"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/printer"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	ReleaseForVersion(productSlug string, releaseVersion string) (pivnet.Release, error)
	UserGroups() ([]pivnet.UserGroup, error)
	UserGroupsForRelease(productSlug string, releaseID int) ([]pivnet.UserGroup, error)
	UserGroup(userGroupID int) (pivnet.UserGroup, error)
	CreateUserGroup(name string, description string, members []string) (pivnet.UserGroup, error)
	UpdateUserGroup(userGroup pivnet.UserGroup) (pivnet.UserGroup, error)
	DeleteUserGroup(userGroupID int) error
	AddUserGroup(productSlug string, releaseID int, userGroupID int) error
	RemoveUserGroup(productSlug string, releaseID int, userGroupID int) error
	AddMemberToGroup(userGroupID int, emailAddress string, admin bool) (pivnet.UserGroup, error)
	RemoveMemberFromGroup(userGroupID int, emailAddress string) (pivnet.UserGroup, error)
}

type UserGroupClient struct {
	pivnetClient PivnetClient
	eh           errorhandler.ErrorHandler
	format       string
	outputWriter io.Writer
	printer      printer.Printer
}

func NewUserGroupClient(
	pivnetClient PivnetClient,
	eh errorhandler.ErrorHandler,
	format string,
	outputWriter io.Writer,
	printer printer.Printer,
) *UserGroupClient {
	return &UserGroupClient{
		pivnetClient: pivnetClient,
		eh:           eh,
		format:       format,
		outputWriter: outputWriter,
		printer:      printer,
	}
}

func (c *UserGroupClient) List(productSlug string, releaseVersion string) error {
	if releaseVersion == "" {
		userGroups, err := c.pivnetClient.UserGroups()
		if err != nil {
			return c.eh.HandleError(err)
		}

		return c.printUserGroups(userGroups)
	}

	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	userGroups, err := c.pivnetClient.UserGroupsForRelease(
		productSlug,
		release.ID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printUserGroups(userGroups)
}

func (c *UserGroupClient) printUserGroups(userGroups []pivnet.UserGroup) error {
	switch c.format {

	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{"ID", "Name", "Description"})

		for _, userGroup := range userGroups {
			userGroupAsString := []string{
				strconv.Itoa(userGroup.ID),
				userGroup.Name,
				userGroup.Description,
			}
			table.Append(userGroupAsString)
		}
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(userGroups)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(userGroups)
	}

	return nil
}

func (c *UserGroupClient) Get(userGroupID int) error {
	userGroup, err := c.pivnetClient.UserGroup(userGroupID)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printUserGroup(userGroup)
}

func (c *UserGroupClient) printUserGroup(userGroup pivnet.UserGroup) error {
	switch c.format {
	case printer.PrintAsTable:
		table := tablewriter.NewWriter(c.outputWriter)
		table.SetHeader([]string{"ID", "Name", "Description", "Members"})

		userGroupAsString := []string{
			strconv.Itoa(userGroup.ID),
			userGroup.Name,
			userGroup.Description,
			strings.Join(userGroup.Members, ", "),
		}
		table.Append(userGroupAsString)
		table.Render()
		return nil
	case printer.PrintAsJSON:
		return c.printer.PrintJSON(userGroup)
	case printer.PrintAsYAML:
		return c.printer.PrintYAML(userGroup)
	}

	return nil
}

func (c *UserGroupClient) AddToRelease(
	productSlug string,
	releaseVersion string,
	userGroupID int,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.AddUserGroup(
		productSlug,
		release.ID,
		userGroupID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"user group %d added successfully to %s/%s\n",
			userGroupID,
			productSlug,
			releaseVersion,
		)
	}

	return nil
}

func (c *UserGroupClient) RemoveFromRelease(
	productSlug string,
	releaseVersion string,
	userGroupID int,
) error {
	release, err := c.pivnetClient.ReleaseForVersion(productSlug, releaseVersion)
	if err != nil {
		return c.eh.HandleError(err)
	}

	err = c.pivnetClient.RemoveUserGroup(
		productSlug,
		release.ID,
		userGroupID,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"user group %d removed successfully from %s/%s\n",
			userGroupID,
			productSlug,
			releaseVersion,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *UserGroupClient) Create(
	name string,
	description string,
	members []string,
) error {
	userGroup, err := c.pivnetClient.CreateUserGroup(
		name,
		description,
		members,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printUserGroup(userGroup)
}

func (c *UserGroupClient) Update(
	userGroupID int,
	name *string,
	description *string,
) error {
	userGroup, err := c.pivnetClient.UserGroup(userGroupID)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if name != nil {
		userGroup.Name = *name
	}

	if description != nil {
		userGroup.Description = *description
	}

	updated, err := c.pivnetClient.UpdateUserGroup(userGroup)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printUserGroup(updated)
}

func (c *UserGroupClient) Delete(userGroupID int) error {
	err := c.pivnetClient.DeleteUserGroup(userGroupID)
	if err != nil {
		return c.eh.HandleError(err)
	}

	if c.format == printer.PrintAsTable {
		_, err = fmt.Fprintf(
			c.outputWriter,
			"user group %d deleted successfully\n",
			userGroupID,
		)
	}

	return nil
}

func (c *UserGroupClient) AddUserGroupMember(
	userGroupID int,
	memberEmailAddress string,
	admin bool,
) error {
	userGroup, err := c.pivnetClient.AddMemberToGroup(
		userGroupID,
		memberEmailAddress,
		admin,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printUserGroup(userGroup)
}

func (c *UserGroupClient) RemoveUserGroupMember(
	userGroupID int,
	memberEmailAddress string,
) error {
	userGroup, err := c.pivnetClient.RemoveMemberFromGroup(
		userGroupID,
		memberEmailAddress,
	)
	if err != nil {
		return c.eh.HandleError(err)
	}

	return c.printUserGroup(userGroup)
}
