package start

import (
	"github.com/gimke/riff/common"
	"github.com/gimke/riff/riff"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
)

func isExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func defaultConfig() *riff.Config {
	id := common.GenerateID(12)
	hostName, _ := os.Hostname()
	c := &riff.Config{
		Id:         id,
		Name:       hostName,
		DataCenter: "dc1",
		AutoPilot:  false,
		Fanout:     3,
		Addresses: &riff.Addresses{
			Http: "127.0.0.1",
			Dns:  "127.0.0.1",
			Rpc:  "",
		},
		Ports: &riff.Ports{
			Http: 8610,
			Dns:  8620,
			Rpc:  8630,
		},
	}
	return c
}
func mergeConfig(src, dest *riff.Config) {
	if src.Id != "" {
		dest.Id = src.Id
	}
	if src.Name != "" {
		dest.Name = src.Name
	}
	if src.DataCenter != "" {
		dest.DataCenter = src.DataCenter
	}
	if src.IP != "" {
		dest.IP = src.IP
	}
	if src.Join != "" {
		dest.Join = src.Join
	}
	if src.Addresses != nil {
		if src.Addresses.Http != "" {
			dest.Addresses.Http = src.Addresses.Http
		}
		if src.Addresses.Dns != "" {
			dest.Addresses.Dns = src.Addresses.Dns
		}
		if src.Addresses.Rpc != "" {
			dest.Addresses.Rpc = src.Addresses.Rpc
		}
	}
	if src.Ports != nil {
		if src.Ports.Http != 0 {
			dest.Ports.Http = src.Ports.Http
		}
		if src.Ports.Dns != 0 {
			dest.Ports.Dns = src.Ports.Dns
		}
		if src.Ports.Rpc != 0 {
			dest.Ports.Rpc = src.Ports.Rpc
		}
	}
}

func loadConfig(cmd *cmd) (*riff.Config, error) {
	var host string
	var port int
	var err error

	c := defaultConfig()
	file := common.BinDir + "/config/" + common.Name + ".yml"
	if isExist(file) {
		//return nil, fmt.Errorf("file not exist %s", file)
		content, _ := ioutil.ReadFile(file)
		var file = &riff.Config{}
		err := yaml.Unmarshal(content, &c)
		if err != nil {
			return nil, err
		}
		mergeConfig(file, c)
	}

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
	if !isExist(file) {
		os.MkdirAll(common.BinDir+"/config", 0755)
		out, err := yaml.Marshal(c)
		if err == nil {
			ioutil.WriteFile(file, out, 0666)
		}
	}
	return c, nil
}
