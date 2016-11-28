package pivnet

import "net/http"

type AuthService struct {
	client Client
}

// Check returns:
// true,nil if the auth attempt was succesful,
// false,nil if the auth attempt failed for 401 or 403,
// false,err if the auth attempt failed for any other reason.
// It is guaranteed never to return true,err.
func (e AuthService) Check() (bool, error) {
	url := "/authentication"

	resp, err := e.client.MakeRequest(
		"GET",
		url,
		0,
		nil,
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusUnauthorized:
		fallthrough
	case http.StatusForbidden:
		return false, nil
	default:
		return false, e.client.handleUnexpectedResponse(resp)
	}
}
