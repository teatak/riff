package riff

import (
	"net"
	"sort"
	"strconv"
	"time"
)

type Services map[string]*Service

type Service struct {
	Name        string
	IP          string
	Port        uint16
	Version     uint64
	State       stateType // Current state
	StateChange time.Time // Time last state change happened
}

func (n *Service) Address() string {
	return net.JoinHostPort(n.IP, strconv.Itoa(int(n.Port)))
}

func (ss *Services) Keys() []string {
	var keys = make([]string, 0, 0)
	for key, _ := range *ss {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
