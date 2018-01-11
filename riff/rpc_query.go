package riff

type Query struct {
	server *Server
}

// Ping is used to just check for connectivity
func (q *Query) SnapShot(_ struct{}, snap *string) error {
	*snap = q.server.SnapShot
	q.server.logger.Printf(infoRpcPrefix+"client get snapshot: %s", *snap)
	return nil
}

func (q *Query) Nodes(_ struct{}, nodes *[]*Node) error {
	*nodes = q.server.Slice()
	q.server.logger.Printf(infoRpcPrefix + "client get nodes list")
	return nil
}
