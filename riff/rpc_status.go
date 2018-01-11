package riff

type Status struct {
	server *Server
}

// Ping is used to just check for connectivity
func (s *Status) Ping(_ struct{}, reply *string) error {
	s.server.Logger.Printf("ping")
	*reply = "pong"
	return nil
}

// Leader is used to get the address of the leader
func (s *Status) Leader(args struct{}, reply *string) error {
	leader := "node1"
	if leader != "" {
		*reply = leader
	} else {
		*reply = ""
	}
	return nil
}
