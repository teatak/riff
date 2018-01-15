package start

import (
	"github.com/gimke/riff/common"
	"github.com/gimke/riff/riff"
	"net"
)

func loadConfig(cmd *cmd) (*riff.Config, error) {
	var host string
	var port int
	var err error

	c := riff.LoadConfig()

	advise, err := common.AdviseRpc()
	if err != nil {
		return nil, err
	}
	c.IP = advise

	//http
	host, port = common.GetIpPort(cmd.http)
	if host != "" {
		c.Addresses.Http = host
	}
	if port != 0 {
		c.Ports.Http = port
	}

	host, port = common.GetIpPort(cmd.dns)
	if host != "" {
		c.Addresses.Dns = host
	}
	if port != 0 {
		c.Ports.Dns = port
	}

	host, port = common.GetIpPort(cmd.rpc)
	if host != "" {
		c.Addresses.Rpc = host
	}
	if port != 0 {
		c.Ports.Rpc = port
	}

	if c.Addresses.Rpc == "" {
		ip, _, _ := net.ParseCIDR(advise)
		c.Addresses.Rpc = ip.String()
	}
	if cmd.join != "" {
		c.Join = cmd.join
	}
	if cmd.name != "" {
		c.Name = cmd.name
	}
	if cmd.dc != "" {
		c.DataCenter = cmd.dc
	}
	return c, nil
}
