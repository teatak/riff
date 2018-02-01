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
	Name string `json:"name"`
	//IP        string    `json:"ip,omitempty"`
	//Port      int       `json:"port,omitempty"`
	//State     StateType `json:"state,omitempty"`
	//Config    string    `json:"config,omitempty"`
	NestNodes `json:"nodes,omitempty"`
}

type NestServices []*NestService

type NestService struct {
	Name   string    `json:"name"`
	IP     string    `json:"ip,omitempty"`
	Port   int       `json:"port,omitempty"`
	State  StateType `json:"state,omitempty"`
	Config string    `json:"config,omitempty"`
}
