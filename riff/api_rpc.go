package riff

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/teatak/riff/api"
	"github.com/teatak/riff/common"
	"gopkg.in/yaml.v2"
)

type Query struct{}

func (q *Query) SnapShot(_ struct{}, snap *string) error {
	*snap = server.SnapShot
	server.Logger.Printf(infoServerPrefix+"client get snapshot: %s", *snap)
	return nil
}

func (q *Query) Nodes(_ struct{}, nodes *api.Nodes) error {
	server.Logger.Printf(infoServerPrefix + "client get nodes")
	*nodes = server.api.Nodes()
	return nil
}

func (q *Query) Node(p api.ParamNode, node *api.Node) error {
	server.Logger.Printf(infoServerPrefix+"client get node %s", p.Name)
	n := server.api.Node(p.Name)
	if n == nil {
		return fmt.Errorf("node %s not found", p.Name)
	}
	*node = *n
	return nil
}

func (q *Query) Services(_ struct{}, services *api.Services) error {
	server.Logger.Printf(infoServerPrefix + "client get services")
	*services = server.api.Services()
	return nil
}

func (q *Query) Service(p api.ParamService, service *api.Service) error {
	server.Logger.Printf(infoServerPrefix+"client get service %s", p.Name)
	s := server.api.Service(p.Name, p.State)
	if s == nil {
		return fmt.Errorf("service %s not found", p.Name)
	}
	*service = *s
	return nil
}

type Mutation struct{}

func (m *Mutation) Reload(_ struct{}, result *bool) error {
	server.initServices()
	*result = true
	return nil
}

func (m *Mutation) RegisteService(config string, result *bool) error {
	serviceConfig := &Config{}
	if err := yaml.Unmarshal([]byte(config), serviceConfig); err != nil {
		return fmt.Errorf("config file error")
	}
	//write to config file
	file := common.BinDir + "/config/" + serviceConfig.Name + ".yml"
	_ = ioutil.WriteFile(file, []byte(config), 0666)

	s := server.Self.LoadService(serviceConfig.Name)
	server.Self.Services[s.Name] = s
	server.Self.Shutter()
	server.Shutter()

	*result = true
	server.Logger.Printf(infoServerPrefix + "client add new service")
	return nil
}

func (m *Mutation) UnregisteService(name string, result *bool) error {
	*result = true
	if server.Self.Services[name] != nil {
		if server.Self.Services[name].State == api.StateAlive {
			if err := server.Self.Services[name].Stop(); err != nil {
				server.Logger.Printf(errorServicePrefix+"client stop service %s error", name)
			}
		}
	}
	delete(server.Self.Services, name)

	file := common.BinDir + "/config/" + name + ".yml"
	run := common.BinDir + "/run/" + name + ".ver"
	update := common.BinDir + "/update/" + name

	_ = os.Remove(file)
	_ = os.Remove(run)
	_ = os.RemoveAll(update)

	server.Logger.Printf(infoServerPrefix+"client remove service %s", name)
	return nil
}

func (m *Mutation) Service(p api.ParamServiceMutation, reply *bool) (err error) {
	server.Logger.Printf(infoServerPrefix+"client %s service %s", strings.ToLower(p.Cmd.Name()), p.Name)
	s := server.Self.Services[p.Name]
	if s == nil {
		*reply = false
		err = fmt.Errorf("service %s not found", p.Name)
		return
	}
	switch p.Cmd {
	case api.CmdStart:
		_ = s.Start()
	case api.CmdStop:
		_ = s.Stop()
	case api.CmdRestart:
		_ = s.Restart()
	}
	*reply = true
	return nil
}

type Riff struct{}

// push request a digest
func (r *Riff) Request(snap string, digests *[]*Digest) error {
	if snap == server.SnapShot {
		*digests = nil
	} else {
		//build digest
		*digests = server.MakeDigest()
	}
	return nil
}

//push changes
func (r *Riff) PushDiff(diff []*Node, remoteDiff *[]*Node) error {
	if len(diff) == 0 {
		*remoteDiff = nil
	} else {
		*remoteDiff = server.MergeDiff(diff)
	}
	return nil
}
