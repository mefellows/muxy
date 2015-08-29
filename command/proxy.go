package command

import (
	"flag"
	"github.com/mefellows/muxy/muxy"
	"strings"
)

type ProxyCommand struct {
	Meta     Meta
	Port     int    // Which port to listen on
	Host     string // Which network host/ip to listen on
	Insecure bool   // Enable/Disable TLS between Muck <-> Proxied Host
	Ssl      bool   // Enable/Disable TLS between client <-> Muck
	Target   string // Proxy target
}

func (c *ProxyCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("sync", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Meta.Ui.Output(c.Help()) }

	cmdFlags.IntVar(&c.Port, "port", 8123, "The http port to listen on")
	cmdFlags.StringVar(&c.Host, "host", "", "The host/ip to bind to. Defaults to 0.0.0.0")
	cmdFlags.BoolVar(&c.Insecure, "insecure", false, "Disable TLS connection between Muck <-> Proxied Host")
	cmdFlags.BoolVar(&c.Ssl, "ssl", false, "Disable TLS connection between Client <-> Muck")
	cmdFlags.StringVar(&c.Target, "target", "", "The proxy target")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	muxy := &muxy.Muxy{}
	muxy.Run()

	return 0
}

func (c *ProxyCommand) Help() string {
	helpText := `
Usage: muck proxy [options] 

  Run the Muck proxy to 
  
Options:

  --target                    The http(s) endpoint to proxy
  --port                      The http(s) port to listen on
  --host                      The IP address to listen on. Defaults to 0.0.0.0
  --ssl                       Enable SSL security on the proxy side of the connection
  --insecure                  Disable security between Muck and proxied server
`

	return strings.TrimSpace(helpText)
}

func (c *ProxyCommand) Synopsis() string {
	return "Run the Muck proxy service"
}
