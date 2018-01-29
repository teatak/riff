package riff

import (
	"github.com/gimke/riff/api"
)

type API struct{}

func (a *API) Nodes() api.Nodes {
	keys := server.Keys()
	nodes := make([]*api.Node, 0, len(keys))
	for _, key := range keys {
		if n := server.GetNode(key); n != nil {
			node := &api.Node{
				Name:       n.Name,
				DataCenter: n.DataCenter,
				IP:         n.IP,
				Port:       n.Port,
				State:      n.State,
				SnapShot:   n.SnapShot,
				Version:    int(n.Version),
			}
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (a *API) Node(name string) *api.Node {
	if n := server.GetNode(name); n != nil {
		node := &api.Node{
			Name:       n.Name,
			DataCenter: n.DataCenter,
			IP:         n.IP,
			Port:       n.Port,
			State:      n.State,
			SnapShot:   n.SnapShot,
			Version:    int(n.Version),
		}
		for _, key := range n.Services.Keys() {
			s := n.Services[key]
			service := &api.NestService{
				Name:   s.Name,
				IP:     s.IP,
				Port:   s.Port,
				State:  s.State,
				Config: s.Config,
			}
			node.NestServices = append(node.NestServices, service)
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
					service := &api.Service{
						Name: n.Services[skey].Name,
					}
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
					if n.State&state == n.State || s.State&state == s.State {
						node := &api.NestNode{
							Name:       n.Name,
							DataCenter: n.DataCenter,
							IP:         n.IP,
							Port:       s.Port,
							State:      s.State,
							Version:    int(n.Version),
							SnapShot:   n.SnapShot,
							Config:     s.Config,
						}
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

//func (a *API) cloneNode(n *Node) (node *api.Node) {
//	node = &api.Node{
//		Name:       n.Name,
//		DataCenter: n.DataCenter,
//		IP:         n.IP,
//		Port:       n.Port,
//		State:      n.State,
//		SnapShot:   n.SnapShot,
//		Version:    int(n.Version),
//	}
//	return
//}

//func (a *API) cloneService(s *Service) (service *api.NestService) {
//	service = &api.NestService{
//		Name:   s.Name,
//		IP:     s.IP,
//		Port:   s.Port,
//		State:  s.State,
//		Config: s.Config,
//	}
//	return service
//}
