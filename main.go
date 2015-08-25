package main

import (
	"fmt"
	"github.com/mefellows/muck/command"
	"github.com/mitchellh/cli"
	"os"
)

func main() {
	cli := cli.NewCLI("muck", "1.0.0")
	cli.Args = os.Args[1:]
	cli.Commands = command.Commands

	exitStatus, err := cli.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)
}
