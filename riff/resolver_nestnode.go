package riff

import "github.com/gimke/riff/api"

type NestNodeResover struct {
	node api.NestNode
}

func (r *NestNodeResover) Config() string {
	return r.node.Config
}
func (r *NestNodeResover) StatusContent() string {
	return r.node.StatusContent
}
func (r *NestNodeResover) DataCenter() string {
	return r.node.DataCenter
}
func (r *NestNodeResover) Ip() string {
	return r.node.IP
}
func (r *NestNodeResover) IsSelf() bool {
	return r.node.IsSelf
}
func (r *NestNodeResover) Name() string {
	return r.node.Name
}
func (r *NestNodeResover) Port() int32 {
	return int32(r.node.Port)
}
func (r *NestNodeResover) RpcPort() int32 {
	return int32(r.node.RpcPort)
}
func (r *NestNodeResover) SnapShot() string {
	return r.node.SnapShot
}
func (r *NestNodeResover) State() string {
	return r.node.State.Name()
}
func (r *NestNodeResover) Version() int32 {
	return int32(r.node.Version)
}
