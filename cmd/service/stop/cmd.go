package stop

import (
	"flag"
	"fmt"
	"github.com/gimke/riff/api"
	"github.com/gimke/riff/common"
	"net"
	"net/rpc"
	"strconv"
	"time"
)

const help = `Usage: riff stop <name> [options]

  stop service

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
	c.flags = flag.NewFlagSet("stop", flag.ContinueOnError)
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
	name := args[0]
	c.Stop(name)
	return 0
}

func (c *cmd) Stop(name string) {
	conn, err := net.DialTimeout("tcp", c.rpc, time.Second*10)
	if err != nil {
		fmt.Println(err)
		return
	}
	codec := api.NewGobClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	var result bool
	err = cmd.Call("Mutation.Service", api.ParamServiceMutation{
		Name: name,
		Cmd:  api.CmdStop,
	}, &result)
	if err != nil {
		fmt.Println(err)
	}
	if result {
		fmt.Printf("service %s stop success\n", name)
	}
}

func (c *cmd) Synopsis() string {
	return "Stop service"
}

func (c *cmd) Help() string {
	return help
}
