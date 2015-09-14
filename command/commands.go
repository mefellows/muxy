package command

import (
	pki "github.com/mefellows/pkigo/command"
	"github.com/mitchellh/cli"
	"os"
)

var Commands map[string]cli.CommandFactory
var Ui cli.Ui

func init() {

	Ui = &cli.ColoredUi{
		Ui:          &cli.BasicUi{Writer: os.Stdout, Reader: os.Stdin, ErrorWriter: os.Stderr},
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorNone,
		ErrorColor:  cli.UiColorRed,
	}

	meta := Meta{
		Ui: Ui,
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
