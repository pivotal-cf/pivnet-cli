package main

import (
	"fmt"
	"github.com/pivotal-cf/pivnet-cli/v2/rc"
	"github.com/pivotal-cf/pivnet-cli/v2/rc/filesystem"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/pivotal-cf/pivnet-cli/v2/commands"
	"github.com/pivotal-cf/pivnet-cli/v2/errorhandler"
	"github.com/pivotal-cf/pivnet-cli/v2/version"
	"github.com/pivotal-cf/pivnet-cli/v2/hostwarning"
)

var (
	// buildVersion is deliberately left uninitialized so it can be set at compile-time
	buildVersion string
)

func main() {
	if buildVersion == "" {
		version.Version = "dev"
	} else {
		version.Version = buildVersion
	}

	parser := flags.NewParser(&commands.Pivnet, flags.HelpFlag)

	if len(os.Args) > 1 && os.Args[1] == "manpage" {
		parser.WriteManPage(os.Stdout)
		return
	}

	_, err := parser.Parse()
	if err != nil {
		if err == commands.ErrShowHelpMessage {
			helpParser := flags.NewParser(&commands.Pivnet, flags.HelpFlag)
			helpParser.NamespaceDelimiter = "-"
			_, _ = helpParser.ParseArgs([]string{"-h"})
			helpParser.WriteHelp(os.Stdout)
			os.Exit(0)
		}

		// Do not consider the built-in help an error
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp {
				fmt.Fprintln(os.Stdout, err.Error())
				os.Exit(0)
			}
		}

		if err == errorhandler.ErrAlreadyHandled {
			os.Exit(1)
		}

		coloredMessage := fmt.Sprintf(errorhandler.RedFunc(err.Error()))
		fmt.Fprintln(os.Stderr, coloredMessage)
		os.Exit(1)
	} else {
		pivnetClient := commands.NewPivnetClient()
		client := commands.NewPivnetVersionsClient(pivnetClient)
		result := client.Warn(version.Version)
		if result != "" {
			fmt.Fprintf(os.Stderr, "\n%s\n", result)
		}

		rcFileReadWriter := filesystem.NewPivnetRCReadWriter(commands.Pivnet.ConfigFile)
		rch := rc.NewRCHandler(rcFileReadWriter)
		if commands.Pivnet.Profile != nil {
		    profile, err := rch.ProfileForName(commands.Pivnet.Profile.Name)
		    if err == nil && profile != nil {
				fmt.Fprint(os.Stderr, hostwarning.NewHostWarning(profile.Host).Warn())
			}
		}
	}
}
