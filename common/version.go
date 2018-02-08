package common

/*
	defined in ../scripts/build.sh

	TYPE="release"
	VERSION="$(cat version)"
	GITSHA="$(git rev-parse HEAD)"
	GITBRANCH="$(git rev-parse --abbrev-ref HEAD)"
*/
var Type = "debug"
var Version = "v0.0.1.dev"
var GitSha = "nil"
var GitBranch = "nil"

func IsRelease() bool {
	if Type == "release" {
		return true
	}
	return false
}
