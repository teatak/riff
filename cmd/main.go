package main

import (
	"fmt"
	"os"

	"github.com/teatak/riff/api"
	"github.com/teatak/riff/cli"
	"github.com/teatak/riff/cmd/cli/daem"
	"github.com/teatak/riff/cmd/cli/query"
	"github.com/teatak/riff/cmd/cli/quit"
	"github.com/teatak/riff/cmd/cli/reload"
	"github.com/teatak/riff/cmd/cli/run"
	"github.com/teatak/riff/cmd/cli/service"
	"github.com/teatak/riff/cmd/cli/update"
	"github.com/teatak/riff/cmd/cli/version"
	"github.com/teatak/riff/common"
)

var Commands cli.Commands

func init() {
	Commands = cli.Commands{
		"version": version.New(common.Version),
		"daem":    daem.New(),
		"quit":    quit.New(),
		"update":  update.New(),
		"reload":  reload.New(),
		"run":     run.New(),
		"query":   query.New(),
		"start":   service.New(api.CmdStart),
		"stop":    service.New(api.CmdStop),
		"restart": service.New(api.CmdRestart),
	}
}
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

const cheers = `
dP""Y88b   d8P''Y8b  .dP""Y88b  o888  o888  .dP""Y88b   d8P''Y8b  .dP""Y88b  
	 ]8P' 888    888       ]8P'  888   888        ]8P' 888    888       ]8P' 
   .d8P'  888    888     .d8P'   888   888      .d8P'  888    888     .d8P'  
 .dP'     888    888   .dP'      888   888    .dP'     888    888   .dP'     
.oP     .o '88b  d88' .oP     .o  888   888  .oP     .o '88b  d88' .oP     .o 
8888888888  'Y8bd8P'  8888888888 o888o o888o 8888888888  'Y8bd8P'  8888888888 
`
