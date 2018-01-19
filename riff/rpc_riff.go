package riff

type Riff struct{}

// push request a digest
func (r *Riff) Request(snap string, digests *[]*Digest) error {
	if snap == server.SnapShot {
		*digests = nil
	} else {
		//build digest
		*digests = server.MakeDigest()
	}
	return nil
}

//push changes
func (r *Riff) PushDiff(diff []*Node, remoteDiff *[]*Node) error {
	if len(diff) == 0 {
		*remoteDiff = nil
	} else {
		*remoteDiff = server.MergeDiff(diff)
	}
	return nil
}
