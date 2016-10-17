package commands

import "github.com/pivotal-cf/go-pivnet/cmd/pivnet/commands/curl"

type CurlCommand struct {
	Method string `long:"request" short:"X" description:"Custom method e.g. PATCH"`
	Data   string `long:"data" short:"d" description:"Request data e.g. '{\"foo\":\"bar\"}'"`
}

//go:generate counterfeiter . CurlClient
type CurlClient interface {
	MakeRequest(method string, body string, args []string) error
}

var NewCurlClient = func() CurlClient {
	return curl.NewCurlClient(
		NewPivnetClient(),
		ErrorHandler,
		Pivnet.Format,
		OutputWriter,
		Printer,
	)
}

func (command *CurlCommand) Execute(args []string) error {
	Init()

	return NewCurlClient().MakeRequest(command.Method, command.Data, args)
}
