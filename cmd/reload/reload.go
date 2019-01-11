package reload

import (
	"fmt"
	"flag"
	"net"
	"strconv"
	"github.com/gimke/riff/common"
	"github.com/gimke/riff/api"
)

const help = `Usage: riff reload

Options:

  -rpc    RPC address of riff (-rpc 192.168.1.1:8630)
`

type cmd struct {
	flags *flag.FlagSet
	rpc string
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("reload", flag.ContinueOnError)
	c.flags.StringVar(&c.rpc, "rpc", "", "usage")

	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
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
	c.Reload()
	return 0
}

func (c *cmd) Reload() {
	client, err := api.NewClient(c.rpc)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	var result bool
	err = client.Call("Mutation.Reload", struct{}{}, &result)
	if err != nil {
		fmt.Println("error", err)
	}
	if result {
		fmt.Println("reload success")
	} else {
		fmt.Println("reload failed")
	}
}

func (c *cmd) Synopsis() string {
	return "Reload riff config"
}

func (c *cmd) Help() string {
	return help
}