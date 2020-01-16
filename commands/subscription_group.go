package commands

import (
	"github.com/pivotal-cf/pivnet-cli/commands/companygroup"
)

type CompanyGroupsCommand struct{}

type CompanyGroupCommand struct {
	CompanyGroupID int `long:"company-group-id" short:"i" description:"Company group ID e.g. 1234" required:"true"`
}

type AddCompanyGroupMemberCommand struct {
	CompanyGroupID int    `long:"company-group-id" short:"i" description:"Company group ID e.g. 1234" required:"true"`
	MemberEmail    string `long:"member-email" description:"Email of member to add e.g. example@example.net" required:"true"`
	IsAdmin        string `long:"admin" description:"Whether the user should be an admin e.g. --admin=[true|false]"`
}

type RemoveCompanyGroupMemberCommand struct {
	CompanyGroupID int    `long:"company-group-id" short:"i" description:"Company group ID e.g. 1234" required:"true"`
	MemberEmail    string `long:"member-email" description:"Email of member to add e.g. example@example.net" required:"true"`
}

//go:generate counterfeiter . CompanyGroupClient
type CompanyGroupClient interface {
	List() error
	Get(companyGroupID int) error
	AddMember(companyGroupID int, memberEmail string, isAdmin string) error
	RemoveMember(companyGroupID int, memberEmail string) error
}

var NewCompanyGroupClient = func(client companygroup.PivnetClient) CompanyGroupClient {
	return companygroup.NewCompanyGroupClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *CompanyGroupsCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewCompanyGroupClient(client).List()
}

func (command *CompanyGroupCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewCompanyGroupClient(client).Get(command.CompanyGroupID)
}

func (command *AddCompanyGroupMemberCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewCompanyGroupClient(client).AddMember(command.CompanyGroupID, command.MemberEmail, command.IsAdmin)
}

func (command *RemoveCompanyGroupMemberCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewCompanyGroupClient(client).RemoveMember(command.CompanyGroupID, command.MemberEmail)
}
