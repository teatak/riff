// +build !windows

package run

import (
	"fmt"
	"github.com/gimke/riff/riff"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"io/ioutil"
	"github.com/gimke/riff/common"
)

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	config, err := loadConfig(c)
	if err != nil {
		fmt.Printf("riff.start error: %v\n", err)
		return 1
	}
	s, err := riff.NewServer(config)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer s.Shutdown()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	go func() {
		for {
			sig := <-sigs
			fmt.Println()
			s.Logger.Printf(infoServerPrefix+"get signal %v\n", sig)
			if sig == syscall.SIGUSR2 {
				s.Shutdown()
			} else {
				s.Shutdown()
			}
		}
	}()
	<-s.ShutdownCh
	return 0
}
