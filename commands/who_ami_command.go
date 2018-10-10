package commands

import "fmt"

type WhoAmICommand struct {
}

func (command *WhoAmICommand) Execute([]string) error {
	Init(false)

	var name = Pivnet.ProfileName
	var profile, err = RC.ProfileForName(name)

	if err != nil || profile == nil {
		fmt.Fprintln(OutputWriter, "No user is logged in")
	} else {
		var host = profile.Host
		output := fmt.Sprintf("Username: %v, Host: %v", name, host)
		fmt.Fprintln(OutputWriter, output)
	}

	return nil
}
