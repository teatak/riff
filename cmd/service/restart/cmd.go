package restart

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

const help = `Usage: riff restart <name> [options]

  restart service

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
	c.flags = flag.NewFlagSet("restart", flag.ContinueOnError)
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
	c.ReStart(name)
	return 0
}

func (c *cmd) ReStart(name string) {
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
		Cmd:  api.CmdRestart,
	}, &result)
	if err != nil {
		fmt.Println(err)
	}
	if result {
		fmt.Printf("service %s restart success\n", name)
	}
}

func (c *cmd) Synopsis() string {
	return "Restart service"
}

func (c *cmd) Help() string {
	return help
}
