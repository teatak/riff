package riff

type Query struct {
	server *Server
}

// Ping is used to just check for connectivity
func (q *Query) SnapShort(_ struct{}, reply *string) error {
	*reply = q.server.SnapShort
	return nil
}

func (q *Query) Nodes(_ struct{}, nodes *Nodes) error {
	*nodes = q.server.Nodes
	return nil
}
