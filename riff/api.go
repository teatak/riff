package riff

import (
	"github.com/teatak/riff/api"
)

type API struct{}

func (a *API) makeNode(n *Node) *api.Node {
	node := &api.Node{
		Name:       n.Name,
		DataCenter: n.DataCenter,
		IP:         n.IP,
		RpcPort:    n.RpcPort,
		HttpPort:   n.HttpPort,
		State:      n.State,
		SnapShot:   n.SnapShot,
		IsSelf:     n.IsSelf,
		Version:    int(n.Version),
	}
	return node
}
func (a *API) makeNestNode(n *Node, s *Service, resolveState api.StateType) *api.NestNode {
	ip := n.IP
	if s.Ip != "" {
		ip = s.Ip
	}
	node := &api.NestNode{
		Name:          n.Name,
		DataCenter:    n.DataCenter,
		IP:            ip,
		Port:          s.Port,
		RpcPort:       n.RpcPort,
		State:         resolveState,
		Version:       int(n.Version),
		SnapShot:      n.SnapShot,
		IsSelf:        n.IsSelf,
		Config:        s.Config,
		StatusContent: s.StatusContent,
		StartTime:     s.StartTime,
		Progress: &api.NestProgress{
			Current:    s.Progress.Current,
			Total:      s.Progress.Total,
			InProgress: s.Progress.InProgress,
		},
	}
	return node
}

func (a *API) makeService(s *Service) *api.Service {
	return &api.Service{
		Name: s.Name,
	}
}

func (a *API) makeNestService(s *Service) *api.NestService {
	service := &api.NestService{
		Name:          s.Name,
		IP:            s.Ip,
		Port:          s.Port,
		State:         s.State,
		Config:        s.Config,
		StatusContent: s.StatusContent,
		StartTime:     s.StartTime,
		Progress: &api.NestProgress{
			Current:    s.Progress.Current,
			Total:      s.Progress.Total,
			InProgress: s.Progress.InProgress,
		},
	}
	return service
}

func (a *API) Nodes() api.Nodes {
	keys := server.Keys()
	nodes := make([]*api.Node, 0, len(keys))
	for _, key := range keys {
		if n := server.GetNode(key); n != nil {
			nodes = append(nodes, a.makeNode(n))
		}
	}
	return nodes
}

func (a *API) Node(name string) *api.Node {
	if n := server.GetNode(name); n != nil {
		node := a.makeNode(n)
		for _, key := range n.Services.Keys() {
			s := n.Services[key]
			node.NestServices = append(node.NestServices, a.makeNestService(s))
		}
		return node
	}
	return nil
}

func (a *API) Services() api.Services {
	keys := server.Keys()
	helper := make(map[string]string, 0)
	services := make([]*api.Service, 0)
	for _, key := range keys {
		if n := server.GetNode(key); n != nil {
			for _, skey := range n.Services.Keys() {
				if _, ok := helper[skey]; !ok {
					helper[skey] = skey
					service := a.makeService(n.Services[skey])
					services = append(services, service)
				}
			}
		}
	}
	return services
}

func (a *API) Service(name string, state api.StateType) *api.Service {
	keys := server.Keys()
	var service *api.Service
	nodes := make(api.NestNodes, 0)
	for _, key := range keys {
		if n := server.GetNode(key); n != nil {
			for _, s := range n.Services {
				if s.Name == name {
					if service == nil {
						service = &api.Service{
							Name: s.Name,
						}
					}

					resolveState := api.StateAlive
					if n.State == api.StateAlive {
						resolveState = s.State
					} else {
						resolveState = n.State
					}

					if resolveState&state == resolveState {
						node := a.makeNestNode(n, s, resolveState)
						nodes = append(nodes, node)
					}
				}
			}
		}
	}
	if service != nil {
		service.NestNodes = nodes
	}
	return service
}

func (a *API) Start(name string) bool {
	if s, ok := server.Self.Services[name]; ok {
		err := s.Start()
		if err != nil {
			return true
		}
	} else {
		return false
	}
	return true
}
func (a *API) Stop(name string) bool {
	if s, ok := server.Self.Services[name]; ok {
		err := s.Stop()
		if err != nil {
			return true
		}
	} else {
		return false
	}
	return true
}
func (a *API) Restart(name string) bool {
	if s, ok := server.Self.Services[name]; ok {
		err := s.Restart()
		if err != nil {
			return true
		}
	} else {
		return false
	}
	return true
}
