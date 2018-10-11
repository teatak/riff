package riff

import (
	"github.com/gimke/riff/api"
)

type NodeResolver struct {
	node *api.Node
}

func (r *NodeResolver) DataCenter() string {
	return r.node.DataCenter
}
func (r *NodeResolver) Ip() string {
	return r.node.IP
}
func (r *NodeResolver) IsSelf() bool {
	return r.node.IsSelf
}
func (r *NodeResolver) Name() string {
	return r.node.Name
}
func (r *NodeResolver) RpcPort() int32 {
	return int32(r.node.RpcPort)
}
func (r *NodeResolver) HttpPort() int32 {
	return int32(r.node.HttpPort)
}
func (r *NodeResolver) Services() *[]*NestServiceResover {
	var l []*NestServiceResover
	if r.node.NestServices != nil {
		for _, service := range r.node.NestServices {
			l = append(l, &NestServiceResover{*service})
		}
		return &l
	} else {
		for _, service := range server.api.Node(r.node.Name).NestServices {
			l = append(l, &NestServiceResover{*service})
		}
		return &l
	}
	//return &[]*NestServiceResover{{}}
}
func (r *NodeResolver) SnapShot() string {
	return r.node.SnapShot
}
func (r *NodeResolver) State() string {
	return r.node.State.Name()
}
func (r *NodeResolver) Version() int32 {
	return int32(r.node.Version)
}
