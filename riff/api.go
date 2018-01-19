package riff

import (
	"github.com/gimke/riff/api"
)

type API struct {}

func (a *API) Nodes() api.Nodes {
	keys := server.Keys()
	nodes := make([]*api.Node, 0, len(keys))
	for _, key := range keys {
		if n := server.GetNode(key); n != nil {
			nodes = append(nodes, a.cloneNode(n))
		}
	}
	return nodes
}

func (a *API) cloneNode(n *Node) (node *api.Node) {
	node = &api.Node{
		Name:       n.Name,
		DataCenter: n.DataCenter,
		IP:         n.IP,
		Port:       n.Port,
		State:      int(n.State),
		SnapShot:   n.SnapShot,
	}
	return
}
