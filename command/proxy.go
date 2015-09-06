package command

import (
	"flag"
	"github.com/mefellows/muxy/muxy"
	"strings"
)

type ProxyCommand struct {
	Meta Meta
}

func (pc *ProxyCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("proxy", flag.ContinueOnError)
	cmdFlags.Usage = func() { pc.Meta.Ui.Output(pc.Help()) }
	c := &muxy.MuxyConfig{}

	cmdFlags.StringVar(&c.ConfigFile, "config", "", "Path to a YAML configuration file")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	muxy := muxy.New(c)
	muxy.Run()

	return 0
}

func (c *ProxyCommand) Help() string {
	helpText := `
Usage: muck proxy [options] 

  Run the Muck proxy.
  
Options:

  --config                    Location of Muxy configuration file
`

	return strings.TrimSpace(helpText)
}

func (c *ProxyCommand) Synopsis() string {
	return "Run the Muxy proxy"
}
