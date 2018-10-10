package riff

import (
	"github.com/gimke/riff/api"
)

type Resolver struct {
}

//Query
func (_ *Resolver) Riff() *RiffResolver {
	return &RiffResolver{}
}

func (_ *Resolver) Nodes() *[]*NodeResolver {
	var l []*NodeResolver
	for _, node := range server.api.Nodes() {
		l = append(l, &NodeResolver{node})
	}
	return &l
}

func (_ *Resolver) Node(args struct{ Name string }) *NodeResolver {
	return &NodeResolver{server.api.Node(args.Name)}
}

func (_ *Resolver) Server() *NodeResolver {
	return &NodeResolver{server.api.Node(server.Self.Name)}
}

func (_ *Resolver) Services() *[]*ServiceResolver {
	var l []*ServiceResolver
	for _, service := range server.api.Services() {
		l = append(l, &ServiceResolver{service})
	}
	return &l
}

func (_ *Resolver) Service(args struct {
	Name  string
	State string
}) *ServiceResolver {
	state := api.StateType_FromName(args.State)
	return &ServiceResolver{server.api.Service(args.Name, state)}
}
