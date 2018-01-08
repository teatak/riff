package riff

type Riff struct {
	server *Server
}

// push request a digest
func (r *Riff) Request(snap string, digest *[]Node) error {
	r.server.logger.Printf(infoRpcPrefix+"riff request snapshot: %s", snap)
	if snap == r.server.SnapShot {
		*digest = nil
	} else {
		//build digest
		*digest = r.server.MakeDigest()
	}
	return nil
}

//push changes
func (r *Riff) PushDiff(diff []Node, remoteDiff *[]Node) error {
	*remoteDiff = r.server.MergeDiff(diff)
	return nil
}
