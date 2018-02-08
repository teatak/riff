package query

import (
	"flag"
	"fmt"
	"github.com/gimke/riff/api"
	"github.com/gimke/riff/common"
	"github.com/ryanuber/columnize"
	"net"
	"strconv"
)

const help = `Usage: riff query <command> [options]

  Query riff service

Available subcommands are:

  snap        Get snap short.
  nodes       Get nodes.
  node        Get node [name].
  services    Get services.
  service     Get service [name].

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
	case "services":
		c.Services()
		return 0
		break
	case "service":
		name := args[1]
		c.Service(name)
		return 0
		break
	}

	return 0
}

func (c *cmd) SnapShot() {
	client,err := api.NewClient(c.rpc)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	var snapshot string
	err = client.Call("Query.SnapShot", struct{}{}, &snapshot)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println(snapshot)
}

func (c *cmd) Nodes() {
	client,err := api.NewClient(c.rpc)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	var nodes api.Nodes
	err = client.Call("Query.Nodes", struct{}{}, &nodes)
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
	client,err := api.NewClient(c.rpc)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	var node api.Node
	err = client.Call("Query.Node", api.Node{Name: name}, &node)
	if err != nil {
		fmt.Println(err)
		return
	}
	info := make([]string, 0, 5)

	info = append(info, "Node:|"+node.Name)
	info = append(info, "Address:|"+net.JoinHostPort(node.IP, strconv.Itoa(node.Port)))
	info = append(info, "Status:|"+node.State.String())
	info = append(info, "DC:|"+node.DataCenter)
	info = append(info, "SnapShot:|"+node.SnapShot)

	output := columnize.SimpleFormat(info)

	fmt.Println(output)

	if len(node.NestServices) > 0 {
		fmt.Println("")
		//output service
		results := make([]string, 0, len(node.NestServices))
		header := "Service|Port|Status"
		results = append(results, header)

		for _, s := range node.NestServices {
			line := fmt.Sprintf("%s|%s|%s",
				s.Name,
				net.JoinHostPort("", strconv.Itoa(s.Port)),
				s.State.String())
			results = append(results, line)
		}

		output = columnize.SimpleFormat(results)
		fmt.Println(output)
	}
}

func (c *cmd) Services() {
	client,err := api.NewClient(c.rpc)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	var services api.Services
	err = client.Call("Query.Services", struct{}{}, &services)
	if err != nil {
		fmt.Println(err)
		return
	}
	results := make([]string, 0, len(services)+1)

	for _, s := range services {
		line := fmt.Sprintf("%s", s.Name)
		results = append(results, line)
	}

	output := columnize.SimpleFormat(results)
	fmt.Println(output)
}

func (c *cmd) Service(name string) {
	client,err := api.NewClient(c.rpc)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	var service api.Service
	err = client.Call("Query.Service", api.ParamService{Name: name, State: api.StateAll}, &service)
	if err != nil {
		fmt.Println(err)
		return
	}
	results := make([]string, 0, len(service.NestNodes)+1)
	header := "Node|Address|Status|DC|SnapShot"
	results = append(results, header)

	for _, n := range service.NestNodes {
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

func (c *cmd) Synopsis() string {
	return "Query"
}

func (c *cmd) Help() string {
	return help
}
