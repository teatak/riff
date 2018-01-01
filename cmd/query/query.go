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

type cmd struct {
	flags *flag.FlagSet
	// flags
	snap  bool
	nodes bool
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}
func (c *cmd) init() {
	c.flags = flag.NewFlagSet("query", flag.ContinueOnError)
	c.flags.BoolVar(&c.snap, "snap", false, "usage")
	c.flags.BoolVar(&c.nodes, "nodes", false, "usage")

	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}
func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	if c.snap {
		//call client
		c.SnapShort()
		return 0
	}
	if c.nodes {
		c.Nodes()
		return 0
	}
	return 0
}

func (c *cmd) SnapShort() {
	conn, err := net.DialTimeout("tcp", "192.168.1.220:8530", time.Second*10)
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
	conn, err := net.DialTimeout("tcp", "192.168.1.220:8530", time.Second*10)
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
	fmt.Printf("%-16s %-16s %-24s %-8v %-48s\n", "Node", "DC", "Address", "Status", "SnapShort")
	for _, n := range nodes {
		fmt.Printf("%-16s %-16s %-24s %-8v %-48s\n",
			n.Name,
			n.DataCenter,
			net.JoinHostPort(n.IP.String(), strconv.Itoa(n.Port)),
			riff.GetState(n.State),
			n.SnapShort)
	}
}

func (c *cmd) Synopsis() string {
	return "Query nodes"
}

func (c *cmd) Help() string {
	return ""
}
