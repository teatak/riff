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

  -rpc    RPC address of riff (-rpc 0.0.0.0:8630)
`

type cmd struct {
	flags *flag.FlagSet
	// flags
	rpc string
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}
func (c *cmd) init() {
	c.flags = flag.NewFlagSet("query", flag.ContinueOnError)
	c.flags.StringVar(&c.rpc, "rpc", "0.0.0.0:8630", "usage")

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
		c.SnapShot()
		return 0
		break
	case "nodes":
		c.Nodes()
		return 0
		break
	}
	//if c.snap {
	//	//call client
	//	c.SnapShot()
	//	return 0
	//}
	//if c.nodes {
	//	c.Nodes()
	//	return 0
	//}
	return 0
}

func (c *cmd) SnapShot() {
	conn, err := net.DialTimeout("tcp", c.rpc, time.Second*10)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	//encBuf := bufio.NewWriter(conn)
	codec := common.NewGobClientCodec(conn)
	//codec := jsonrpc.NewClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	var snapshot string
	err = cmd.Call("Query.SnapShot", struct{}{}, &snapshot)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println(snapshot)
}

func (c *cmd) Nodes() {
	conn, err := net.DialTimeout("tcp", c.rpc, time.Second*10)
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
	fmt.Printf("%-12s %-24s %-10s %-24s %-8v %-12s\n", "Id", "Node", "DC", "Address", "Status", "SnapShot")
	for _, n := range nodes {
		fmt.Printf("%-12s %-24s %-10s %-24s %-8v %-12s\n",
			n.Id,
			n.Name,
			n.DataCenter,
			net.JoinHostPort(n.IP, strconv.Itoa(n.Port)),
			n.State.String(),
			n.SnapShot[0:9]+"...")
	}

}

func (c *cmd) Synopsis() string {
	return "Query nodes"
}

func (c *cmd) Help() string {
	return help
}
