package auth

import (
	"fmt"

	"github.com/pivotal-cf/pivnet-cli/v2/errorhandler"
)

//go:generate counterfeiter . AuthClient
type AuthClient interface {
	Auth() (bool, error)
}

type Authenticator struct {
	eh errorhandler.ErrorHandler
}

func NewAuthenticator(eh errorhandler.ErrorHandler) *Authenticator {
	return &Authenticator{
		eh: eh,
	}
}

func (a *Authenticator) AuthenticateClient(client AuthClient) error {
	ok, err := client.Auth()
	if err != nil {
		return a.eh.HandleError(err)
	}

	if !ok {
		err = fmt.Errorf("Credentials rejected - please login again")
		return a.eh.HandleError(err)
	}

	return nil
}
