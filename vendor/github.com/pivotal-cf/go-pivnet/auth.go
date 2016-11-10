package pivnet

import "net/http"

type AuthService struct {
	client Client
}

func (e AuthService) Check() error {
	url := "/authentication"

	resp, err := e.client.MakeRequest(
		"GET",
		url,
		http.StatusOK,
		nil,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
