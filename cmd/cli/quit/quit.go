package quit

import (
	"flag"
	"fmt"
	"github.com/teatak/riff/cmd/cli"
	"syscall"
	"time"
)

const help = `Usage: riff quit
`

type cmd struct {
	flags *flag.FlagSet
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("quit", flag.ContinueOnError)

	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	c.Quit()
	return 0
}

func (c *cmd) Quit() {
	pid := cli.GetPid()
	if pid == 0 {
		fmt.Println("can't find riff")
	} else {
		if p, find := cli.ProcessExist(pid); find {
			err := p.Signal(syscall.SIGINT)
			if err != nil {
				fmt.Println(err)
			} else {
				quitStop := make(chan bool)
				go func() {
					for {
						if pid := cli.GetPid(); pid == 0 {
							quitStop <- true
							break
						}
						time.Sleep(1 * time.Second)
					}
				}()
				<-quitStop
				fmt.Println("quit riff success")
			}
		}
	}
}

func (c *cmd) Synopsis() string {
	return "Quit Riff"
}

func (c *cmd) Help() string {
	return help
}
