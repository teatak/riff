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
		for _, node := range server.api.Service(r.service.Name,api.StateAll).NestNodes {
			l = append(l, &NestNodeResover{*node})
		}
		return &l
	}}