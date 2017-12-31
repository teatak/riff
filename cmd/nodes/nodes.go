package nodes

import (
	"fmt"
	"flag"
	"net"
	"time"
	"github.com/gimke/riff/common"
	"net/rpc"
)

type cmd struct {
	flags *flag.FlagSet
	// flags
	list bool
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}
func (c *cmd) init() {
	c.flags = flag.NewFlagSet("nodes", flag.ContinueOnError)
	c.flags.BoolVar(&c.list, "list", false,
		"usage")
	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}
func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	if c.list {
		//call client
		c.List()
		return 0
	}
	return 0
}

func (c *cmd) List() {
	conn, err := net.DialTimeout("tcp", ":8530", time.Second*10)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	//encBuf := bufio.NewWriter(conn)
	codec := common.NewGobClientCodec(conn)
	//codec := jsonrpc.NewClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	var reply string
	err = cmd.Call("Nodes.List", struct{}{}, &reply)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println(reply)
}

func (c *cmd) Synopsis() string {
	return "Get nodes"
}

func (c *cmd) Help() string {
	return ""
}