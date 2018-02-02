package api

type Nodes []*Node

type Node struct {
	Name         string    `json:"name"`
	DataCenter   string    `json:"dataCenter"`
	IP           string    `json:"ip"`
	Port         int       `json:"port,omitempty"`
	Version      int       `json:"version"`
	State        StateType `json:"state"`
	SnapShot     string    `json:"snapShot,omitempty"`
	IsSelf       bool      `json:"isSelf,omitempty"`
	NestServices `json:"services,omitempty"`
}

type NestNodes []*NestNode

type NestNode struct {
	Name       string    `json:"name"`
	DataCenter string    `json:"dataCenter"`
	IP         string    `json:"ip"`
	Port       int       `json:"port,omitempty"`
	Version    int       `json:"version"`
	State      StateType `json:"state"`
	SnapShot   string    `json:"snapShot,omitempty"`
	IsSelf     bool      `json:"isSelf,omitempty"`
	Config     string    `json:"config,omitempty"`
}

type Services []*Service

type Service struct {
	Name      string `json:"name"`
	NestNodes `json:"nodes,omitempty"`
}

type NestServices []*NestService

type NestService struct {
	Name   string    `json:"name"`
	Port   int       `json:"port,omitempty"`
	State  StateType `json:"state,omitempty"`
	Config string    `json:"config,omitempty"`
}

type ParamNode struct {
	Name  string    `json:"name"`
	State StateType `json:"state,omitempty"`
}

type ParamServices struct {
	Name  string    `json:"name"`
	State StateType `json:"state,omitempty"`
}
