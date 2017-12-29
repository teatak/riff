package main

import (
	"github.com/gimke/riff/cli"
	"github.com/gimke/riff/cmd/start"
	"github.com/gimke/riff/cmd/version"
	"github.com/gimke/riff/common"
)

var Commands cli.Commands

func init() {
	Commands = cli.Commands{
		"version": version.New(common.Version),
		"start":   start.New(),
	}
}
