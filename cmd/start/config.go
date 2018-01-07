package start

import (
	"fmt"
	"github.com/gimke/riff/common"
	"github.com/gimke/riff/riff"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"strconv"
)

var configText = ""

func init() {
	//	hostName, _ := os.Hostname()
	//	id := common.GenerateID(12)
	//	configText = `#id: auto generated id
	//#name: node name
	//#data_center: data center
	//
	//#addresses:
	//#  http: 172.0.0.1
	//#  dns: 172.0.0.1
	//#  rpc: 192.169.1.2
	//
	//#ports:
	//#  http: 8610
	//#  dns: 8620
	//#  rpc: 8630
	//
	//id: ` + id + `
	//name: ` + hostName + `
	//data_center: dc1
	//`
	//	initConfig()
}

func isExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func initConfig() {
	//file := common.BinDir + "/config/" + common.Name + ".yml"
	//if !isExist(file) {
	//	os.MkdirAll(common.BinDir+"/config", 0755)
	//	ioutil.WriteFile(file, []byte(configText), 0666)
	//}
}

func defaultConfig() *riff.Config {
	id := common.GenerateID(12)
	hostName, _ := os.Hostname()
	c := &riff.Config{
		Id:         id,
		Name:       hostName,
		DataCenter: "dc1",
		AutoPilot:  false,
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

func adviseRpc(addr string) string {
	var advise string
	if common.IsAny(addr) {
		var addrs []*net.IPNet
		var err error
		//detect ip
		var addrtype string

		switch {
		case common.IsAnyV4(addr):
			addrtype = "private IPv4"
			addrs, err = common.GetPrivateIPv4()
			if err != nil {
				fmt.Println("Error detecting %s address: %s", addrtype, err)
			}
			break
		case common.IsAnyV6(addr):
			addrtype = "public IPv6"
			addrs, err = common.GetPublicIPv6()
			if err != nil {
				fmt.Println("Error detecting %s address: %s", addrtype, err)
			}
			break
		}
		if len(addrs) > 0 {
			advise = addrs[0].String()
		}
	}
	return advise
}

func getIpPort(ipPort string) (ip string, port int, err error) {
	index := strings.LastIndex(ipPort,":")
	if ipPort == "" {
		err = fmt.Errorf("empty ip and port\n")
		return
	}
	if index > -1 && index < len(ipPort) {
		ip = ipPort[0:index]
		port,err = strconv.Atoi(ipPort[index+1:])
		if err != nil {
			ip = ipPort
			port = 0
			err = nil
		}
	} else {
		ip = ipPort
		port = 0
		err = nil
	}
	if ip != "" {
		ipaddr := net.ParseIP(ip)
		if ipaddr == nil {
			err = fmt.Errorf("error ip and port\n")
		}
	}
	return
}

func loadConfig(cmd *cmd) (*riff.Config, error) {
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
	advise := adviseRpc("0.0.0.0")
	c.IP = advise

	//if cmd.http != "" {
	var host string
	var port int
	var err error
	//http
	host,port,err = getIpPort(cmd.http)
	if err == nil {
		if host != "" {
			c.Addresses.Http = host
		}
		if port != 0 {
			c.Ports.Http = port
		}
	}
	host,port,err = getIpPort(cmd.dns)
	if err == nil {
		if host != "" {
			c.Addresses.Dns = host
		}
		if port != 0 {
			c.Ports.Dns = port
		}
	}
	host,port,err = getIpPort(cmd.rpc)
	if err == nil {
		if host != "" {
			c.Addresses.Rpc = host
		}
		if port != 0 {
			c.Ports.Rpc = port
		}
	}
	if c.Addresses.Rpc == "" {
		ip,_,_ := net.ParseCIDR(advise)
		c.Addresses.Rpc = ip.String()
	}
	return c, nil
}
