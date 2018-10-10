package riff

import "github.com/gimke/riff/common"

type RiffResolver struct {
}

func (r *RiffResolver) GitBranch() string {
	return common.GitBranch
}
func (r *RiffResolver) GitSha() string {
	return common.GitSha
}
func (r *RiffResolver) Version() string {
	return common.Version
}
