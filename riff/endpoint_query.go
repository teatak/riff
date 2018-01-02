package riff

import "log"

type Query struct {
	server *Server
}

// Ping is used to just check for connectivity
func (q *Query) SnapShot(_ struct{}, snap *string) error {
	*snap = q.server.SnapShot
	log.Printf(infoRpcPrefix+"client get snapshot: %s",*snap)
	return nil
}

func (q *Query) Nodes(_ struct{}, nodes *Nodes) error {
	*nodes = q.server.Nodes
	log.Printf(infoRpcPrefix+"client get nodes list")
	return nil
}
