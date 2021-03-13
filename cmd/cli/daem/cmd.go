package daem

import (
	"flag"
	"fmt"
	"github.com/teatak/riff/cmd/cli"
	"github.com/teatak/riff/common"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const synopsis = "Run Riff as service "
const help = `Usage: daem [options]

  Run riff as daemon service

Options:

  -name       Node name
  -dc         DataCenter name
  -http       Http address of riff (-http 127.0.0.1:8610)
  -rpc        RPC address of riff (-rpc [::]:8630)
  -join       Join RPC address (-join 192.168.1.1:8630,192.168.1.2:8630,192.168.1.3:8630)
`

const infoServerPrefix = "[INFO] riff.server: "

type cmd struct {
	flags *flag.FlagSet
	help  string
	// flags
	name string
	dc   string
	http string
	dns  string
	rpc  string
	join string
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("start", flag.ContinueOnError)
	c.flags.StringVar(&c.http, "http", "", "usage")
	c.flags.StringVar(&c.dns, "dns", "", "usage")
	c.flags.StringVar(&c.rpc, "rpc", "", "usage")
	c.flags.StringVar(&c.name, "name", "", "usage")
	c.flags.StringVar(&c.join, "join", "", "usage")
	c.flags.StringVar(&c.dc, "dc", "", "usage")

	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}

func (c *cmd) Run(args []string) int {
	if err := c.Start(args); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

func (c *cmd) Start(args []string) error {
	if cli.GetPid() != 0 {
		return fmt.Errorf("%s is already running", "riff")
	}
	command := c.resoveCommand(common.BinDir + "/riff")
	dir, _ := filepath.Abs(filepath.Dir(command))

	newArgs := append([]string{}, "run")
	newArgs = append(newArgs, args...)
	args = newArgs
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	out := common.MakeFile(common.BinDir + "/logs/riff/stdout.log")
	cmd.Stdout = out
	cmd.Stderr = out

	err := cmd.Start()
	if err != nil {
		return err
	} else {
		cli.SetPid(cmd.Process.Pid)
		fmt.Println("start riff success")
	}
	return nil
}

func (c *cmd) resoveCommand(path string) string {
	if filepath.IsAbs(path) {
		return path
	} else {
		if strings.HasPrefix(path, "."+string(os.PathSeparator)) {
			return common.BinDir + path[1:]
		} else {
			return path
		}
	}
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return strings.TrimSpace(help)
}
