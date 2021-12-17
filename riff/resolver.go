package riff

import (
	"errors"
	"net"
	"strconv"

	"github.com/teatak/riff/api"
)

type Resolver struct {
}

//Query
func (*Resolver) Riff() *RiffResolver {
	return &RiffResolver{}
}

func (*Resolver) Nodes() *[]*NodeResolver {
	var l []*NodeResolver
	for _, node := range server.api.Nodes() {
		l = append(l, &NodeResolver{node})
	}
	return &l
}

func (*Resolver) Node(args struct{ Name string }) *NodeResolver {
	node := server.api.Node(args.Name)
	if node == nil {
		return nil
	} else {
		return &NodeResolver{node}
	}
}

func (*Resolver) Server() *NodeResolver {
	return &NodeResolver{server.api.Node(server.Self.Name)}
}

func (*Resolver) Services() *[]*ServiceResolver {
	var l []*ServiceResolver
	for _, service := range server.api.Services() {
		l = append(l, &ServiceResolver{service})
	}
	return &l
}

func (*Resolver) Service(args struct {
	Name  string
	State *string
}) *ServiceResolver {
	state := api.StateAll
	if args.State != nil {
		state = api.StateType_FromName(*args.State)
	}
	service := server.api.Service(args.Name, state)
	if service == nil {
		return nil
	} else {
		return &ServiceResolver{service}
	}
}

func (*Resolver) MutationService(args struct {
	Services *[]*MutationServiceInput
}) *[]*MutationService {
	var l []*MutationService
	for _, service := range *args.Services {
		result := &MutationService{
			cmd:  service.Cmd,
			ip:   service.Ip,
			port: int(service.Port),
		}
		if err := mutationService(service.Name, net.JoinHostPort(service.Ip, strconv.Itoa(int(service.Port))), api.CmdType_FromName(service.Cmd)); err != nil {
			result.error = err.Error()
			result.success = false
		} else {
			result.error = ""
			result.success = true
		}
		l = append(l, result)
	}
	return &l
}

func (*Resolver) RegisteService(args struct {
	Node struct {
		Ip   string
		Port int32
	}
	Config string
}) (*bool, error) {
	result := true
	if err := registeService(net.JoinHostPort(args.Node.Ip, strconv.Itoa(int(args.Node.Port))), args.Config); err != nil {
		result = false
	} else {
		result = true
	}
	if result {
		return &result, nil
	} else {
		return &result, errors.New("register service fail")
	}
}

func (*Resolver) UnregisteService(args struct {
	Node struct {
		Ip   string
		Port int32
	}
	Name string
}) *bool {
	result := true
	if err := unregisteService(net.JoinHostPort(args.Node.Ip, strconv.Itoa(int(args.Node.Port))), args.Name); err != nil {
		result = false
	} else {
		result = true
	}
	return &result
}
