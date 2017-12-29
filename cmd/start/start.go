package start

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"strings"
	"time"
	"github.com/gimke/riff/common"
	"github.com/gimke/riff/riff"
)

const synopsis = "Start a service"
const help = `Usage: serf info [options]

  Provides debugging information for operators

Options:

  -format     If provided, output is returned in the specified
              format. Valid formats are 'json', and 'text' (default)
  -rpc-addr   RPC address of the Serf agent.
  -rpc-auth   RPC auth token of the Serf agent.
`

type cmd struct {
	flags *flag.FlagSet
	help  string
	// flags
	ping bool
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("start", flag.ContinueOnError)
	c.flags.BoolVar(&c.ping, "ping", false,
		"usage")
	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}
func (c *cmd) Run(args []string) int {
	var exit = make(chan bool)
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	if c.ping {
		//call client
		c.Ping()
		return 0
	}
	s, err := riff.NewServer()
	if err != nil {
		fmt.Println(err)
	}
	defer s.Shutdown()
	<-exit
	return 0
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return strings.TrimSpace(help)

}

func (c *cmd) Ping() {
	conn, err := net.DialTimeout("tcp", ":8530", time.Second*10)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	//encBuf := bufio.NewWriter(conn)
	var exit = make(chan bool)
	codec := common.NewGobClientCodec(conn)
	//codec := jsonrpc.NewClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	var reply string
	for {
		err = cmd.Call("Status.Ping", struct{}{}, &reply)
		if err != nil {
			fmt.Println("error", err)
			close(exit)
			break
		}
		fmt.Println(reply)
		time.Sleep(5*time.Second)
	}
	<-exit
	//cmd.Close()
	//if err != nil && errc != nil {
	//	return fmt.Errorf("%s %s", err, errc)
	//}
	//if err != nil {
	//	return err
	//} else {
	//	return errc
	//}
}
