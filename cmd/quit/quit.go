package quit

import (
	"flag"
	"fmt"
	"github.com/gimke/riff/common"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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

func (c *cmd) GetPid() int {
	content, err := ioutil.ReadFile(common.BinDir + "/run/riff.pid")
	if err != nil {
		return 0
	} else {
		pid, _ := strconv.Atoi(strings.Trim(string(content), "\n"))
		if _, find := c.processExist(pid); find {
			return pid
		} else {
			return 0
		}
	}
}

func (c *cmd) processExist(pid int) (*os.Process, bool) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return nil, false
	} else {
		err := process.Signal(syscall.Signal(0))
		if err != nil {
			return nil, false
		}
	}
	return process, true
}

func (c *cmd) Quit() {
	pid := c.GetPid()
	if pid == 0 {
		fmt.Println("can't find riff")
	} else {
		if p, find := c.processExist(pid); find {
			err := p.Signal(syscall.SIGINT)
			if err != nil {
				fmt.Println(err)
			} else {
				quitStop := make(chan bool)
				go func() {
					for {
						if pid := c.GetPid(); pid == 0 {
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
