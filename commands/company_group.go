package commands

import (
	"github.com/pivotal-cf/pivnet-cli/commands/companygroup"
)

type CompanyGroupsCommand struct{}

//go:generate counterfeiter . CompanyGroupClient
type CompanyGroupClient interface {
	List() error
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
