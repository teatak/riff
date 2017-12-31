package riff

import (
	"fmt"
	"net"
)

type Config struct {
	IP         net.IP
	Port       int
	Name       string
	DataCenter string
}

func NewConfig(addr, name, dataCenter string) (*Config, error) {
	ipAddr, err := net.ResolveTCPAddr("", addr)
	if err != nil {
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}
	return &Config{
		IP:         ipAddr.IP,
		Port:       ipAddr.Port,
		Name:       name,
		DataCenter: dataCenter,
	}, nil

}
