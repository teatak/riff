package riff

import "sort"

type Nodes map[string]*Node

type Node struct {
	Services
	Name string
}

func (ns *Nodes) sort() []string {
	var keys = make([]string, 0, 0)
	for key, _ := range *ns {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func (n *Node) AddService(s *Service) {
	if n.Services == nil {
		n.Services = make(map[string]*Service)
	}
	n.Services[s.Name] = s
}
