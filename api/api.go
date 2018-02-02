package api

type StateType int

const (
	StateAlive StateType = 1 << iota
	StateSuspect
	StateDead
	StateAll = StateAlive | StateSuspect | StateDead
)

func (s StateType) String() string {
	switch s {
	case StateAlive:
		return "Alive"
		break
	case StateSuspect:
		return "Suspect"
		break
	case StateDead:
		return "Dead"
		break
	case StateAll:
		return "All"
		break
	}
	return "Unknow"
}

type CmdType int

const (
	CmdStart CmdType = 1 << iota
	CmdStop
	CmdRestart
)

func (c CmdType) String() string {
	switch c {
	case CmdStart:
		return "Start"
		break
	case CmdStop:
		return "Stop"
		break
	case CmdRestart:
		return "Restart"
		break
	}
	return "Unknow"
}

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
	Name string `json:"name"`
}

type ParamService struct {
	Name  string    `json:"name"`
	State StateType `json:"state,omitempty"`
}

type ParamServiceMutation struct {
	Name string  `json:"name"`
	Cmd  CmdType `json:"cmd,omitempty"`
}
