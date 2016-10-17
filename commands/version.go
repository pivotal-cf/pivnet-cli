package commands

type VersionCommand struct{}

func (command *VersionCommand) Execute(args []string) error {
	Pivnet.VersionFunc()
	return nil
}
