package commands

import "errors"

var ErrShowHelpMessage = errors.New("help command invoked")

type HelpCommand struct{}

func (command *HelpCommand) Execute(args []string) error {
	// Reset flags to avoid overriding defaults
	Pivnet.APIToken = ""
	Pivnet.Host = DefaultHost

	return ErrShowHelpMessage
}
