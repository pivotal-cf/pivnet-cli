package commands

import "github.com/pivotal-cf/go-pivnet/cmd/pivnet/commands/usergroup"

type UserGroupsCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1"`
}

type UserGroupCommand struct {
	UserGroupID int `long:"user-group-id" short:"i" description:"User group ID e.g. 1234" required:"true"`
}

type CreateUserGroupCommand struct {
	Name        string   `long:"name" description:"Name e.g. all_users" required:"true"`
	Description string   `long:"description" description:"Description e.g. 'All users in the world'" required:"true"`
	Members     []string `long:"member" description:"Email addresses of members to be added"`
}

type UpdateUserGroupCommand struct {
	UserGroupID int     `long:"user-group-id" short:"i" description:"User group ID e.g. 1234" required:"true"`
	Name        *string `long:"name" description:"Name e.g. all_users"`
	Description *string `long:"description" description:"Description e.g. 'All users in the world'"`
}

type AddUserGroupCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	UserGroupID    int    `long:"user-group-id" short:"i" description:"User Group ID e.g. 1234" required:"true"`
}

type RemoveUserGroupCommand struct {
	ProductSlug    string `long:"product-slug" short:"p" description:"Product slug e.g. p-mysql" required:"true"`
	ReleaseVersion string `long:"release-version" short:"r" description:"Release version e.g. 0.1.2-rc1" required:"true"`
	UserGroupID    int    `long:"user-group-id" short:"i" description:"User Group ID e.g. 1234" required:"true"`
}

type DeleteUserGroupCommand struct {
	UserGroupID int `long:"user-group-id" short:"i" description:"User group ID e.g. 1234" required:"true"`
}

type AddUserGroupMemberCommand struct {
	UserGroupID        int    `long:"user-group-id" short:"i" description:"User group ID e.g. 1234" required:"true"`
	MemberEmailAddress string `long:"member-email" description:"Member email address e.g. 1234" required:"true"`
	Admin              bool   `long:"admin" description:"Whether the user should be an admin"`
}

type RemoveUserGroupMemberCommand struct {
	UserGroupID        int    `long:"user-group-id" short:"i" description:"User group ID e.g. 1234" required:"true"`
	MemberEmailAddress string `long:"member-email" description:"Member email address e.g. 1234" required:"true"`
}

//go:generate counterfeiter . UserGroupClient
type UserGroupClient interface {
	List(productSlug string, releaseVersion string) error
	Get(userGroupID int) error
	Create(name string, description string, members []string) error
	Update(userGroupID int, name *string, description *string) error
	AddToRelease(productSlug string, releaseVersion string, userGroupID int) error
	Delete(userGroupID int) error
	RemoveFromRelease(productSlug string, releaseVersion string, userGroupID int) error
	AddUserGroupMember(userGroupID int, memberEmailAddress string, admin bool) error
	RemoveUserGroupMember(userGroupID int, memberEmailAddress string) error
}

var NewUserGroupClient = func() UserGroupClient {
	return usergroup.NewUserGroupClient(
		NewPivnetClient(),
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *UserGroupCommand) Execute([]string) error {
	Init()

	return NewUserGroupClient().Get(command.UserGroupID)
}

func (command *UserGroupsCommand) Execute([]string) error {
	Init()

	return NewUserGroupClient().List(command.ProductSlug, command.ReleaseVersion)
}

func (command *CreateUserGroupCommand) Execute([]string) error {
	Init()

	return NewUserGroupClient().Create(
		command.Name,
		command.Description,
		command.Members,
	)
}

func (command *DeleteUserGroupCommand) Execute([]string) error {
	Init()

	return NewUserGroupClient().Delete(command.UserGroupID)
}

func (command *UpdateUserGroupCommand) Execute([]string) error {
	Init()

	return NewUserGroupClient().Update(
		command.UserGroupID,
		command.Name,
		command.Description,
	)
}

func (command *AddUserGroupCommand) Execute([]string) error {
	Init()

	return NewUserGroupClient().AddToRelease(
		command.ProductSlug,
		command.ReleaseVersion,
		command.UserGroupID,
	)
}

func (command *RemoveUserGroupCommand) Execute([]string) error {
	Init()

	return NewUserGroupClient().RemoveFromRelease(
		command.ProductSlug,
		command.ReleaseVersion,
		command.UserGroupID,
	)
}

func (command *AddUserGroupMemberCommand) Execute([]string) error {
	Init()

	return NewUserGroupClient().AddUserGroupMember(
		command.UserGroupID,
		command.MemberEmailAddress,
		command.Admin,
	)
}

func (command *RemoveUserGroupMemberCommand) Execute([]string) error {
	Init()

	return NewUserGroupClient().RemoveUserGroupMember(
		command.UserGroupID,
		command.MemberEmailAddress,
	)
}
