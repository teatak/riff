package start

import (
	"fmt"
	"github.com/gimke/riff/common"
	"github.com/gimke/riff/riff"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var configText = ""

func init() {
	hostName, _ := os.Hostname()
	id := common.GenerateID(12)
	configText = `#id: auto generated id
#name: node name
#data_center: data center

#addresses:
#  http: 172.0.0.1
#  dns: 172.0.0.1
#  rpc: 192.169.1.2

#ports:
#  http: 8610
#  dns: 8620
#  rpc: 8630

id: ` + id + `
name: ` + hostName + `
data_center: dc1
`
	initConfig()
}

func isExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func initConfig() {
	file := common.BinDir + "/config/" + common.Name + ".yml"
	if !isExist(file) {
		os.MkdirAll(common.BinDir+"/config", 0755)
		ioutil.WriteFile(file, []byte(configText), 0666)
	}
}

func loadConfig() (*riff.Config, error) {
	file := common.BinDir + "/config/" + common.Name + ".yml"
	if !isExist(file) {
		return nil, fmt.Errorf("file not exist %s", file)
	}
	content, _ := ioutil.ReadFile(file)
	var c = &riff.Config{}
	err := yaml.Unmarshal(content, &c)
	if err != nil {
		return nil, err
	}
	if c.Addresses == nil {
		c.Addresses = &riff.Addresses{}
	}
	if c.Ports == nil {
		c.Ports = &riff.Ports{
			Http: 8610,
			Dns:  8620,
			Rpc:  8630,
		}
	} else {
		if c.Ports.Http == 0 {
			c.Ports.Http = 8610
		}
		if c.Ports.Dns == 0 {
			c.Ports.Dns = 8620
		}
		if c.Ports.Rpc == 0 {
			c.Ports.Rpc = 8630
		}
	}

	return c, nil
}
