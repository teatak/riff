package riff

import (
	"fmt"
	"github.com/gimke/riff/api"
	"strings"
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

func (q *Mutation) Service(p api.ParamServiceMutation, reply *bool) (err error) {
	server.Logger.Printf(infoServerPrefix+"client %s service %s", strings.ToLower(p.Cmd.String()), p.Name)
	s := server.Self.Services[p.Name]
	if s == nil {
		*reply = false
		err = fmt.Errorf("service %s not found", p.Name)
		return
	}
	switch p.Cmd {
	case api.CmdStart:
		err = s.Start()
		break
	case api.CmdStop:
		err = s.Stop()
		break
	case api.CmdRestart:
		err = s.Restart()
		break
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
