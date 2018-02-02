package query

import (
	"flag"
	"fmt"
	"github.com/gimke/riff/api"
	"github.com/gimke/riff/common"
	"github.com/ryanuber/columnize"
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

  -rpc    RPC address of riff (-rpc 192.168.1.1:8630)
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
	c.flags.StringVar(&c.rpc, "rpc", "", "usage")

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
	advise, _ := common.AdviseRpc()
	host, port := common.GetIpPort(c.rpc)
	if host == "" {
		ip, _, _ := net.ParseCIDR(advise)
		host = ip.String()
	}
	if port == 0 {
		port = common.DefaultRpcPort
	}
	c.rpc = net.JoinHostPort(host, strconv.Itoa(port))

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
	case "node":
		name := args[1]
		c.Node(name)
		return 0
		break
	}

	return 0
}

func (c *cmd) SnapShot() {
	conn, err := net.DialTimeout("tcp", c.rpc, time.Second*10)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	codec := api.NewGobClientCodec(conn)
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
	codec := api.NewGobClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	var nodes api.Nodes
	err = cmd.Call("Query.Nodes", struct{}{}, &nodes)
	if err != nil {
		fmt.Println(err)
		return
	}
	results := make([]string, 0, len(nodes)+1)
	header := "Node|Address|Status|DC|SnapShot"
	results = append(results, header)

	for _, n := range nodes {
		line := fmt.Sprintf("%s|%s|%s|%s|%s",
			n.Name,
			net.JoinHostPort(n.IP, strconv.Itoa(n.Port)),
			n.State.String(),
			n.DataCenter,
			n.SnapShot[0:9]+"...")
		results = append(results, line)
	}

	output := columnize.SimpleFormat(results)
	fmt.Println(output)
}

func (c *cmd) Node(name string) {
	conn, err := net.DialTimeout("tcp", c.rpc, time.Second*10)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	codec := api.NewGobClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	var node api.Node
	err = cmd.Call("Query.Node", api.Node{Name:name}, &node)
	if err != nil {
		fmt.Println(err)
		return
	}
	results := make([]string, 0, 1)
	header := "Node|Address|Status|DC|SnapShot"
	results = append(results, header)

	line := fmt.Sprintf("%s|%s|%s|%s|%s",
		node.Name,
		net.JoinHostPort(node.IP, strconv.Itoa(node.Port)),
		node.State.String(),
		node.DataCenter,
		node.SnapShot[0:9]+"...")
	results = append(results, line)

	output := columnize.SimpleFormat(results)
	fmt.Println(output)
}

func (c *cmd) Synopsis() string {
	return "Query nodes"
}

func (c *cmd) Help() string {
	return help
}
