package riff

import (
	"github.com/gimke/riff/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	IP         string     `yaml:"ip"`          //server ip
	Name       string     `yaml:"name"`        //server random name
	DataCenter string     `yaml:"data_center"` //server data center
	Join       string     `yaml:"join"`        //join address
	AutoPilot  bool       `yaml:"auto_pilot"`  //auto join node
	Addresses  *Addresses `yaml:"addresses"`
	Ports      *Ports     `yaml:"ports"`
	Fanout     int        `yaml:"fan_out"`
}
type Addresses struct {
	Http string `yaml:"http"` //http address
	Dns  string `yaml:"dns"`  //dns address
	Rpc  string `yaml:"rpc"`  //rpc address
}

type Ports struct {
	Http int `yaml:"http"` //http port default 8610
	Dns  int `yaml:"dns"`  //dns port default 8620
	Rpc  int `yaml:"rpc"`  //rpc port defalut 8630
}

func mergeConfig(src, dest *Config) {
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

func DefaultConfig() *Config {
	hostName, _ := os.Hostname()
	c := &Config{
		Name:       hostName,
		DataCenter: "dc1",
		AutoPilot:  false,
		Fanout:     3,
		Addresses: &Addresses{
			Http: "127.0.0.1",
			Dns:  "127.0.0.1",
			Rpc:  "",
		},
		Ports: &Ports{
			Http: common.DefaultHttpPort,
			Dns:  common.DefaultDnsPort,
			Rpc:  common.DefaultRpcPort,
		},
	}
	return c
}
func LoadConfig() *Config {
	file := common.BinDir + "/config/" + common.Name + ".yml"

	var c = DefaultConfig()

	if common.IsExist(file) {
		content, _ := ioutil.ReadFile(file)
		fc := DefaultConfig()
		err := yaml.Unmarshal(content, &fc)
		if err == nil {
			mergeConfig(fc, c)
		}
	} else {
		os.MkdirAll(common.BinDir+"/config", 0755)
		out, err := yaml.Marshal(c)
		if err == nil {
			ioutil.WriteFile(file, out, 0666)
		}
	}
	return c

}
