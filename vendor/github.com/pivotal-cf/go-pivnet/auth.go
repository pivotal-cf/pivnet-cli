package pivnet

import "net/http"

type AuthService struct {
	client Client
}

func (e AuthService) Check() error {
	url := "/authentication"

	var response EULAsResponse
	_, _, err := e.client.MakeRequest(
		"GET",
		url,
		http.StatusOK,
		nil,
		&response,
	)
	if err != nil {
		return err
	}

	return nil
}
