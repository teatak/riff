package riff

import "github.com/gimke/riff/api"

type ServiceResolver struct {
	service *api.Service
}

func (r *ServiceResolver) Name() string {
	return r.service.Name
}

func (r *ServiceResolver) Nodes() *[]*NestNodeResover {
	var l []*NestNodeResover
	if r.service.NestNodes != nil {
		for _, node := range r.service.NestNodes {
			l = append(l, &NestNodeResover{*node})
		}
		return &l
	} else {
		for _, node := range server.api.Service(r.service.Name, api.StateAll).NestNodes {
			l = append(l, &NestNodeResover{*node})
		}
		return &l
	}
}

type MutationServiceInput struct {
	Port int32
	Cmd  string
	Name string
	Ip   string
}

type MutationService struct {
	cmd     string
	error   string
	ip      string
	name    string
	port    int
	success bool
}

func (r *MutationService) Cmd() string {
	return r.cmd
}
func (r *MutationService) Error() string {
	return r.error
}
func (r *MutationService) Ip() string {
	return r.ip
}
func (r *MutationService) Name() string {
	return r.name
}
func (r *MutationService) Port() int32 {
	return int32(r.port)
}
func (r *MutationService) Success() bool {
	return r.success
}
