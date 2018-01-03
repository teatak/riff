package start

import (
	"flag"
	"fmt"
	"github.com/gimke/riff/common"
	"github.com/gimke/riff/riff"
	"net"
	"os"
	"strings"
)

const synopsis = "Start Riff service"
const help = `Usage: start [options]

  Start riff service

Options:

  -http       Http address of riff (-http 127.0.0.1)
  -dns        Dns address of riff (-dns 127.0.0.1)
  -rpc        RPC address of riff (-rpc 0.0.0.0)
  -name       Node name.
  -dc         DataCenter name.
`

type cmd struct {
	flags *flag.FlagSet
	help  string
	// flags
	name string
	dc   string
	http string
	dns  string
	rpc  string
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}

func (c *cmd) init() {
	hostName, _ := os.Hostname()
	c.flags = flag.NewFlagSet("start", flag.ContinueOnError)
	c.flags.StringVar(&c.http, "http", "127.0.0.1", "usage")
	c.flags.StringVar(&c.dns, "dns", "127.0.0.1", "usage")
	c.flags.StringVar(&c.rpc, "rpc", "0.0.0.0", "usage")
	c.flags.StringVar(&c.name, "name", hostName, "usage")
	c.flags.StringVar(&c.dc, "dc", "dc1", "usage")

	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}
func (c *cmd) Run(args []string) int {
	var exit = make(chan bool)
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("riff.start error: %v\n", err)
		return 1
	}

	var adviseRpc string
	if common.IsAny(c.rpc) {
		var addrs []*net.IPAddr
		//detect ip
		var addrtyp string

		switch {
		case common.IsAnyV4(c.rpc):
			addrtyp = "private IPv4"
			addrs, err = common.GetPrivateIPv4()
			if err != nil {
				fmt.Println("Error detecting %s address: %s", addrtyp, err)
			}
			break
		case common.IsAnyV6(c.rpc):
			addrtyp = "public IPv6"
			addrs, err = common.GetPublicIPv6()
			if err != nil {
				fmt.Println("Error detecting %s address: %s", addrtyp, err)
			}
			break
		}
		adviseRpc = addrs[0].String()
	}
	if config.Addresses.Rpc == "" {
		config.IP = adviseRpc
		config.Addresses.Rpc = adviseRpc
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

//func (c *cmd) Ping() {
//	conn, err := net.DialTimeout("tcp", ":8530", time.Second*10)
//	if err != nil {
//		fmt.Println("error", err)
//		return
//	}
//	//encBuf := bufio.NewWriter(conn)
//	var exit = make(chan bool)
//	codec := common.NewGobClientCodec(conn)
//	//codec := jsonrpc.NewClientCodec(conn)
//	cmd := rpc.NewClientWithCodec(codec)
//	var reply string
//	for {
//		err = cmd.Call("Status.Ping", struct{}{}, &reply)
//		if err != nil {
//			fmt.Println("error", err)
//			close(exit)
//			break
//		}
//		fmt.Println(reply)
//		time.Sleep(5 * time.Second)
//	}
//	<-exit
//	//cmd.Close()
//	//if err != nil && errc != nil {
//	//	return fmt.Errorf("%s %s", err, errc)
//	//}
//	//if err != nil {
//	//	return err
//	//} else {
//	//	return errc
//	//}
//}
