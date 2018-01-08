package riff

type Riff struct {
	server *Server
}

// push request a digest
func (r *Riff) Request(snap string, digests *Digests) error {
	if snap == r.server.SnapShot {
		*digests = nil
	} else {
		//build digest
		*digests = r.server.MakeDigest()
	}
	return nil
}

//push changes
func (r *Riff) PushDiff(diff Nodes, remoteDiff *Nodes) error {
	if len(diff) == 0 {
		*remoteDiff = nil
	} else {
		*remoteDiff = r.server.MergeDiff(diff)
	}
	return nil
}
