package command

import (
	"flag"
	"strings"

	"github.com/mefellows/muxy/muxy"
)

// ProxyCommand enables an http proxy for http tampering
type ProxyCommand struct {
	Meta Meta
}

// Run the HTTP Proxy CLI command
func (pc *ProxyCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("proxy", flag.ContinueOnError)
	cmdFlags.Usage = func() { pc.Meta.UI.Output(pc.Help()) }
	c := &muxy.Config{}

	cmdFlags.StringVar(&c.ConfigFile, "config", "", "Path to a YAML configuration file")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	muxy := muxy.New(c)
	muxy.Run()

	return 0
}

// Help prints out detailed help for this command
func (pc *ProxyCommand) Help() string {
	helpText := `
Usage: muck proxy [options]

  Run the Muck proxy.

Options:

  --config                    Location of Muxy configuration file
`

	return strings.TrimSpace(helpText)
}

// Synopsis prints out help for this command
func (pc *ProxyCommand) Synopsis() string {
	return "Run the Muxy proxy"
}
