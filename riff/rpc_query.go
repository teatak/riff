package riff

import (
	"fmt"
	"github.com/gimke/riff/api"
)

type Query struct{}

// Ping is used to just check for connectivity
func (q *Query) SnapShot(_ struct{}, snap *string) error {
	*snap = server.SnapShot
	server.Logger.Printf(infoServerPrefix+"client get snapshot: %s", *snap)
	return nil
}

func (q *Query) Nodes(_ struct{}, nodes *api.Nodes) error {
	*nodes = server.api.Nodes()
	server.Logger.Printf(infoServerPrefix + "client get nodes list")
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
