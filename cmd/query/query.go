package query

import (
	"flag"
	"fmt"
	"github.com/gimke/riff/common"
	"github.com/gimke/riff/riff"
	"net"
	"net/rpc"
	"strconv"
	"time"
)

const help = `Usage: riff query <command> [options]

  Query riff service

Available subcommands are:

  nodes       Get nodes list
  snap        Get snap short.

Options:

  -addr   RPC address of riff (-addr 127.0.0.1:8530)
`

type cmd struct {
	flags *flag.FlagSet
	// flags
	addr  string
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}
func (c *cmd) init() {
	c.flags = flag.NewFlagSet("query", flag.ContinueOnError)
	c.flags.StringVar(&c.addr, "addr", "127.0.0.1:8530", "usage")

	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}
func (c *cmd) Run(args []string) int {
	if len(args) > 1 {
		if err := c.flags.Parse(args[1:]); err != nil {
			return 1
		}
	}
	//get args 0
	command := args[0]
	switch command {
	case "snap":
		c.SnapShort()
		return 0
		break
	case "nodes":
		c.Nodes()
		return 0
		break
	}
	//if c.snap {
	//	//call client
	//	c.SnapShort()
	//	return 0
	//}
	//if c.nodes {
	//	c.Nodes()
	//	return 0
	//}
	return 0
}

func (c *cmd) SnapShort() {
	conn, err := net.DialTimeout("tcp", c.addr, time.Second*10)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	//encBuf := bufio.NewWriter(conn)
	codec := common.NewGobClientCodec(conn)
	//codec := jsonrpc.NewClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	var snapshort string
	err = cmd.Call("Query.SnapShort", struct{}{}, &snapshort)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println(snapshort)
}

func (c *cmd) Nodes() {
	conn, err := net.DialTimeout("tcp", c.addr, time.Second*10)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	//encBuf := bufio.NewWriter(conn)
	codec := common.NewGobClientCodec(conn)
	//codec := jsonrpc.NewClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	var nodes riff.Nodes
	err = cmd.Call("Query.Nodes", struct{}{}, &nodes)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	fmt.Printf("%-16s %-10s %-24s %-8v %-48s\n", "Node", "DC", "Address", "Status", "SnapShort")
	for _, n := range nodes {
		fmt.Printf("%-16s %-10s %-24s %-8v %-48s\n",
			n.Name,
			n.DataCenter,
			net.JoinHostPort(n.IP.String(), strconv.Itoa(n.Port)),
			n.State.String(),
			n.SnapShort[0:10]+"...")
	}

}

func (c *cmd) Synopsis() string {
	return "Query nodes"
}

func (c *cmd) Help() string {
	return help
}
