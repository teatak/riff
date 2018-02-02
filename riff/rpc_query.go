package riff

import (
	"fmt"
	"github.com/gimke/riff/api"
)

type Query struct{}

// Ping is used to just check for connectivity
func (q *Query) SnapShot(_ struct{}, snap *string) error {
	server.Logger.Printf(infoServerPrefix+"client get snapshot: %s", *snap)
	*snap = server.SnapShot
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
		return fmt.Errorf("node %s Not found", p.Name)
	}
	*node = *n
	return nil
}

func (q *Query) Services(p api.ParamNode, services *api.Services) error {
	server.Logger.Printf(infoServerPrefix + "client get services")
	*services = server.api.Services()
	return nil
}

func (q *Query) Service(p api.ParamService, service *api.Service) error {
	server.Logger.Printf(infoServerPrefix+"client get service %s", p.Name)
	s := server.api.Service(p.Name, p.State)
	if s == nil {
		return fmt.Errorf("service %s Not found", p.Name)
	}
	*service = *s
	return nil
}
