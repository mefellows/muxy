package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mefellows/muxy/command"
	_ "github.com/mefellows/muxy/middleware"
	_ "github.com/mefellows/muxy/protocol"
	_ "github.com/mefellows/muxy/symptom"
	"github.com/mitchellh/cli"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	cli := cli.NewCLI(strings.ToLower(ApplicationName), Version)
	cli.Args = os.Args[1:]
	cli.Commands = command.Commands

	exitStatus, err := cli.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	return exitStatus
}
