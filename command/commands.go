// Package command contains the CLI options for Muxy.
package command

import (
	"os"

	m "github.com/mefellows/muxy/muxy"
	pki "github.com/mefellows/pkigo/command"
	"github.com/mitchellh/cli"
)

// Commands contains all CLI commands available
var Commands map[string]cli.CommandFactory

// UI wraps the commands available to the CLI
var UI cli.Ui

var muxy *m.Muxy
var c *m.Config

func init() {
	setup()
}

func setup() {
	c = &m.Config{}
	muxy = m.New(c)

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
