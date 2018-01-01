package riff

import "log"

type Query struct {
	server *Server
}

// Ping is used to just check for connectivity
func (q *Query) SnapShort(_ struct{}, snap *string) error {
	*snap = q.server.SnapShort
	log.Printf(infoRpcPrefix+"client get snapshort: %s",*snap)
	return nil
}

func (q *Query) Nodes(_ struct{}, nodes *Nodes) error {
	*nodes = q.server.Nodes
	log.Printf(infoRpcPrefix+"client get nodes list")
	return nil
}
