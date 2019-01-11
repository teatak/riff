package api

type StateType int

const (
	StateAlive StateType = 1 << iota
	StateSuspect
	StateDead
	StateAll = StateAlive | StateSuspect | StateDead
)

func StateType_FromName(name string) StateType {
	switch name {
	case "ALIVE":
		return StateAlive
	case "SUSPECT":
		return StateSuspect
	case "DEAD":
		return StateDead
	case "ALL":
		return StateAll
	}
	return StateAll
}
func (s StateType) Name() string {
	switch s {
	case StateAlive:
		return "ALIVE"
	case StateSuspect:
		return "SUSPECT"
	case StateDead:
		return "DEAD"
	case StateAll:
		return "ALL"
	default:
		return "UNKOWN"
	}
}

func (s StateType) Value() int {
	return int(s)
}

type CmdType int

const (
	CmdStart CmdType = 1 << iota
	CmdStop
	CmdRestart
)

func CmdType_FromName(name string) CmdType {
	switch name {
	case "START":
		return CmdStart
	case "STOP":
		return CmdStop
	case "RESTART":
		return CmdRestart
	}
	return CmdStart
}
func (c CmdType) Name() string {
	switch c {
	case CmdStart:
		return "START"
	case CmdStop:
		return "STOP"
	case CmdRestart:
		return "RESTART"
	default:
		return "UNKOWN"
	}
}

func (c CmdType) Value() int {
	return int(c)
}

type Nodes []*Node

type Node struct {
	Name         string    `json:"name"`
	DataCenter   string    `json:"dataCenter"`
	IP           string    `json:"ip"`
	RpcPort      int       `json:"rpcPort,omitempty"`
	HttpPort     int       `json:"httpPort,omitempty"`
	Version      int       `json:"version"`
	State        StateType `json:"state"`
	SnapShot     string    `json:"snapShot,omitempty"`
	IsSelf       bool      `json:"isSelf,omitempty"`
	NestServices `json:"services,omitempty"`
}

type NestNodes []*NestNode

type NestNode struct {
	Name          string    `json:"name"`
	DataCenter    string    `json:"dataCenter"`
	IP            string    `json:"ip"`
	Port          int       `json:"port,omitempty"`
	RpcPort       int       `json:"rpcPort,omitempty"`
	Version       int       `json:"version"`
	State         StateType `json:"state"`
	SnapShot      string    `json:"snapShot,omitempty"`
	IsSelf        bool      `json:"isSelf,omitempty"`
	Config        string    `json:"config,omitempty"`
	StatusContent string    `json:"statusContent,omitempty"`
}

type Services []*Service

type Service struct {
	Name      string `json:"name"`
	NestNodes `json:"nodes,omitempty"`
}

type NestServices []*NestService

type NestService struct {
	Name          string    `json:"name"`
	IP            string    `json:"ip"`
	Port          int       `json:"port,omitempty"`
	State         StateType `json:"state,omitempty"`
	Config        string    `json:"config,omitempty"`
	StatusContent string    `json:"statusContent,omitempty"`
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
