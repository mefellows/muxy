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

	cmdFlags.IntVar(&c.Port, "port", 8123, "The http port to listen on")
	cmdFlags.StringVar(&c.Host, "host", "0.0.0.0", "The host/ip to bind to")
	cmdFlags.IntVar(&c.ProxyPort, "proxyPort", -1, "The proxied hosts http port")
	cmdFlags.StringVar(&c.ProxyHost, "proxyHost", "", "The proxied hosts ip/hostname")
	cmdFlags.StringVar(&c.ProxyProtocol, "proxyProtocol", "http", "The proxied hosts protocol (http)")
	cmdFlags.StringVar(&c.ConfigFile, "config", "", "Path to a YAML configuration file")
	// TODO: SSL
	//cmdFlags.BoolVar(&c.Insecure, "insecure", false, "Disable TLS connection between Muxy <-> Proxied Host")
	//cmdFlags.BoolVar(&c.Ssl, "ssl", false, "Disable TLS connection between Client <-> Muxy")

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

  Run the Muck proxy to 
  
Options:

  --port                      The http(s) port to listen on
  --host                      The IP address to listen on. Defaults to 0.0.0.0
  --proxyHost                 The http(s) host to proxy
  --proxyPort                 The http(s) port to proxy
  --proxyProtocol             The protocol to proxy (e.g. http)
  --config                    Location of Muxy configuration file
`

	return strings.TrimSpace(helpText)
}

func (c *ProxyCommand) Synopsis() string {
	return "Run the Muxy proxy"
}
