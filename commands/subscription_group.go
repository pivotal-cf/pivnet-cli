package commands

import (
	"github.com/pivotal-cf/pivnet-cli/commands/subscriptiongroup"
)

type SubscriptionGroupsCommand struct{}

type SubscriptionGroupCommand struct {
	SubscriptionGroupID int `long:"subscription-group-id" short:"i" description:"Subscription group ID e.g. 1234" required:"true"`
}

type AddSubscriptionGroupMemberCommand struct {
	SubscriptionGroupID int    `long:"subscription-group-id" short:"i" description:"Subscription group ID e.g. 1234" required:"true"`
	MemberEmail    string `long:"member-email" description:"Email of member to add e.g. example@example.net" required:"true"`
	IsAdmin        string `long:"admin" description:"Whether the user should be an admin e.g. --admin=[true|false]"`
}

type RemoveSubscriptionGroupMemberCommand struct {
	SubscriptionGroupID int    `long:"subscription-group-id" short:"i" description:"Subscription group ID e.g. 1234" required:"true"`
	MemberEmail    string `long:"member-email" description:"Email of member to add e.g. example@example.net" required:"true"`
}

//go:generate counterfeiter . SubscriptionGroupClient
type SubscriptionGroupClient interface {
	List() error
	Get(subscriptionGroupID int) error
	AddMember(subscriptionGroupID int, memberEmail string, isAdmin string) error
	RemoveMember(subscriptionGroupID int, memberEmail string) error
}

var NewSubscriptionGroupClient = func(client subscriptiongroup.PivnetClient) SubscriptionGroupClient {
	return subscriptiongroup.NewSubscriptionGroupClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *SubscriptionGroupsCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewSubscriptionGroupClient(client).List()
}

func (command *SubscriptionGroupCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewSubscriptionGroupClient(client).Get(command.SubscriptionGroupID)
}

func (command *AddSubscriptionGroupMemberCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewSubscriptionGroupClient(client).AddMember(command.SubscriptionGroupID, command.MemberEmail, command.IsAdmin)
}

func (command *RemoveSubscriptionGroupMemberCommand) Execute([]string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewSubscriptionGroupClient(client).RemoveMember(command.SubscriptionGroupID, command.MemberEmail)
}
