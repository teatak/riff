package riff

import "log"

type Riff struct {
	server *Server
}

// push request a digest
func (r *Riff) Pull(snap string, digest *Nodes) error {
	log.Printf(infoRpcPrefix+"riff pull snapshot: %s", snap)
	if snap == r.server.SnapShot {
		digest = nil
	} else {
		//build digest
		*digest = r.server.MakeDigest()
	}
	return nil
}

//push changes
func (r *Riff) PushDiff(diff Nodes, remoteDiff *Nodes) error {
	log.Printf(infoRpcPrefix+"riff push diff: %v", diff)
	return nil
}
