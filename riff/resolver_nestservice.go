package riff

import "github.com/gimke/riff/api"

type NestServiceResover struct {
	service api.NestService
}

func (r *NestServiceResover) Config() string {
	return r.service.Config
}
func (r *NestServiceResover) StatusContent() string {
	return r.service.StatusContent
}
func (r *NestServiceResover) Name() string {
	return r.service.Name
}
func (r *NestServiceResover) Ip() string {
	return r.service.IP
}
func (r *NestServiceResover) Port() int32 {
	return int32(r.service.Port)
}
func (r *NestServiceResover) State() string {
	return r.service.State.Name()
}
func (r *NestServiceResover) StartTime() int32 {
	return int32(r.service.StartTime.Unix())
}
func (r *NestServiceResover) Progress() *NestProgressResover {
	return &NestProgressResover{r.service.Progress}
}

type NestProgressResover struct {
	progress *api.NestProgress
}

func (r *NestProgressResover) Current() int32 {
	return int32(r.progress.Current)
}

func (r *NestProgressResover) Total() int32 {
	return int32(r.progress.Total)
}

func (r *NestProgressResover) InProgress() bool {
	return r.progress.InProgress
}
