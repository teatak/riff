package service

import (
	"flag"
	"fmt"
	"github.com/gimke/riff/api"
	"github.com/gimke/riff/common"
	"net"
	"net/rpc"
	"strconv"
	"strings"
	"time"
)

var help = `Usage: riff %s <name> [options]

  %s service

Options:

  -rpc    RPC address of riff (-rpc 192.168.1.1:8630)
`

type cmd struct {
	flags *flag.FlagSet
	// flags
	rpc     string
	cmdType api.CmdType
}

func New(cmdType api.CmdType) *cmd {
	c := &cmd{
		cmdType: cmdType,
	}
	c.init()
	return c
}
func (c *cmd) init() {
	c.flags = flag.NewFlagSet("start", flag.ContinueOnError)
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
	c.Cmd(name)
	return 0
}

func (c *cmd) Cmd(name string) {
	conn, err := net.DialTimeout("tcp", c.rpc, time.Second*10)
	if err != nil {
		fmt.Println(err)
		return
	}
	codec := api.NewGobClientCodec(conn)
	client := rpc.NewClientWithCodec(codec)
	defer client.Close()

	var result bool
	err = client.Call("Mutation.Service", api.ParamServiceMutation{
		Name: name,
		Cmd:  c.cmdType,
	}, &result)
	if err != nil {
		fmt.Println(err)
	}
	if result {
		fmt.Printf("service %s %s success\n", name, strings.ToLower(c.cmdType.String()))
	}
}

func (c *cmd) Synopsis() string {
	return c.cmdType.String() + " service"
}

func (c *cmd) Help() string {
	cmdName := strings.ToLower(c.cmdType.String())
	return fmt.Sprintf(help, cmdName, cmdName)
}
