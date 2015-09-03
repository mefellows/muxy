package main

import (
	"fmt"
	"github.com/mefellows/muxy/command"
	_ "github.com/mefellows/muxy/middleware"
	_ "github.com/mefellows/muxy/symptom"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	cli := cli.NewCLI(strings.ToLower(APPLICATION_NAME), VERSION)
	cli.Args = os.Args[1:]
	cli.Commands = command.Commands

	exitStatus, err := cli.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	return exitStatus
}
