package command

import "testing"

func TestCommands_Proxy(t *testing.T) {
	setup()
	meta := Meta{
		UI: UI,
	}

	pc := ProxyCommand{Meta: meta}
	pc.Help()
	pc.Synopsis()
}
