package riff

import (
	"net"
	"sort"
	"strconv"
)

type Services map[string]*Service

type Service struct {
	Name    string
	Addr    net.IP
	Port    uint16
	Version uint64
}

func (n *Service) Address() string {
	return net.JoinHostPort(n.Addr.String(), strconv.Itoa(int(n.Port)))
}

func (ss *Services) sort() []string {
	var keys = make([]string, 0, 0)
	for key, _ := range *ss {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
