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
