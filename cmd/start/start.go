package start

import (
	"flag"
	"fmt"
	"github.com/gimke/riff/common"
	"github.com/gimke/riff/riff"
	"net"
	"net/rpc"
	"os"
	"strings"
	"time"
)

const synopsis = "Start riff"
const help = `Usage: start [options]

  Start riff service

Options:

  -bind       RPC address of riff (-bind [::]:8530)
  -name       Node name.
`

type cmd struct {
	flags *flag.FlagSet
	help  string
	// flags
	bind string
	name string
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}

func (c *cmd) init() {
	hostName, _ := os.Hostname()
	c.flags = flag.NewFlagSet("start", flag.ContinueOnError)
	c.flags.StringVar(&c.bind, "bind", ":8530", "usage")
	c.flags.StringVar(&c.name, "name", hostName, "usage")

	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}
func (c *cmd) Run(args []string) int {
	var exit = make(chan bool)
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	config, err := riff.NewConfig(c.bind, c.name)
	if err != nil {
		fmt.Printf("riff.start error:%v\n", err)
		return 1
	}
	s, err := riff.NewServer(config)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer s.Shutdown()
	<-exit
	return 0
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return strings.TrimSpace(help)

}

func (c *cmd) Ping() {
	conn, err := net.DialTimeout("tcp", ":8530", time.Second*10)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	//encBuf := bufio.NewWriter(conn)
	var exit = make(chan bool)
	codec := common.NewGobClientCodec(conn)
	//codec := jsonrpc.NewClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	var reply string
	for {
		err = cmd.Call("Status.Ping", struct{}{}, &reply)
		if err != nil {
			fmt.Println("error", err)
			close(exit)
			break
		}
		fmt.Println(reply)
		time.Sleep(5 * time.Second)
	}
	<-exit
	//cmd.Close()
	//if err != nil && errc != nil {
	//	return fmt.Errorf("%s %s", err, errc)
	//}
	//if err != nil {
	//	return err
	//} else {
	//	return errc
	//}
}
