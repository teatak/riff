package riff

import "github.com/gimke/riff/api"

type Query struct{}

// Ping is used to just check for connectivity
func (q *Query) SnapShot(_ struct{}, snap *string) error {
	*snap = server.SnapShot
	server.Logger.Printf(infoRpcPrefix+"client get snapshot: %s", *snap)
	return nil
}

func (q *Query) Nodes(_ struct{}, nodes *api.Nodes) error {
	*nodes = server.api.Nodes()
	server.Logger.Printf(infoRpcPrefix + "client get nodes list")
	return nil
}
