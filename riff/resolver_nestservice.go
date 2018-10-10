package riff

import "github.com/gimke/riff/api"

type NestServiceResover struct {
	service api.NestService
}
func (r *NestServiceResover) Config() string {
	return r.service.Config
}
func (r *NestServiceResover) Name() string {
	return r.service.Name
}
func (r *NestServiceResover) Port() int32 {
	return int32(r.service.Port)
}
func (r *NestServiceResover) State() string {
	return r.service.State.Name()
}