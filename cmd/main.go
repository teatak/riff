package main

import (
	"fmt"
	"github.com/gimke/riff/cli"
	"github.com/gimke/riff/common"
	"os"
)

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "cheers" {
			fmt.Println(cheers)
			return
		}
		if arg == "--" {
			break
		}

		if arg == "-v" || arg == "--version" {
			args = []string{"version"}
			break
		}
	}

	c := cli.NewCLI("riff", common.Version)
	c.Args = args
	c.Commands = Commands
	exitCode, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
	}

	os.Exit(exitCode)
}
