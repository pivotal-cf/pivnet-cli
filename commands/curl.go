package commands

import "github.com/pivotal-cf/pivnet-cli/commands/curl"

type Args struct {
	URL string `positional-arg-name:"URL" description:"URL without host or API prefix e.g. /products/p-mysql/releases/3451"`
}

type CurlCommand struct {
	Method string `long:"request" short:"X" description:"Custom method e.g. PATCH"`
	Data   string `long:"data" short:"d" description:"Request data e.g. '{\"foo\":\"bar\"}'"`
	Args   Args   `positional-args:"yes" required:"true"`
}

//go:generate counterfeiter . CurlClient
type CurlClient interface {
	MakeRequest(method string, body string, url string) error
}

var NewCurlClient = func(client curl.PivnetClient) CurlClient {
	return curl.NewCurlClient(
		client,
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *CurlCommand) Execute(args []string) error {
	err := Init(true)
	if err != nil {
		return err
	}

	client := NewPivnetClient()
	err = Auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	return NewCurlClient(client).MakeRequest(command.Method, command.Data, command.Args.URL)
}
