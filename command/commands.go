package command

import (
	"os"

	pki "github.com/mefellows/pkigo/command"
	"github.com/mitchellh/cli"
)

// Commands contains all CLI commands available
var Commands map[string]cli.CommandFactory

// UI wraps the commands available to the CLI
var UI cli.Ui

func init() {

	UI = &cli.ColoredUi{
		Ui:          &cli.BasicUi{Writer: os.Stdout, Reader: os.Stdin, ErrorWriter: os.Stderr},
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorNone,
		ErrorColor:  cli.UiColorRed,
	}

	meta := Meta{
		UI: UI,
	}

	Commands = map[string]cli.CommandFactory{
		"proxy": func() (cli.Command, error) {
			return &ProxyCommand{
				Meta: meta,
			}, nil
		},
		"pki": func() (cli.Command, error) {
			return &pki.PkiCommand{}, nil
		},
	}
}
