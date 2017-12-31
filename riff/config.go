package riff

import (
	"net"
	"fmt"
)

type Config struct {
	IP   net.IP
	Port int
	Name string
}

func NewConfig(addr,name string) (*Config, error) {
	ipAddr, err := net.ResolveTCPAddr("", addr)
	if err != nil {
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}
	return &Config{
		IP:   ipAddr.IP,
		Port: ipAddr.Port,
		Name: name,
	}, nil

}
